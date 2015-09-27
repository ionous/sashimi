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
	started, quit, completed bool
	lastInput                string
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
func NewStandardGame(model *M.Model, frame R.EventFrame, output R.IOutput) (ret StandardStart, err error) {
	if game, e := R.NewGame(model, frame, output); e != nil {
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
				core.SetRight(fmt.Sprintf(`"%s" by %s.`, story.Text("name"), story.Text("author")))
				ret = StandardStart{core}
			}
		}
	}
	return ret, err
}

// Start sends commencing, and returns a new StandardGame.
func (sg *StandardStart) Start() (ret *StandardGame, err error) {
	return &StandardGame{_StandardCore: sg._StandardCore}, err
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
			if in == "start" && !sg.started {
				sg.endTurn("commencing")
			} else {
				if in == "commence" {
					in = sg.lastInput
				} else {
					sg.lastInput = in
				}
				if matcher, e := sg.ParseInput(in); e != nil {
					// FIXFIXFIXFIX
					// change some "report"
					sg.output.Println(e)
				} else if e := matcher.OnMatch(); e != nil {
					sg.output.Println(e)
				} else {
					sg.EndTurn()
				}
			}
		}
	}
	return !sg.IsFinished()
}

// EndTurn finishes the turn for the player.
// ( This is normally called automatically by Input )
func (sg *StandardGame) EndTurn() {
	sg.endTurn("ending the turn")
}
func (sg *StandardGame) endTurn(event string) {
	game := sg.Game
	game.SendEvent(event, sg.story.Id())
	if e := game.ProcessEvents(); e != nil {
		log.Println(e)
	} else {
		if sg.story.Is("completed") {
			sg.completed = true
		}
	}
}
