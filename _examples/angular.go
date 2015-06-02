package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/web"
	"github.com/ionous/sashimi/web/commands"
	"github.com/ionous/sashimi/web/session"
	"github.com/ionous/sashimi/web/support"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//
// Run a story in an angular web-app friendly way.
// See also: https://github.com/ionous/sashimi-testapp
//
func main() {
	story := flag.String("story", "", "select the story to play.")
	flag.Parse()

	root := flag.Arg(0)
	if !stories.Stories.Select(*story) || root == "" {
		fmt.Println("Expected a story and a directory of angular files to serve.")
		fmt.Println("ex. angular -story sushi .")
		fmt.Println("Please select one of the following example stories:")
		for _, nick := range stories.Stories.List() {
			fmt.Println(" ", nick)
		}
	} else {
		fmt.Println("listening on http://localhost:8080")
		server := NewServer(
			":8080", root,
			support.Dir("/app/"),
			support.Dir("/bower_components/"),
			support.Dir("/media/"))
		log.Fatal(server.ListenAndServe())
	}
}

//
// Create the server.
//
func NewServer(addr string, root string, dirs ...support.FilePair) *http.Server {
	// FIX: probably need to implement the ISession here,
	// so that we can have our own data associated with it.
	sessions := session.NewSessions("application/json",
		func(id string) (ret session.ISession, err error) {
			if m, e := web.NewGameModel(); e != nil {
				return ret, e
			} else {
				return commands.NewSession(id, m)
			}
		})
	//
	handler := support.NewServeMux()

	// file serving
	handler.HandleFilePatterns(root, dirs)

	// game session command:
	handler.HandleFunc("/game/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling", r.URL.Path)
		if res, e := ParseResourceUrl(r.URL); e != nil {
			http.NotFound(w, r)
		} else {
			// log.Printf("received request: gid:'%s', type:'%s', rid:'%s'",
			// 	res.gameId, res.resourceType, res.resourceId)
			if e := handle(w, r, res, sessions); e != nil {
				log.Println(e)
				http.Error(w, e.Error(), http.StatusNotImplemented)
			}
		}
	})
	return &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func handle(w http.ResponseWriter, r *http.Request, res Resource, sessions session.Sessions) (err error) {
	if res.resourceType == "" && res.gameId == "new" {
		if id, e := sessions.NewSession(); e != nil {
			err = e
		} else if sess, ok := sessions.Session(id); !ok {
			err = fmt.Errorf("session not found after creation'%s'", id)
		} else {
			log.Println("created session", id)
			sess.Handle("start", w)
		}
	} else if ses, ok := res.Session(sessions); !ok {
		err = fmt.Errorf("session not found %s", res.gameId)
	} else {
		game := ses.Session().(*commands.CommandSession).Game()
		switch res.resourceType {
		case "class":
			// FIX: although this loop isnt a lot of work in and of itself
			// i feel like it should be able to be generic, via database or Json markup of the xlass
			classes := []commands.Dict{}
			for k, c := range game.Model.Classes {
				//rel, _ := inner.FindRelation(model.Relations)
				relations := []*M.RelativeProperty{}
				for _, v := range c.Properties() {
					if rel, ok := v.(*M.RelativeProperty); ok {
						relations = append(relations, rel)
					}
				}
				// slow:
				actions := []string{}
				for _, act := range game.Model.Actions {
					if act.Source().Id() == "Actors" && act.Target() == c {
						actions = append(actions, act.Name())
					}
				}
				d := commands.Dict{
					"id":        k.String(),
					"name":      c.Name(),
					"relations": relations,
					"actions":   actions,
				}
				classes = append(classes, d)
			}
			err = json.NewEncoder(w).Encode(classes)
		case "action":
			err = json.NewEncoder(w).Encode([]commands.Dict{})
		case "relation":
			err = json.NewEncoder(w).Encode([]commands.Dict{})
		case "":
			type Input struct {
				Input string `json:"input"`
			}
			v := Input{}
			decoder := json.NewDecoder(r.Body)
			if e := decoder.Decode(&v); e != nil {
				err = e
			} else {
				log.Println("running", v)
				ses.Handle(v.Input, w)
			}
		default:
			err = fmt.Errorf("unknown resource %s", res.resourceType)
		}
	}

	return err
}

// scheme://userinfo@host/path?query#fragment
type Resource struct {
	gameId, resourceType, resourceId string
}

func ParseResourceUrl(url *url.URL) (ret Resource, err error) {
	parts := strings.Split(url.Path[1:], "/")
	switch len(parts) {
	case 4:
		ret.resourceId = parts[3]
		fallthrough
	case 3:
		ret.resourceType = parts[2]
		fallthrough
	case 2:
		ret.gameId = parts[1]
	default:
		err = fmt.Errorf("error: too many parts %v", parts)
	}
	return
}

func (this *Resource) Session(sessions session.Sessions) (*session.Session, bool) {
	return sessions.Session(this.gameId)
}
