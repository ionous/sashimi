package standard

import (
	"flag"
	"fmt"
	"github.com/ionous/sashimi/console"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/minicon"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
)

// implement IConsole for MiniCon
type MiniConsole struct {
	*minicon.MiniCon
}

func (this MiniConsole) Readln() (string, bool) {
	return this.Update(), true
}

//
func RunGame() {
	var cons console.IConsole

	// command line parsing
	verbose := flag.Bool("verbose", false, "prints log output when true.")
	text := flag.Bool("text", false, "uses the simpler text console when true.")
	dump := flag.Bool("dump", false, "dump the model.")
	flag.Parse()
	//
	if *text || *dump {
		cons = console.NewConsole()
	} else {
		mini := MiniConsole{minicon.NewMiniCon()}
		defer mini.Close()
		cons = mini
	}
	RunGameWithConsole(cons, *verbose, *dump)
}

//
func RunGameWithInput(input []string) {
	cons := console.NewBufCon(input)
	RunGameWithConsole(cons, true, false)
}

//
func RunGameWithConsole(cons console.IConsole, verbose bool, dump bool) {
	s := InitScripts()
	var err error
	logger := console.NewLogger(verbose)
	if model, e := s.Compile(logger); e != nil {
		err = e
	} else {
		if dump {
			model.PrintModel(func(args ...interface{}) { fmt.Println(args...) })
			return
		}
		if game, e := R.NewGame(model, cons, logger); e != nil {
			err = e
		} else {
			if parser, e := R.NewParser(game); e != nil {
				err = e
			} else {
				// FIX: move parser source into the model/parser
				game.PushParserSource(func(g G.Play) (ret G.IObject) {
					return g.The("player")
				})
				game.PushParentLookup(func(g G.Play, o G.IObject) (ret G.IObject) {
					if parent, where := DirectParent(o); where != "" {
						ret = parent
					}
					return ret
				})
				story := game.FindFirstOf(model.Classes.FindClass("stories"))
				if story == nil {
					err = fmt.Errorf("couldnt find story")
				} else {
					if obj, okay := game.FindObject("status bar"); !okay {
						err = fmt.Errorf("couldnt find status bar")
					} else {
						story_ := R.NewObjectAdapter(game, story)
						statusBar := R.NewObjectAdapter(game, obj)
						left, right := story_.Name(), fmt.Sprint("by ", story_.Text("author"))
						statusBar.SetText("left", left)
						statusBar.SetText("right", right)

						// FIX: shouldnt the interface be Go("commence")?
						if e := game.SendEvent("starting to play", story.String()); e != nil {
							err = e
						} else {
							for {
								// process all existing messages in the queue first
								game.Update()

								// update the status bar as needed....
								// ( FIX: status change before the text associated with the change has been teletyped )
								// ( control code handlers -- needed for changing color in text -- might be better )
								newleft, newright := statusBar.Text("left"), statusBar.Text("right")
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
									if story_.Is("completed") {
										break
									}

									if mini, ok := cons.(MiniConsole); ok {
										mini.Flush() // print all remaining teletype text
										mini.Println()
										mini.Println(">", s)
									}
									// input:
									in := parser.NormalizeInput(s)
									if in == "q" || in == "quit" {
										break
									}
									// run some commands:
									if _, res, e := parser.Parse(in); e != nil {
										cons.Println(e)
										continue
									} else if e := res.Run(); e != nil {
										cons.Println(e)
										continue
									}
									game.SendEvent("ending the turn", story.String())
								}
							}
						}
					}
				}
			}
		}
	}

	if err != nil {
		panic(err)
	}
}
