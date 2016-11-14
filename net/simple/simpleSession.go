package simple

import (
	"fmt"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/play"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/api"
	_ "github.com/ionous/sashimi/standard" // init
	"github.com/ionous/sashimi/standard/framework"
)

func NewSimpleSession(modelApi meta.Model, calls api.LookupCallbacks) (ret *SimpleSession, err error) {
	out := &SimpleOutput{}
	cfg := play.NewConfig().SetCalls(calls).SetOutput(out).SetParentLookup(framework.NewParentLookup(modelApi))
	game := cfg.MakeGame(modelApi)

	if game, e := framework.NewStandardGame(game); e != nil {
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
	game  *framework.StandardGame
	out   *SimpleOutput
	lines []string
}

//
func (s *SimpleSession) HandleInput(in string) (err error) {
	if s.game.IsQuit() {
		err = fmt.Errorf("session closed: player quit.")
	} else if s.game.IsComplete() {
		err = fmt.Errorf("session closed: player finished game.")
	} else {
		s.game.Input(in)
		newLines := s.out.Flush()
		s.lines = append(s.lines, newLines...)
	}
	return err
}
