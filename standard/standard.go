package standard

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	R "github.com/ionous/sashimi/runtime"
	"log"
)

type StandardCore struct {
	*R.Game
	*R.ObjectParser
	output        R.IOutput
	story, status R.ObjectAdapter
}
type StandardStart struct {
	StandardCore
}
type StandardGame struct {
	StandardCore
	quit, completed bool
}

func (sc *StandardCore) Left() string {
	return sc.status.Text("left")
}

func (sc *StandardCore) Right() string {
	return sc.status.Text("right")
}

func (sc *StandardCore) SetLeft(status string) {
	sc.status.SetText("left", status)
}

func (sc *StandardCore) SetRight(status string) {
	sc.status.SetText("right", status)
}

func NewStandardGame(model *M.Model, output R.IOutput) (ret StandardStart, err error) {
	if game, e := R.NewGame(model, output); e != nil {
		err = e
	} else if parser, e := R.NewParser(game); e != nil {
		err = e
	} else {
		//
		parser.PushParserSource(func(g G.Play) G.IObject {
			return g.The("player")
		})
		//
		game.PushParentLookup(func(g G.Play, o G.IObject) (ret G.IObject) {
			if parent, where := DirectParent(o); where != "" {
				ret = parent
			}
			return ret
		})
		//
		storyObject := game.FindFirstOf(model.Classes.FindClass("stories"))
		if storyObject == nil {
			err = fmt.Errorf("couldn't find story")
		} else {
			if statusObject, ok := game.FindObject("status bar"); !ok {
				err = fmt.Errorf("couldn't find status bar")
			} else {
				story := R.NewObjectAdapter(game, storyObject)
				status := R.NewObjectAdapter(game, statusObject)
				//
				core := StandardCore{game, parser, output, story, status}
				core.SetLeft(story.Text("name"))
				core.SetRight(fmt.Sprint(story.Text("name"), "by ", story.Text("author")))
				ret = StandardStart{core}
			}
		}
	}
	return ret, err
}

//
// sends starting to play, and returns a new game.
//
func (sg *StandardStart) Start() (ret StandardGame, err error) {
	// FIX: shouldnt the interface be Go("commence")?
	if e := sg.SendEvent("starting to play", sg.story.Id()); e != nil {
		err = e
	} else {
		// process all existing messages in the queue first
		if e := sg.ProcessEvents(); e != nil {
			err = e
		}
	}
	return StandardGame{sg.StandardCore, false, false}, err
}

func (sg *StandardGame) IsQuit() bool {
	return sg.quit
}

func (sg *StandardGame) IsFinished() bool {
	return sg.quit || sg.completed
}

//
// return false if the game has finished
// (automatically ends the turn )
//
func (sg *StandardGame) Input(s string) bool {
	if !sg.IsFinished() {
		in := sg.NormalizeInput(s)
		if in == "q" || in == "quit" {
			sg.quit = true
		} else {
			if _, e := sg.Parse(in); e != nil {
				sg.output.Println(e)
			}
			sg.EndTurn()
		}
	}
	return !sg.IsFinished()
}

func (sg *StandardGame) EndTurn() {
	game := sg.Game
	game.SendEvent("ending the turn", sg.story.Id())
	if e := game.ProcessEvents(); e != nil {
		log.Println(e)
	} else {
		if sg.story.Is("completed") {
			sg.completed = true
		}
	}
}
