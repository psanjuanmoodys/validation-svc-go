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
	// fmt.Println("===body ", str)
	
	unmarshalErr := json.Unmarshal([]byte(str), &m)
	if unmarshalErr != nil {
		panic(unmarshalErr)
	}
	
	// fmt.Println("===body ", m)
	// fmt.Println("===body ", m["dataSetOne"])

	for k, v := range m {
		rType := reflect.TypeOf(v)
		if rType.Kind() == reflect.Slice {
			fmt.Println(v, ": ")
			for ksub, vsub := range v {
				fmt.Println(reflect.TypeOf(vsub).Kind())
				fmt.Println("k: ", ksub, " v: ", vsub)
			}
		} else {
			fmt.Println("k: ", k, " v: ", v)
		}
	}
	
	isValid := ValidData{Valid: true};
	json.NewEncoder(w).Encode(isValid)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/validate", validateData).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}