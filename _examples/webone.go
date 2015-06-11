package main

import (
	"flag"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/net/simple"
	"net/http"
)

func main() {
	story := flag.String("story", "", "select the story to play.")
	flag.Parse()
	if !stories.Select(*story) {
		flag.PrintDefaults()
		fmt.Println("Please select one of the following stories:")
		for _, nick := range stories.List() {
			fmt.Println(" ", nick)
		}
	} else {
		fmt.Println("serving", "http://localhost:8080/")
		http.ListenAndServe(":8080", simple.NewSimpleServer())
	}
}
