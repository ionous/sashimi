package app

import (
	"encoding/json"
	"github.com/ionous/sashimi/net/resource"
	"log"
	"net/http"
)

// where is a good place for this?
func Encode(w http.ResponseWriter, r *http.Request, doc resource.Document) {
	w.Header().Set("Content-Type", "application/json")
	prettyBytes, _ := json.Marshal(doc)
	log.Println("returning", string(prettyBytes))
	if e := json.NewEncoder(w).Encode(doc); e != nil {
		log.Println(e)
	}
}
