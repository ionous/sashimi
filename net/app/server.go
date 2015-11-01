package app

import (
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/net/session"
	"log"
	"net/http"
)

//
// Create the command based json-server.
//
// Some example uris:
// 	POST /game/new, create new session
// 	POST /game/<session>, send new input
// 	 GET /game/<session>/rooms/<name>/contains, list of objects
// 	 GET /game/<session>/classes/rooms/actions
//
func NewGameHandler(sessions *session.Sessions) http.HandlerFunc {
	rootGame := GameResource(sessions)

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling", r.URL.Path, r.Method)
		if res, err := resource.FindResource(rootGame, r.URL.Path[1:]); err != nil {
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
