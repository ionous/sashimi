package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/web/commands"
	"github.com/ionous/sashimi/web/session"
	"github.com/ionous/sashimi/web/support"
	"io/ioutil"
	"log"
	"net/http"
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
	sessions := session.NewSessions(
		func(id string) (ret session.ISession, err error) {
			// FIX: it's very silly to have to init and compile each time.
			if m, e := S.InitScripts().Compile(ioutil.Discard); e != nil {
				err = e
			} else {
				ret, err = commands.NewSession(id, m)
			}
			return
		})
	//
	handler := support.NewServeMux()

	// file serving
	handler.HandleFilePatterns(root, dirs)

	// game session command:
	handler.HandleFunc("/game/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling", r.URL.Path)
		if req, e := ParseResourceUrl(r); e != nil {
			http.NotFound(w, r)
		} else {
			// log.Printf("received request: gid:'%s', type:'%s', rid:'%s'",
			// 	res.gameId, res.resourceType, res.resourceId)
			if data, e := handle(req, sessions); e == nil {
				w.Header().Set("Content-Type", "application/json")
				prettyBytes, _ := json.Marshal(data)
				log.Println("returning", string(prettyBytes))
				if e := json.NewEncoder(w).Encode(data); e != nil {
					log.Println(e)
				}
			} else {
				log.Println(e)
				if _, ok := e.(notFound); ok {
					http.NotFound(w, r)
				} else {
					http.Error(w, e.Error(), http.StatusNotImplemented)
				}
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

func handle(req Request, sessions session.Sessions) (ret interface{}, err error) {
	if req.resourceType == "" && req.gameId == "new" {
		if id, e := sessions.NewSession(); e != nil {
			err = e
		} else if sess, ok := sessions.Session(id); !ok {
			err = fmt.Errorf("session not found after creation'%s'", id)
		} else {
			log.Println("created session", id)
			//
			ret, err = sess.WriteRead("start")
		}
	} else if ses, ok := req.Session(sessions); !ok {
		err = NotFound("session not found %s", req)
	} else {
		game := ses.Session().(*commands.CommandSession).Game()
		switch req.resourceType {
		// FIX: although this loop isnt a lot of work in and of itself
		// i feel like it should be able to be generic, via database or Json markup of the class
		// at the very least we probably need an IResource
		case "class":
			if req.resourceId == "" {
				classes := []commands.Dict{}
				for k, c := range game.Model.Classes {
					d := commands.Dict{
						"id":   k.String(),
						"name": c.Name(),
					}
					classes = append(classes, d)
				}
				ret = classes
			} else if cls, ok := game.Model.Classes[M.StringId(req.resourceId)]; !ok {
				err = NotFound("class not found %s", req)
			} else {
				switch req.dataType {
				case "":
					// really we should return a list of urls for more restiness
					ret = []string{}
				case "relations":
					relations := []*M.RelativeProperty{}
					//rel, _ := inner.FindRelation(model.Relations)
					for _, v := range cls.AllProperties() {
						if rel, ok := v.(*M.RelativeProperty); ok {
							relations = append(relations, rel)
						}
					}
					ret = relations
				case "actions":
					// slow:
					actions := []M.StringId{}
					for _, act := range game.Model.Actions {
						if act.Source().Id() == "Actors" && (act.Target() == cls || cls.HasParent(act.Target())) {
							actions = append(actions, act.Id())
						}
					}
					ret = actions
				default:
					err = NotFound("unknown data type %s", req)
				}
			}
		// case "action":
		// 	ret = []commands.Dict{}
		// case "relation":
		// 	ret = []commands.Dict{}
		case "":
			decoder, in := json.NewDecoder(req.Body), commands.CommandInput{}
			if e := decoder.Decode(&in); e != nil {
				err = e
			} else {
				log.Println("running", in)
				ret, err = ses.WriteRead(in)
			}
		default:
			err = NotFound("unknown resource %s", req)
		}
	}

	return ret, err
}

// FIX? better would be, allow the request to return a function
// if that fails, it means the resource is not found, otherwise execute the function
// and an error means internal server error.
type notFound struct {
	reason string
}

func NotFound(format string, args ...interface{}) notFound {
	return notFound{fmt.Sprintf(format, args...)}
}

func (this notFound) Error() string {
	return this.reason
}

// /game/:gameId/class/:clsId/dataType
type Request struct {
	*http.Request
	gameId, resourceType, resourceId, dataType string
}

func ParseResourceUrl(r *http.Request) (ret Request, err error) {
	// scheme://userinfo@host/path?query#fragment
	parts := strings.Split(r.URL.Path[1:], "/")
	switch len(parts) {
	case 5:
		ret.dataType = parts[4]
		fallthrough
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
	ret.Request = r
	return
}
func (this *Request) String() string {
	return fmt.Sprint("resource: %s,%s,%s,%s", this.gameId, this.resourceType, this.resourceId, this.dataType)
}
func (this *Request) Session(sessions session.Sessions) (*session.Session, bool) {
	return sessions.Session(this.gameId)
}
