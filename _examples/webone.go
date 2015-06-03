package main

import (
	"flag"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/web/simple"
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
	} else {
		fmt.Println("serving", "http://localhost:8080/")
		simple.NewSimpleServer(":8080").ListenAndServe()
	}
}
