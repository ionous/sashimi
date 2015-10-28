package tests

import (
	"github.com/ionous/sashimi/compiler/call"
	C "github.com/ionous/sashimi/console"
	P "github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard" // :(
	"strings"
	"testing"
)

type LogOutput struct {
	t *testing.T
}

func Log(t *testing.T) LogOutput {
	return LogOutput{t}
}

func (out LogOutput) Write(bytes []byte) (int, error) {
	out.t.Log(strings.TrimSpace(string(bytes)))
	return len(bytes), nil
}

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
	out.t.Log(strings.TrimSpace(s))
}

//
func NewTestGame(t *testing.T, s *Script) (ret TestGame, err error) {
	if model, e := s.Compile(Log(t)); e != nil {
		err = e
	} else {
		cons := TestOutput{t, &C.BufferedOutput{}}
		cfg := R.Config{Calls: model.Calls, Output: cons}
		if game, e := cfg.NewGame(model.Model); e != nil {
			err = e
		} else if parser, e := standard.NewParser(game); e != nil {
			err = e
		} else {
			ret = TestGame{t, game, cons, parser, model.Calls}
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
	call call.MemoryStorage
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
