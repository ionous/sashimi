package main

import (
	"flag"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/web"
	"log"
)

func main() {
	story := flag.String("story", "", "select the story to play.")
	flag.Parse()
	if !stories.Stories.Select(*story) {
		flag.PrintDefaults()
		fmt.Println("Please select one of the following stories:")
		for _, nick := range stories.Stories.List() {
			fmt.Println(" ", nick)
		}
	} else if root := flag.Arg(0); root == "" {
		fmt.Println("expected a directory of files to serve")
	} else {
		server := web.NewServer(
			":8080", root,
			support.Dir("/app/"),
			support.Dir("/bower_components/"),
			support.Dir("/media/"))
		log.Fatal(server.ListenAndServe())
	}
}
