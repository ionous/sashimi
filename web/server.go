package web

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/web/commands"
	"github.com/ionous/sashimi/web/session"
	"github.com/ionous/sashimi/web/simple"
	"github.com/ionous/sashimi/web/support"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	sessions map[string]session.Sessions
	*http.Server
}

// FIX: it's very silly to have to init and compile each time.
func NewGameModel() (model *M.Model, err error) {
	return S.InitScripts().Compile(ioutil.Discard /*os.Stderr*/)
}

// FIX: its not very useful to have this in the "server" code
// might as well simplify and move directly into the examples.
func NewServer(addr string, root string, dirs ...support.FilePair) *Server {
	handler := support.NewServeMux()
	sessions := map[string]session.Sessions{
		"text": session.NewSessions("text/html",
			func(id string) (ret session.ISession, err error) {
				if model, e := NewGameModel(); e != nil {
					return ret, e
				} else {
					return simple.NewSimpleSession(model)
				}
			}),
		"game": session.NewSessions("application/json",
			func(id string) (ret session.ISession, err error) {
				if model, e := NewGameModel(); e != nil {
					return ret, e
				} else {
					return commands.NewSession(id, model)
				}
			}),
	}
	handler.HandleFilePatterns(root, dirs)

	this := &Server{
		sessions,
		&http.Server{
			Addr:           addr,
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}

	handler.HandleFunc("/text/", this.handleGame)
	handler.HandleFunc("/game/", this.handleGame)
	return this
}

//
//
//
func (this *Server) HandleText(pattern, text string) {
	handler := this.Handler.(support.ServeMux)
	handler.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != pattern {
			http.NotFound(w, r)
		} else {
			w.Write([]byte(text))
		}
	})
}

//
//
//
func (this *Server) HandleFilePatterns(root string, dirs ...support.FilePair) {
	handler := this.Handler.(support.ServeMux)
	handler.HandleFilePatterns(root, dirs)
}

//
//
//
func (this *Server) handleGame(w http.ResponseWriter, r *http.Request) {
	//scheme://userinfo@host/path?query#fragment
	// 1: to skip past first slash
	log.Println("received", r.URL)
	parts := strings.Split(r.URL.Path[1:], "/")

	if cnt := len(parts); cnt > 1 {
		path := parts[0] // ex. game or simple

		if sessions, ok := this.sessions[path]; !ok {
			http.NotFound(w, r)
		} else {
			switch cnt {
			case 2:
				if r.Method != "POST" {
					http.Error(w, "unsupported get", http.StatusMethodNotAllowed)
				} else {
					if act := parts[1]; act != "new" {
						http.NotFound(w, r)
					} else {
						log.Println("creating session...")
						if id, e := sessions.NewSession(); e != nil {
							log.Println("session creation error", e)
							http.Error(w, e.Error(), http.StatusInternalServerError)
						} else {
							dest := fmt.Sprintf("/%s/%s/run?q=start", path, id)
							log.Println("redirecting to", dest)
							http.Redirect(w, r, dest, http.StatusSeeOther)
						}
					}
				}
			case 3:
				if sess, ok := sessions.Session(parts[1]); !ok {
					http.NotFound(w, r)
				} else {
					switch act := parts[2]; act {
					case "run":
						q := r.FormValue("q")
						log.Println("running", q)
						sess.Handle(q, w)
					default:
						http.NotFound(w, r)
					}
				}
			default:
				http.NotFound(w, r)
			}
		}
	}
}
