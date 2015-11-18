package simple

import (
	"fmt"
	"github.com/ionous/sashimi/compiler/call"
	"github.com/ionous/sashimi/net/session"
	"github.com/ionous/sashimi/net/support"
	"github.com/ionous/sashimi/script"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func NewSimpleServer() *http.ServeMux {
	calls := call.MakeMemoryStorage()
	sessions := session.NewSessions(
		func(string) (ret session.SessionData, err error) {
			// FIX: it's very silly to have to init and compile each time.
			if model, e := script.InitScripts().CompileCalls(ioutil.Discard, calls); e != nil {
				err = e
			} else {
				model.PrintModel(func(a ...interface{}) { fmt.Println(a...) })
				ret, err = NewSimpleSession(calls, model)
			}
			return ret, err
		})

	handler := support.NewServeMux()
	handler.HandleText("/", index)
	handler.HandleFunc("/game/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("received", r.URL.Path, r.Method)
		parts := strings.Split(r.URL.Path[1:], "/")
		//
		if cnt := len(parts); cnt != 2 {
			http.NotFound(w, r)
		} else if part := parts[1]; part == "new" {
			if r.Method != "POST" {
				log.Println("Method not allowed:", r.Method)
				http.Error(w, r.Method, http.StatusMethodNotAllowed)
			} else if id, _, e := sessions.NewSession(); e != nil {
				log.Println("Failed to create session:", e)
				http.Error(w, e.Error(), http.StatusInternalServerError)
			} else {
				// redirect so that the client can see the session id.
				dest := fmt.Sprintf("/game/%s?q=start", id)
				http.Redirect(w, r, dest, http.StatusSeeOther)
			}
		} else if sd, ok := sessions.Session(part); !ok {
			log.Println("Couldnt find session:", part)
			http.NotFound(w, r)
		} else {
			sd := sd.(*SimpleSession)
			// post input
			if r.Method == "POST" {
				if e := sd.HandleInput(r.FormValue("q")); e != nil {
					log.Println("Couldn't handle input:", e)
					http.Error(w, e.Error(), http.StatusInternalServerError)
					return
				}
			}
			// return current output
			if r.Method != "POST" && r.Method != "GET" {
				log.Println("Unknown method:", r.Method)
				http.Error(w, r.Method, http.StatusMethodNotAllowed)
			} else {
				w.Header().Set("Content-Type", "text/html")
				if e := page.ExecuteTemplate(w, "simple.html", sd.lines); e != nil {
					log.Println(e)
				}
			}
		}
	})
	return handler.ServeMux
}
