package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {

	query := r.FormValue("query")
	fmt.Println(query)

	query = r.URL.Query()["query"]
	fmt.Println(query)

}

func serve() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9000", nil)
}
