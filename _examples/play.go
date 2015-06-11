package main

import (
	"flag"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/standard"
)

//
// Boilerplate to run the story.
//
func main() {
	story := flag.String("story", "", "select the story to play.")
	options := standard.ParseCommandLine()

	if !stories.Select(*story) {
		flag.PrintDefaults()
		fmt.Println("Please select one of the following stories:")
		for _, nick := range stories.List() {
			fmt.Println(" ", nick)
		}
	} else {
		standard.RunGame(options)
	}
}
