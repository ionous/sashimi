package main

import (
	"flag"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/standard/output"
)

// Boilerplate to run the story.
func main() {
	story := flag.String("story", "", "select the story to play.")
	options := output.ParseCommandLine()

	if s, ok := stories.Select(*story); !ok {
		flag.PrintDefaults()
		fmt.Println("Please select one of the following stories:")
		for _, nick := range stories.List() {
			fmt.Println(" ", nick)
		}
	} else {
		output.RunGame(s, options)
	}
}
