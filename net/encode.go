package net

import (
	"encoding/json"
	"github.com/ionous/sashimi/net/resource"
	"log"
	"net/http"
)

// Encode the passed resource to the http writer.
func Encode(w http.ResponseWriter, r *http.Request, doc resource.Document) {
	w.Header().Set("Content-Type", "application/json")
	prettyBytes, _ := json.Marshal(doc)
	log.Println("returning", string(prettyBytes))
	if e := json.NewEncoder(w).Encode(doc); e != nil {
		log.Println(e)
	}
}
