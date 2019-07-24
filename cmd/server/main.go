package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/ejfitzgerald/annals"
)

var entriesLock = sync.Mutex{}
var entries = []annals.CompilationMetadata{}

func main() {
	http.HandleFunc("/compilation", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var metadata annals.CompilationMetadata

			data, err := ioutil.ReadAll(r.Body)

			if err != nil {
				fmt.Println("Error: Unable to read compilation metadata submission: ", err)
				return
			}

			err = json.Unmarshal(data, &metadata)

			if err != nil {
				fmt.Println("Error: Unable to decode metadata submission: ", err)
				return
			}

			entriesLock.Lock()
			defer entriesLock.Unlock()

			// add the entry to the list
			entries = append(entries, metadata)

		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			entriesLock.Lock()
			defer entriesLock.Unlock()

			data, err := json.Marshal(entries)
			if err != nil {
				fmt.Println("Failed to marshal entries data")
				return
			}

			// clear the entries
			entries = []annals.CompilationMetadata{}

			w.Write(data)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	http.ListenAndServe(":9100", nil)
}
