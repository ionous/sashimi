package app

import (
	"encoding/json"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/net/session"
	"github.com/ionous/sashimi/net/support"
	"github.com/ionous/sashimi/script"
	"io/ioutil"
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
func NewServer() *support.ServeMux {
	handler := support.NewServeMux()

	// game session command:
	sessions := session.NewSessions(
		func(id string) (ret session.SessionData, err error) {
			// FIX: it's very silly to have to init and compile each time.
			// the reason is because relations change the original model.
			if m, e := script.InitScripts().Compile(ioutil.Discard); e != nil {
				err = e
			} else {
				ret, err = NewCommandSession(id, m)
			}
			return ret, err
		})
	rootGame := GameResource(sessions)

	handler.HandleFunc("/game/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling", r.URL.Path, r.Method)
		if res, err := resource.NewResourcePath(r.URL.Path[1:]).FindResource(rootGame); err != nil {
			log.Println(err)
			http.NotFound(w, r)
		} else if r.Method != "GET" && r.Method != "POST" {
			http.Error(w, r.Method, http.StatusMethodNotAllowed)
		} else {
			if r.Method == "GET" {
				encode(w, r, res.Query())
			} else if doc, e := res.Post(r.Body); e != nil {
				log.Println(e.Error())
				http.Error(w, e.Error(), http.StatusInternalServerError)
			} else {
				encode(w, r, doc)
			}
		}
	})
	return &handler
}

func encode(w http.ResponseWriter, r *http.Request, doc resource.Document) {
	w.Header().Set("Content-Type", "application/json")
	prettyBytes, _ := json.Marshal(doc)
	log.Println("returning", string(prettyBytes))
	if e := json.NewEncoder(w).Encode(doc); e != nil {
		log.Println(e)
	}
}
