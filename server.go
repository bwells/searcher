package main

import (
	"encoding/json"
	"github.com/bwells/trie"
	"net/http"
	"text/template"
)

func queryHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	query := r.FormValue("term")

	results := multiMatch(search_trie, query)

	if len(results) == 0 {
		w.Write([]byte("[]"))
		return
	}

	output, err := json.Marshal(results)

	if err != nil {
		w.Write([]byte("[]"))
		return
	}

	w.Write(output)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

var search_trie *trie.Trie

var tpl *template.Template

func serve() {

	tpl, _ = template.ParseFiles("form.html")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", queryHandler)
	http.ListenAndServe(":9000", nil)
}
