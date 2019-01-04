package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"reflect"
	"sort"
	"strconv"
)

const port string = ":8000" // Configurable localhost port

type validData struct {
	DataHashOne string `json:"DataHashOne"`
	DataHashTwo string `json:"DataHashTwo"`
	Valid       bool   `json:"valid"`
}

type byKey [][]string

type byFirstValue [][][]string

func (k byKey) Len() int {
	return len(k)
}

func (k byKey) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

func (k byKey) Less(i, j int) bool {
	return k[i][0] < k[j][0]
}

func (v byFirstValue) Len() int {
	return len(v)
}

func (v byFirstValue) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v byFirstValue) Less(i, j int) bool {
	return v[i][0][1] < v[j][0][1]
}

func loggingMiddleWare(next http.Handler) http.Handler {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		} else {
			var method string
			if(r.Method == "DELETE") {
				method = red(r.Method)
			} else {
				method = green(r.Method)
			}
			log.Println(method, "-", r.RequestURI)
		}
		next.ServeHTTP(w, r)
	})
}

func validateData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m map[string][]interface{}

	body, readErr := ioutil.ReadAll(r.Body)
	str := string(body)
	if readErr != nil {
		panic(readErr)
	}

	unmarshalErr := json.Unmarshal([]byte(str), &m)
	if unmarshalErr != nil {
		panic(unmarshalErr)
	}

	dataSlice := make([][][][]string, 0)

	// Map datasets into two arrays
	for _, v := range m {
		rType := reflect.TypeOf(v)
		if rType.Kind() == reflect.Slice {
			dataSet := make([][][]string, 0)
			tmpSlice := make([][]string, 0)
			for _, vsub := range v {
				vmap := vsub.(map[string]interface{})
				for k, val := range vmap {
					var strVal string
					if reflect.TypeOf(val).Kind() == reflect.Float64 {
						strVal = strconv.FormatFloat(vmap[k].(float64), 'f', 6, 64)
					} else {
						strVal = vmap[k].(string)
					}
					pair := []string{k, strVal}
					tmpSlice = append(tmpSlice, pair)
				}
				dataSet = append(dataSet, tmpSlice)
				tmpSlice = make([][]string, 0) // reset tmp arr
			}
			dataSlice = append(dataSlice, dataSet)
			dataSet = make([][][]string, 0) // reset dataset tmp arr
		}
	}

	// Sort datasets
	for idx := range dataSlice[0] {
		sort.Sort(byKey(dataSlice[0][idx]))
	}

	sort.Sort(byFirstValue(dataSlice[0]))

	for idx := range dataSlice[1] {
		sort.Sort(byKey(dataSlice[1][idx]))
	}

	sort.Sort(byFirstValue(dataSlice[1]))

	// Create hashes
	dataSetOne := sha256.New()
	dataSetTwo := sha256.New()

	dataSetOne.Write([]byte(fmt.Sprintf("%b", dataSlice[0])))
	dataSetTwo.Write([]byte(fmt.Sprintf("%b", dataSlice[1])))

	hashOne := fmt.Sprintf("%x", dataSetOne.Sum(nil))
	hashTwo := fmt.Sprintf("%x", dataSetTwo.Sum(nil))

	// Compare and validate
	isValid := validData{
		Valid:       hashOne == hashTwo,
		DataHashOne: hashOne,
		DataHashTwo: hashTwo,
	}
	json.NewEncoder(w).Encode(isValid)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/validate", validateData).Methods("POST")

	r.Use(loggingMiddleWare)

	color.Set(color.FgCyan)
	log.Println("Now running validation-svc-go on", port)
	color.Unset()

	log.Fatal(http.ListenAndServe(port, r))
}