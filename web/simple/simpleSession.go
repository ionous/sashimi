package simple

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/web/session"
	"io"
	"log"
)

func NewSimpleSession(model *M.Model) (ret *SimpleSession, err error) {
	out := &SimpleOutput{}
	if game, e := standard.NewStandardGame(model, out); e != nil {
		err = e
	} else if game, e := game.Start(); e != nil {
		err = e
	} else {
		ret = &SimpleSession{&game, out, nil}
	}
	return ret, err
}

//
// a single game run by the server
//
type SimpleSession struct {
	game           *standard.StandardGame
	bufferedOutput *SimpleOutput
	lastError      error
}

//
// ISession implementation
//
func (this *SimpleSession) Read(in string) session.ISession {
	if this.lastError == nil {
		if this.game.IsQuit() {
			this.lastError = session.SessionClosed{"player quit"}
		} else if this.game.IsFinished() {
			this.lastError = session.SessionClosed{"game finished"}
		} else {
			this.game.Input(in)
		}
	}
	return this
}

func (this *SimpleSession) Write(w io.Writer) (err error) {
	if e := this.lastError; e != nil {
		err, this.lastError = e, nil
	} else {
		lines := this.bufferedOutput.Flush()
		log.Println("here", lines)
		if e := page.ExecuteTemplate(w, "simple.html", lines); e != nil {
			err = e
		}
	}
	return err
}
