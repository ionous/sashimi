package tests

import (
	"fmt"
	"github.com/ionous/sashimi/compiler"
	C "github.com/ionous/sashimi/console"
	"github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/api"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard" // :(
	"strings"
	"testing"
)

var _ = fmt.Print

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

func (out TestOutput) ActorSays(whose api.Instance, lines []string) {
	prop, _ := whose.GetProperty("Name")
	name := prop.GetValue().GetText()

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
		cfg := R.RuntimeConfig{Calls: model.Calls, Output: cons}
		if game, e := cfg.NewGame(model.Model); e != nil {
			err = e
		} else if parser, e := standard.NewStandardParser(game); e != nil {
			err = e
		} else {
			ret = TestGame{t, game, model, cons, parser}
		}
	}
	return ret, err
}

type TestGame struct {
	t *testing.T
	*R.Game
	compiler.MemoryResult
	out            TestOutput
	StandardParser *standard.StandardParser
}

//
// For testing:
//
func (test *TestGame) RunInput(s string) (err error) {
	if e := test.ProcessEvents(); e != nil {
		err = e
	} else {
		in := parser.NormalizeInput(s)
		if m, e := test.StandardParser.ParseInput(in); e != nil {
			test.out.Log(fmt.Sprintf("RunInput: failed parse: %v orig: '%s' in: '%s' e: '%s'", m, s, in, e))
			err = e
		} else if e := m.OnMatch(); e != nil {
			test.out.Log(fmt.Sprint("RunInput: no match: ", s, e))
			err = e
		}
	}
	return err
}

func (test *TestGame) FlushOutput() []string {
	if e := test.ProcessEvents(); e != nil {
		test.out.Println(e)
	}
	return test.out.Flush()
}
