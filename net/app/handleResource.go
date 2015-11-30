package app

import (
	"github.com/ionous/sashimi/net/resource"
	"log"
	"net/http"
)

// HandleResource turns an IResource into an http.HandlerFunc.
func HandleResource(root resource.IResource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling", r.URL.Path, r.Method)
		if res, err := resource.FindResource(root, r.URL.Path[1:]); err != nil {
			log.Println(err)
			http.NotFound(w, r)
		} else if r.Method != "GET" && r.Method != "POST" {
			http.Error(w, r.Method, http.StatusMethodNotAllowed)
		} else {
			if r.Method == "GET" {
				Encode(w, r, res.Query())
			} else if doc, e := res.Post(r.Body); e != nil {
				log.Println(e.Error())
				http.Error(w, e.Error(), http.StatusInternalServerError)
			} else {
				Encode(w, r, doc)
			}
		}
	}
}
