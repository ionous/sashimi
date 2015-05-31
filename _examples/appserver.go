package main

import (
	"flag"
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
	if directory := flag.Arg(0); directory == "" {
		log.Fatal("expected a directory of files to serve")
	}

	//http://localhost:8080/game/new
	AddScript(func(s *Script) {
		s.The("story",
			Called("testing"),
			Has("author", "me"),
			Has("headline", "extra extra"))
		s.The("room",
			Called("somewhere"),
			Has("description", "an empty room"),
		)
	})

	server := web.NewServer(":8080", root,
		support.Dir("/app/"),
		support.Dir("/bower_components/"),
		support.Dir("/media/"))
	log.Fatal(server.ListenAndServe())
}
	server := web.NewServer(":8080", "")
	server.HandleText("/", index)
	log.Println("serving", "http://localhost:8080/")
	log.Fatal(server.ListenAndServe())
}

//http://localhost:8080/app/
/*func main() {
	// FIX: via command line.
	root := "/Users/ionous/Dev/ngmockup/"
	
*/
