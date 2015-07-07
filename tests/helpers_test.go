package tests

import (
	C "github.com/ionous/sashimi/console"
	P "github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
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
	name := whose.Name()
	for _, l := range lines {
		out.Println(name, ":", l)
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
		} else if parser, e := R.NewParser(game); e != nil {
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
	*R.ObjectParser
}

//
// For testing:
//
func (test *TestGame) RunInput(s string) *TestGame {
	if e := test.ProcessEvents(); e != nil {
		test.out.Println(e)
	} else if _, e := test.Parse(P.NormalizeInput(s)); e != nil {
		test.out.Println(e)
	}
	return test
}

//
// For testing:
//
func (test *TestGame) RunTest(input []string) *TestGame {
	for _, in := range input {
		test.RunInput(in)
	}
	return test
}

func (test *TestGame) FlushOutput() []string {
	if e := test.ProcessEvents(); e != nil {
		test.out.Println(e)
	}
	return test.out.Flush()
}
