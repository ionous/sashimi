package standard

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	R "github.com/ionous/sashimi/runtime"
	"log"
)

// StardardStart assists the creation of a standard game.
// see: NewStandardGame()
type StandardStart struct {
	_StandardCore
}

// StandardGame wraps the runtime.Game with the standard rules.
type StandardGame struct {
	_StandardCore
	quit, completed bool
}

// _StandardCore assists the transformation of a StandardStart into a StandardGame.
type _StandardCore struct {
	*R.Game
	*Parser
	output        R.IOutput
	story, status G.IObject
}

// Left status bar text.
func (sc *_StandardCore) Left() string {
	return sc.status.Text("left")
}

// Right status bar text.
func (sc *_StandardCore) Right() string {
	return sc.status.Text("right")
}

// SetLeft status bar text.
func (sc *_StandardCore) SetLeft(status string) {
	sc.status.SetText("left", status)
}

// SetRight status bar text.
func (sc *_StandardCore) SetRight(status string) {
	sc.status.SetText("right", status)
}

// NewStandardGame creates a game which is based on the standard rules.
func NewStandardGame(model *M.Model, output R.IOutput) (ret StandardStart, err error) {
	if game, e := R.NewGame(model, output); e != nil {
		err = e
	} else if parser, e := NewParser(game); e != nil {
		err = e
	} else {
		g := R.NewGameAdapter(game)
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
		if story, found := G.Any(g, "stories"); !found {
			err = fmt.Errorf("couldn't find story")
		} else {
			if status, found := G.Any(g, "status bar instances"); !found {
				err = fmt.Errorf("couldn't find status bar")
			} else {
				core := _StandardCore{game, parser, output, story, status}
				core.SetLeft(story.Text("name"))
				core.SetRight(fmt.Sprint(story.Text("name"), "by ", story.Text("author")))
				ret = StandardStart{core}
			}
		}
	}
	return ret, err
}

// Start sends starting to play, and returns a new StandardGame.
func (sg *StandardStart) Start() (ret *StandardGame, err error) {
	// FIX: shouldnt the interface be Go("commence")?
	if e := sg.SendEvent("starting to play", sg.story.Id()); e != nil {
		err = e
	} else {
		// process all existing messages in the queue first
		if e := sg.ProcessEvents(); e != nil {
			err = e
		}
	}
	return &StandardGame{sg._StandardCore, false, false}, err
}

// IsQuit when the user has requested to quit the game.
func (sg *StandardGame) IsQuit() bool {
	return sg.quit
}

// IsFinished when the user has completed the game or quit the game.
func (sg *StandardGame) IsFinished() bool {
	return sg.quit || sg.completed
}

// Input turns the passed user input to a game command.
// Returns false if the game IsFinished.
// (automatically ends the turn )
func (sg *StandardGame) Input(s string) bool {
	if !sg.IsFinished() {
		in := sg.NormalizeInput(s)
		if in == "q" || in == "quit" {
			sg.quit = true
		} else {
			if matcher, e := sg.ParseInput(in); e != nil {
				sg.output.Println(e)
			} else if e := matcher.OnMatch(); e != nil {
				sg.output.Println(e)
			} else {
				sg.EndTurn()
			}
		}
	}
	return !sg.IsFinished()
}

// EndTurn finishes the turn for the player.
// ( This is normally called automatically by Input )
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
