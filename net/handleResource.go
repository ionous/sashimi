package net

import (
	"encoding/json"
	"fmt"
	"github.com/ionous/sashimi/net/resource"
	"log"
	"net/http"
)

// HandleResource turns an IResource into an http.HandlerFunc;
// providing responses to http get and post requests.
func HandleResource(root resource.IResource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling", r.URL.Path, r.Method)
		if e := handleResponse(w, r, root); e != nil {
			log.Println(e)
		}
	}
}

// NOTE: the error, if any, is automatically passed to http.Error
func handleResponse(w http.ResponseWriter, r *http.Request, root resource.IResource) (err error) {
	if res, e := resource.FindResource(root, r.URL.Path[1:]); e != nil {
		http.NotFound(w, r)
		err = e
	} else if r.Method != "GET" && r.Method != "POST" {
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
		err = fmt.Errorf("method %s not allowed", r.Method)
	} else {
		if r.Method == "GET" {
			doc := res.Query()
			encode(w, r, doc)
		} else {
			if doc, e := res.Post(r.Body); e != nil {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				err = e
			} else {
				encode(w, r, doc)
			}
		}
	}
	return
}

// Encode the passed resource to the http writer.
func encode(w http.ResponseWriter, r *http.Request, doc resource.Document) {
	w.Header().Set("Content-Type", "application/json")
	prettyBytes, _ := json.Marshal(doc)
	log.Println("returning", string(prettyBytes))
	if e := json.NewEncoder(w).Encode(doc); e != nil {
		log.Println(e)
	}
}
