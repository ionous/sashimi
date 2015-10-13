package simple

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/session"
	"github.com/ionous/sashimi/standard"
)

func NewSimpleSession(model *M.Model) (ret *SimpleSession, err error) {
	out := &SimpleOutput{}
	if game, e := standard.NewStandardGame(model, nil, out); e != nil {
		err = e
	} else {
		immediate := true
		if game, e := game.Start(immediate); e != nil {
			err = e
		} else {
			ret = &SimpleSession{game, out, out.Flush()}
		}
	}
	return ret, err
}

//
// a single game run by the server
//
type SimpleSession struct {
	game  *standard.StandardGame
	out   *SimpleOutput
	lines []string
}

//
func (this *SimpleSession) HandleInput(in string) (err error) {
	if this.game.IsQuit() {
		err = session.SessionClosed{"player quit"}
	} else if this.game.IsFinished() {
		err = session.SessionClosed{"game finished"}
	} else {
		this.game.Input(in)
		newLines := this.out.Flush()
		this.lines = append(this.lines, newLines...)
	}
	return err
}
