package web

import (
	"encoding/json"
	"fmt"
	"github.com/ionous/sashimi/web/support"
	"html"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	*http.Server
	sessions Sessions
}

func NewServer(addr string, root string, dirs ...support.FilePair) *Server {
	handler := support.NewServeMux()

	this := &Server{
		&http.Server{
			Addr:           addr,
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}, NewSessions(),
	}
	handler.HandleFunc("/game/", this.handleGame)
	//handler.HandleFunc("/game/data.json", handleData)
	//handler.HandleFunc("/action/", handleAction)
	return this
}

func (this *Server) HandleText(pattern, text string) {
	handler := this.Handler.(support.ServeMux)
	handler.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(text))
	})
}

func (this *Server) HandleFilePatterns(root string, dirs ...support.FilePair) {
	handler := this.Handler.(support.ServeMux)
	handler.HandleFilePatterns(root, dirs)
}

var simplest = `<html>
<body><form action="input" id="f"><input id="q" name="q"></form></body>
</html>`

func (this *Server) handleGame(w http.ResponseWriter, r *http.Request) {
	//scheme://userinfo@host/path?query#fragment
	// 1: to skip past first slash
	log.Println("received", r.URL)
	parts := strings.Split(r.URL.Path[1:], "/")
	switch cnt := len(parts); cnt {
	case 2:
		if r.Method != "POST" {
			http.Error(w, "unsupported get", http.StatusMethodNotAllowed)
		} else {
			if act := parts[1]; act != "new" {
				http.NotFound(w, r)
			} else {
				log.Println("creating session...")
				if sess, e := this.sessions.NewSession(); e != nil {
					log.Println("session creation error", e)
					http.Error(w, e.Error(), http.StatusInternalServerError)
				} else {
					dest := fmt.Sprintf("/game/%s/run?q=start", sess.id)
					log.Println("redirecting to", dest)
					http.Redirect(w, r, dest, http.StatusSeeOther)
				}
			}
		}
	case 3:
		if sess, ok := this.sessions.Session(parts[1]); !ok {
			http.NotFound(w, r)
		} else {
			switch act := parts[2]; act {
			case "run":
				q := r.FormValue("q")
				log.Println("running", q)
				sess.Handle(w, q)
			default:
				http.NotFound(w, r)
			}
		}
	default:
		http.NotFound(w, r)
	}
}

// okay: index.html shol

//
// copied from boilerplate :(
//

//
func handleAction(w http.ResponseWriter, r *http.Request) {
	s := r.FormValue("obj")
	o := fmt.Sprintf("Hello, %s %s %q", r.Method, s, html.EscapeString(r.URL.Path))
	fmt.Fprint(w, o)
	fmt.Println(o)
	switch r.Method {
	case "GET":
		// Serve the resource.
	case "POST":
		// Create a new record.
	case "PUT":
		// Update an existing record.
	case "DELETE":
		// Remove the record.
	default:
		// Give an error message.
	}
}

type Object struct {
	Id       string   `json:"id"`
	States   []string `json:"states,omitempty"`
	Contains []Object `json:"contains,omitempty"`
}

type Array []interface{}
type Dict map[string]interface{}

func handleData(w http.ResponseWriter, r *http.Request) {
	studio := Array{
		Dict{"present": []Object{
			Object{
				Id: "Studio",
				Contains: []Object{
					Object{
						Id:     "Cabinet",
						States: []string{"Closed"},
					},
					Object{
						Id: "Easel",
					},
					Object{
						Id: "Table",
						Contains: []Object{
							Object{
								Id: "Vase",
							},
						},
					},
					Object{
						Id: "Aquarium",
						Contains: []Object{
							Object{
								Id: "Gravel",
							},
							Object{
								Id: "Seaweed",
							},
							Object{
								Id: "EvilFish",
							},
						},
					},
				},
			},
		}}}
	//	encoder:Encoder writes JSON objects to an output stream.
	// http://thenewstack.io/make-a-restful-json-api-go/

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(studio)
}
