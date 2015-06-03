package simple

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/web/session"
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
func (this *SimpleSession) Write(in interface{}) session.ISession {
	if this.lastError == nil {
		if this.game.IsQuit() {
			this.lastError = session.SessionClosed{"player quit"}
		} else if this.game.IsFinished() {
			this.lastError = session.SessionClosed{"game finished"}
		} else if str, ok := in.(string); !ok {
			this.lastError = fmt.Errorf("unknown input %v(%T); expected string.", in, in)
		} else {
			this.game.Input(str)
		}
	}
	return this
}

//
// ISession implemenation: sends lines of text to the response handler.
//
func (this *SimpleSession) Read() (ret interface{}, err error) {
	if e := this.lastError; e != nil {
		err, this.lastError = e, nil
	} else {
		ret = this.bufferedOutput.Flush()
	}
	return ret, err
}
