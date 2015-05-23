package web

import (
	"fmt"
	"github.com/ionous/sashimi/standard"
	"html/template"
	"log"
	"net/http"
)

type Signal chan struct{}

// channel helper to serialize concurrent requests
type CallResponse struct {
	input  string
	output http.ResponseWriter
	done   Signal
}

type Session struct {
	id             string
	game           standard.StandardGame
	inOut          chan CallResponse
	bufferedOutput *SessionOutput
	started        bool
}

// starts a new game session.
func StartSession(id string, game standard.StandardStart, buffer *SessionOutput) (ret *Session, err error) {
	if game, e := game.Start(); e != nil {
		err = e
	} else {
		ret = &Session{id, game, make(chan CallResponse), buffer, false}
		// FIX: some sort of timeout to evenually kill sessions?
		go func() {
			for ret.readWrite() {
			}
		}()
	}
	return ret, err
}

// send new input to the game session.
func (this *Session) Handle(output http.ResponseWriter, input string) {
	done := make(Signal)
	this.inOut <- CallResponse{input, output, done}
	<-done
}

var simple = template.Must(template.New("simple.html").Parse(`<!DOCTYPE html>
<html lang="en">
<body onload='setFocus()'>
    <div id="story">{{ range .Lines }}
        <p>{{ . }}</p>{{ end }}
    </div>
    <div id="input">
        <form action="run" id="f" method="POST">
            <input id="q" name="q"　type="text">
        </form>
    </div>
    <script>
			function setFocus(){
			    document.getElementById("q").focus();
			}
		</script>
</body>
</html>`))

// true if game/session is still going
func (this *Session) readWrite() (okay bool) {
	cr := <-this.inOut
	in, w, done := cr.input, cr.output, cr.done
	if ok := this.handleInput(in); ok {
		type Lines struct {
			Lines []string
		}
		lines := Lines{this.bufferedOutput.Flush()}
		if e := simple.ExecuteTemplate(w, "simple.html", lines); e != nil {
			log.Println("template error", e)
			http.Error(w, e.Error(), http.StatusInternalServerError)
		}
		done <- struct{}{}
		okay = true
	}
	return okay
}

// the dance areound "started" isnt very nice
// it has to do with the fact the initial post redirects to a get
// so we want data on that first get, but no input.
func (this *Session) handleInput(in string) (okay bool) {
	start := in == "start"
	if !this.started {
		this.started = start
		fmt.Println("starting", in, "for", this.id)
		okay = true // start doesnt send input to game.
	} else if start {
		fmt.Println("ignoring", in, "for", this.id)
		okay = true // start doesnt send input to game.
	} else {
		fmt.Println("handling", in, "for", this.id)
		okay = this.game.Input(in)
	}
	return okay
}
