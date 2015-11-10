package tests

import (
	"fmt"
	"github.com/ionous/sashimi/compiler"
	C "github.com/ionous/sashimi/console"
	"github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/api"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/util/ident"
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

func NewTestGameSource(t *testing.T, s *Script, source string) (ret TestGame, err error) {
	if model, e := s.Compile(Log(t)); e != nil {
		err = e
	} else {
		cons := TestOutput{t, &C.BufferedOutput{}}
		cfg := R.RuntimeConfig{Calls: model.Calls, Output: cons}
		if game, e := cfg.NewGame(model.Model); e != nil {
			err = e
		} else if parser, e := R.NewObjectParser(game.ModelApi, ident.MakeId(source)); e != nil {
			err = e
		} else {
			ret = TestGame{t, game, model, cons, parser}
		}
	}
	return ret, err
}

//
func NewTestGame(t *testing.T, s *Script) (ret TestGame, err error) {
	return NewTestGameSource(t, s, "player")
}

type TestGame struct {
	t *testing.T
	*R.Game
	compiler.MemoryResult
	out    TestOutput
	Parser parser.P
}

//
// For testing:
//
func (test *TestGame) RunInput(s string) (ret []string, err error) {
	if e := test.ProcessEvents(); e != nil {
		err = e
	} else {
		in := parser.NormalizeInput(s)
		if p, m, e := test.Parser.ParseInput(in); e != nil {
			test.out.Log(fmt.Sprintf("RunInput: failed parse: %v orig: '%s' in: '%s' e: '%s'", p, s, in, e))
			err = e
		} else if act, objs, e := m.(*R.ObjectMatcher).GetMatch(); e != nil {
			test.out.Log(fmt.Sprint("RunInput: no match: ", s, e))
			err = e
		} else {
			test.QueueAction(act, objs)
			// the standard rules send an "ending the turn", we do not have to.
			ret, err = test.FlushOutput()
		}
	}
	return
}

func (test *TestGame) FlushOutput() (ret []string, err error) {
	if e := test.ProcessEvents(); e != nil {
		test.out.Println(e)
		err = e
	} else {
		ret = test.out.Flush()
	}
	return
}
