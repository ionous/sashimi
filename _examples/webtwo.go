package main

import (
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/web"
	"log"
)

var index = `
<!DOCTYPE html>
<html lang="en">
<body>
<h1>New Game</h1>
    <div id="input">
        <form action="/game/new" method="POST">
            <button>Start</button>
        </form>
    </div>
</body>
</html>`

func main() {
	script.AddScript(stories.An_Empty_Room)
	//http://localhost:8080/game/new
	server := web.NewServer(":8080", "")
	server.HandleText("/", index)
	log.Println("serving", "http://localhost:8080/")
	log.Fatal(server.ListenAndServe())
}
