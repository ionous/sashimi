package sashimi

import (
	C "github.com/ionous/sashimi/console"
	P "github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
	"io"
	"log"
	"os"
)

//
func NewLogger() *log.Logger {
	return log.New(os.Stderr, "test:", log.Lshortfile)
}

//
func CompileGameWithConsole(s *Script, cons C.IConsole) (ret GameHelper, err error) {
	if model, e := s.Compile(os.Stderr); e != nil {
		err = e
	} else {
		type Output struct {
			C.IConsole
			io.Writer
		}
		if game, e := R.NewGame(model, Output{cons, os.Stderr}); e != nil {
			err = e
		} else {
			ret = GameHelper{cons, game}
		}
	}
	return ret, err
}

type GameHelper struct {
	console C.IConsole
	*R.Game
}

//
// For testing:
//
func (this GameHelper) RunForever() {
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
}
