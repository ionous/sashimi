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
func (this TestOutput) ScriptSays(lines []string) {
	for _, l := range lines {
		this.Println(l)
	}
}

func (this TestOutput) ActorSays(whose *R.GameObject, lines []string) {
	name := whose.Name()
	for _, l := range lines {
		this.Println(name, ":", l)
	}
}

func (this TestOutput) Log(s string) {
	this.t.Log(s)
}

//
func NewTestGame(t *testing.T, s *Script) (ret TestGame, err error) {
	if model, e := s.Compile(os.Stderr); e != nil {
		err = e
	} else {
		cons := TestOutput{t, &C.BufferedOutput{}}
		if game, e := R.NewGame(model, cons); e != nil {
			err = e
		} else {
			ret = TestGame{t, game, cons}
		}
	}
	return ret, err
}

type TestGame struct {
	t *testing.T
	*R.Game
	out TestOutput
}

//
// For testing:
//
func (this *TestGame) RunInput(s string) *TestGame {
	if e := this.ProcessEvents(); e != nil {
		this.out.Println(e)
	} else if _, e := this.Parser.Parse(P.NormalizeInput(s)); e != nil {
		this.out.Println(e)
	}
	return this
}

//
// For testing:
//
func (this *TestGame) RunTest(input []string) *TestGame {
	for _, in := range input {
		this.RunInput(in)
	}
	return this
}

func (this *TestGame) FlushOutput() []string {
	if e := this.ProcessEvents(); e != nil {
		this.out.Println(e)
	}
	return this.out.Flush()
}
