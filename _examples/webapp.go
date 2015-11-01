package main

import (
	"flag"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/net/app"
	"github.com/ionous/sashimi/net/support"
	"log"
	"net/http"
)

//
// Run a story in an web-app friendly way.
// See also: https://github.com/ionous/sashimi-testapp
//
func main() {
	story := flag.String("story", "", "select the story to play.")
	flag.Parse()

	root := flag.Arg(0)
	if !stories.Select(*story) || root == "" {
		fmt.Println("Expected a story and a directory of files to serve.")
		fmt.Println("ex. webapp -story sushi .")
		fmt.Println("The following example stories are available:")
		for _, nick := range stories.List() {
			fmt.Println(" ", nick)
		}
	} else {
		fmt.Println("listening on http://localhost:8080")
		handler := support.NewServeMux()
		calls := call.MakeMemoryStorage()

		sessions := session.NewSessions(
			func(id string) (ret session.SessionData, err error) {
				// FIX: it's very silly to have to init and compile each time.
				// the reason is because relations change the original model.
				if m, e := script.InitScripts().CompileCalls(ioutil.Discard, calls); e != nil {
					err = e
				} else {
					ret, err = NewCommandSession(id, m, calls)
				}
				return ret, err
			})

		handler.HandleFunc("/game/", app.NewGameHandler(sessions))
		handler.HandleFilePatterns(root,
			support.Dir("/app/"),
			support.Dir("/bower_components/"),
			support.Dir("/media/"))
		go support.OpenBrowser("http://localhost:8080/app/")
		log.Fatal(http.ListenAndServe(":8080", handler))
	}
}
