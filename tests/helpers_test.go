package tests

import (
	C "github.com/ionous/sashimi/console"
	P "github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard" // :(
	"os"
	"testing"
)

//
type TestOutput struct {
	t *testing.T
	*C.BufferedOutput
}

//
// Standard output.
//
func (out TestOutput) ScriptSays(lines []string) {
	for _, l := range lines {
		out.Println(l)
	}
}

func (out TestOutput) ActorSays(whose *R.GameObject, lines []string) {
	name := whose.Value("Name")
	for _, l := range lines {
		out.Println(name, ": ", l)
	}
}

func (out TestOutput) Log(s string) {
	out.t.Log(s)
}

//
func NewTestGame(t *testing.T, s *Script) (ret TestGame, err error) {
	if model, e := s.Compile(os.Stderr); e != nil {
		err = e
	} else {
		cons := TestOutput{t, &C.BufferedOutput{}}
		if game, e := R.NewGame(model, cons); e != nil {
			err = e
		} else if parser, e := standard.NewParser(game); e != nil {
			err = e
		} else {
			ret = TestGame{t, game, cons, parser}
		}
	}
	return ret, err
}

type TestGame struct {
	t *testing.T
	*R.Game
	out TestOutput
	//*R.ObjectParser
	*standard.Parser
}

//
// For testing:
//
func (test *TestGame) RunInput(s string) (err error) {
	if e := test.ProcessEvents(); e != nil {
		err = e
	} else if m, e := test.ParseInput(P.NormalizeInput(s)); e != nil {
		err = e
	} else if e := m.OnMatch(); e != nil {
		err = e
	}
	return err
}

func (test *TestGame) FlushOutput() []string {
	if e := test.ProcessEvents(); e != nil {
		test.out.Println(e)
	}
	return test.out.Flush()
}
