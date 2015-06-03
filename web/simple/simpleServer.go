package simple

import (
	"fmt"
	S "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/web/session"
	"github.com/ionous/sashimi/web/support"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	sessions session.Sessions
	*http.Server
}

func NewSimpleServer(addr string) *Server {
	handler := support.NewServeMux()
	sessions := session.NewSessions(
		func(id string) (ret session.ISession, err error) {
			// FIX: it's very silly to have to init and compile each time.
			if model, e := S.InitScripts().Compile(ioutil.Discard); e != nil {
				err = e
			} else {
				ret, err = NewSimpleSession(model)
			}
			return ret, err
		})
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
	handler.HandleText("/", index)
	handler.HandleFunc("/game/", this.handleGame)
	return this
}

//
// POST: game/new
// POST: game/(session)?q=(command)
//
func (this *Server) handleGame(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path[1:], "/")
	log.Println("Received", parts)
	if cnt := len(parts); cnt != 2 {
		http.NotFound(w, r)
	} else if part := parts[1]; part == "new" {
		log.Println("creating session...")
		if id, e := this.sessions.NewSession(); e != nil {
			log.Println("session creation error", e)
			http.Error(w, e.Error(), http.StatusInternalServerError)
		} else {
			dest := fmt.Sprintf("/game/%s?q=start", id)
			log.Println("redirecting to", dest)
			http.Redirect(w, r, dest, http.StatusSeeOther)
		}
	} else if sess, ok := this.sessions.Session(part); !ok {
		http.NotFound(w, r)
	} else {
		q := r.FormValue("q")
		if q != "start" && r.Method != "POST" {
			http.Error(w, "unsupported get", http.StatusMethodNotAllowed)
		} else if lines, e := sess.WriteRead(q); e != nil {
			log.Println("!!! Error", e)
			http.Error(w, e.Error(), http.StatusInternalServerError)
		} else {
			log.Println("here", lines)
			w.Header().Set("Content-Type", "text/html")
			if e := page.ExecuteTemplate(w, "simple.html", lines); e != nil {
				log.Println("!!! Error", e)
			}
		}
	}
}
