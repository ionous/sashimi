package standard

import (
	"fmt"
	"github.com/ionous/sashimi/compiler/metal"
	C "github.com/ionous/sashimi/console"
	"github.com/ionous/sashimi/minicon"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/script"
	"io"
	"io/ioutil"
	"os"
)

// implement IConsole for MiniCon
// FIX-this shouldn't be in root-standard;
// standard.miniconsole maybe?
type MiniConsole struct {
	*minicon.MiniCon
}

func (this MiniConsole) Readln() (string, bool) {
	return this.Update(), true
}

// simplest interface:
func Run(scriptCallback script.InitCallback) {
	script.AddScript(scriptCallback)
	RunGame(ParseCommandLine())
}

func RunGame(opt Options) (err error) {
	script := script.InitScripts()
	return RunScript(script, opt)
}

func GetWriter(opt Options) (ret io.Writer) {
	if opt.verbose {
		ret = os.Stderr
	} else {
		ret = ioutil.Discard
	}
	return
}

func GetConsole(opt Options) (ret C.IConsole) {
	if opt.hasConsole {
		ret = opt.cons
	} else if opt.text {
		ret = C.NewConsole()
	} else {
		ret = MiniConsole{minicon.NewMiniCon()}
	}
	return
}

func RunScript(script *script.Script, opt Options) (err error) {
	writer := GetWriter(opt)
	if model, e := script.Compile(writer); e != nil {
		err = e
	} else if opt.dump {
		model.Model.PrintModel(func(args ...interface{}) { fmt.Println(args...) })
	} else {
		cons := GetConsole(opt)
		defer cons.Close()

		cfg := R.NewConfig().SetCalls(model.Calls).SetOutput(NewStandardOutput(cons, writer)).SetParentLookup(ParentLookup{})
		modelApi := metal.NewMetal(model.Model, make(metal.ObjectValueMap))
		if g, e := cfg.NewGame(modelApi); e != nil {
			err = e
		} else {
			err = PlayGame(cons, g)
		}
	}
	if err != nil {
		panic(err)
	}
	return
}

func PlayGame(cons C.IConsole, g R.Game) (err error) {
	return PlayGameUpdate(cons, g, nil)
}

func PlayGameUpdate(cons C.IConsole, g R.Game, endFrame func()) (err error) {
	if game, e := NewStandardGame(g); e != nil {
		err = e
	} else {
		left := game.title.GetText()
		right := fmt.Sprint(" by ", game.author.GetText())
		game.SetLeft(left)
		game.SetRight(right)
		immediate := true
		if game, e := game.Start(immediate); e != nil {
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

					if e := game.Input(s); e != nil {
						cons.Println(e.Error())
					} else if game.IsFinished() {
						if useMini && !game.IsQuit() {
							mini.Update()
						}
						break
					}
					if endFrame != nil {
						endFrame()
					}
				}
			}
		}
	}
	return err
}
