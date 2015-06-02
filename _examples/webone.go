package main

import (
	"flag"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/web"
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
		server := web.NewServer(":8080", "")
		server.HandleText("/", index)
		fmt.Println("serving", "http://localhost:8080/")
		server.ListenAndServe()
	}
}

//
var index = `
<!DOCTYPE html>
<html lang="en">
<body>
<h1>New Game</h1>
    <div id="input">
        <form action="/text/new" method="POST">
            <button>Start</button>
        </form>
    </div>
</body>
</html>`
