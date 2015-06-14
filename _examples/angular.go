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
// Run a story in an angular web-app friendly way.
// See also: https://github.com/ionous/sashimi-testapp
//
func main() {
	story := flag.String("story", "", "select the story to play.")
	flag.Parse()

	root := flag.Arg(0)
	if !stories.Select(*story) || root == "" {
		fmt.Println("Expected a story and a directory of angular files to serve.")
		fmt.Println("ex. angular -story sushi .")
		fmt.Println("Please select one of the following example stories:")
		for _, nick := range stories.List() {
			fmt.Println(" ", nick)
		}
	} else {
		fmt.Println("listening on http://localhost:8080")
		handler := app.NewServer()
		handler.HandleFilePatterns(root,
			support.Dir("/app/"),
			support.Dir("/bower_components/"),
			support.Dir("/media/"))
		log.Fatal(http.ListenAndServe(":8080", handler))
	}
}
