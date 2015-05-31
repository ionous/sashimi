package standard

import (
	"fmt"
	C "github.com/ionous/sashimi/console"
	"github.com/ionous/sashimi/minicon"
	S "github.com/ionous/sashimi/script"
	"io"
	"io/ioutil"
	"os"
)

// implement IConsole for MiniCon
type MiniConsole struct {
	*minicon.MiniCon
}

func (this MiniConsole) Readln() (string, bool) {
	return this.Update(), true
}

// simplest interface:
func Run(cb S.InitCallback) (err error) {
	S.AddScript(cb)
	return RunGame(ParseCommandLine())
}

// second simplest interface:
func RunGame(opt Options) (err error) {
	// tease out options settings:
	cons, verbose, dump := opt.cons, opt.verbose, opt.dump
	if !opt.hasConsole {
		if opt.text {
			cons = C.NewConsole()
		} else {
			mini := MiniConsole{minicon.NewMiniCon()}
			defer mini.Close()
			cons = mini
		}
	}
	var writer io.Writer
	if verbose {
		writer = os.Stderr
	} else {
		writer = ioutil.Discard
	}
	if model, e := S.InitScripts().Compile(writer); e != nil {
		err = e
	} else {
		if dump {
			model.PrintModel(func(args ...interface{}) { fmt.Println(args...) })
			return
		}
		if game, e := NewStandardGame(model, StandardOutput{cons, writer}); e != nil {
			err = e
		} else {
			left, right := game.story.Name(), fmt.Sprint("by ", game.story.Text("author"))
			game.SetLeft(left)
			game.SetRight(right)
			if game, e := game.Start(); e != nil {
				err = e
			} else {
				for {
					// update the status bar as needed....
					// ( FIX: status change before the text associated with the change has been teletyped )
					// ( control code handlers -- needed for changing color in text -- might be better )
					newleft, newright := game.Left(), game.Right()
					if left != newleft || right != newright {
						if mini, ok := cons.(MiniConsole); ok {
							mini.Status.Left, mini.Status.Right = newleft, newright
							mini.RefreshDisplay()
						}
					}

					// read new input
					if s, ok := cons.Readln(); !ok {
						break
					} else {
						mini, useMini := cons.(MiniConsole)
						if useMini {
							mini.Flush() // print all remaining teletype text
							mini.Println()
							mini.Println(">", s)
						}

						if !game.Input(s) {
							if useMini && !game.IsQuit() {
								mini.Update()
							}
							break
						}
					}
				}
			}
		}
	}

	if err != nil {
		panic(err)
	}
	return err
}
