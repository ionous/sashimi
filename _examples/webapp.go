package main

import (
	"flag"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/net/app"
	"github.com/ionous/sashimi/net/support"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"
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
		handler := app.NewServer()
		handler.HandleFilePatterns(root,
			support.Dir("/app/"),
			support.Dir("/bower_components/"),
			support.Dir("/media/"))
		go Open("http://localhost:8080/app/")
		log.Fatal(http.ListenAndServe(":8080", handler))
	}
}

// From Miki Tebeka: Open calls the OS default program for uri
// http://go-wise.blogspot.com/2012/04/open-fileurls-with-default-program.html
func Open(uri string) {
	var commands = map[string]string{
		"windows": "start",
		"darwin":  "open",
		"linux":   "xdg-open",
	}
	time.Sleep(300 * time.Millisecond)
	if run, ok := commands[runtime.GOOS]; !ok {
		log.Printf("don't know how to open the browser on %s platform", runtime.GOOS)
	} else {
		cmd := exec.Command(run, uri)
		if e := cmd.Start(); e != nil {
			log.Printf("error opening browser %s", e.Error())
		}
	}
}
