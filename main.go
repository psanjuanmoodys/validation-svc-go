package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	// "crypto/sha256"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"reflect"
)

type ValidData struct {
	Valid bool `json:"valid"`
}

func validateData (w http.ResponseWriter, r *http.Request) {
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

	dataSlice := make([][][]string, 0)

	for _, v := range m {
		rType := reflect.TypeOf(v)
		if rType.Kind() == reflect.Slice {
			dataSet := make([][]string, 0)
			tmpSlice := make([]string, 0)
			for _, vsub := range v {
				// fmt.Println(reflect.TypeOf(vsub).Kind())
				vmap := vsub.(map[string]interface{})
				for k := range vmap {
					tmpSlice = append(tmpSlice, k)
				}
				dataSet = append(dataSet, tmpSlice)
				tmpSlice = make([]string, 0) // reset tmp arr
			}
			dataSlice = append(dataSlice, dataSet)
			dataSet = make([][]string, 0) // reset dataset tmp arr
		}
	}

	fmt.Println(dataSlice)
	
	isValid := ValidData{Valid: true};
	json.NewEncoder(w).Encode(isValid)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/validate", validateData).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}