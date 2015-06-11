package sashimi

import (
	C "github.com/ionous/sashimi/console"
	P "github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
	"log"
	"os"
)

//
type TestOutput struct {
	*C.BufCon
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
	os.Stderr.WriteString(s)
}

//
func NewTestGame(s *Script, input []string) (ret TestGame, err error) {
	if model, e := s.Compile(os.Stderr); e != nil {
		err = e
	} else {
		cons := TestOutput{C.NewBufCon(input)}
		if game, e := R.NewGame(model, cons); e != nil {
			err = e
		} else {
			ret = TestGame{cons, game}
		}
	}
	return ret, err
}

type TestGame struct {
	console TestOutput
	*R.Game
}

func (this TestGame) FlushOutput() []string {
	return this.console.Flush()
}

//
// For testing:
//
func (this TestGame) RunTest() []string {
	for {
		// process queue
		if e := this.ProcessEvents(); e != nil {
			log.Println(e)
		} else {
			// read input
			if s, ok := this.console.Readln(); !ok {
				break
			} else {
				in := P.NormalizeInput(s)
				if in == "q" || in == "quit" {
					break
				}
				if e := this.RunCommand(in); e != nil {
					log.Println(e)
				}
			}
		}
	}
	return this.FlushOutput()
}
