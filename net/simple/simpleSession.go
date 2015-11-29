package simple

import (
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/metal"
	"github.com/ionous/sashimi/net/session"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/standard"
)

func NewSimpleSession(calls api.LookupCallbacks, model *M.Model) (ret *SimpleSession, err error) {
	out := &SimpleOutput{}
	cfg := R.NewConfig().SetCalls(calls).SetOutput(out).SetParentLookup(standard.ParentLookup{})
	modelApi := metal.NewMetal(model, make(metal.ObjectValueMap))
	if game, e := cfg.NewGame(modelApi); e != nil {
		err = e
	} else if game, e := standard.NewStandardGame(game); e != nil {
		err = e
	} else if game, e := game.Start(); e != nil {
		err = e
	} else {
		ret = &SimpleSession{game, out, out.Flush()}

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
func (s *SimpleSession) HandleInput(in string) (err error) {
	if s.game.IsQuit() {
		err = session.SessionClosed{"player quit"}
	} else if s.game.IsComplete() {
		err = session.SessionClosed{"game finished"}
	} else {
		s.game.Input(in)
		newLines := s.out.Flush()
		s.lines = append(s.lines, newLines...)
	}
	return err
}
