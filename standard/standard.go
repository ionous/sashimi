package standard

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
	"log"
)

// FIX: we have the concept floating in other fixes of "function" globals
// and that might be needed for this, where we really dont want *shared* globals
// you would want this tied to session, if at all possible.
var Debugging bool

// StardardStart assists the creation of a standard game.
// see: NewStandardGame()
type StandardStart struct {
	StandardCore
}

// StandardGame wraps the runtime.Game with the standard rules.
type StandardGame struct {
	StandardCore
	started, quit, completed bool
	lastInput                string
}

// StandardCore assists the transformation of a StandardStart into a StandardGame.
type StandardCore struct {
	*R.Game
	Parser        parser.P
	story, status G.IObject
}

// Left status bar text.
func (sc *StandardCore) Left() string {
	return sc.status.Text("left")
}

// Right status bar text.
func (sc *StandardCore) Right() string {
	return sc.status.Text("right")
}

// SetLeft status bar text.
func (sc *StandardCore) SetLeft(status string) {
	sc.status.SetText("left", status)
}

// SetRight status bar text.
func (sc *StandardCore) SetRight(status string) {
	sc.status.SetText("right", status)
}

// NewStandardGame creates a game which is based on the standard rules.
func NewStandardGame(game *R.Game) (ret StandardStart, err error) {
	if parser, e := R.NewObjectParser(game, game.ModelApi, "player"); e != nil {
		err = e
	} else {
		g := R.NewGameAdapter(game)
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
				core := StandardCore{game, parser, story, status}
				core.SetLeft(story.Text("name"))
				core.SetRight(fmt.Sprintf(`"%s" by %s.`, story.Text("name"), story.Text("author")))
				ret = StandardStart{core}
			}
		}
	}
	return ret, err
}

// Start sends commencing, and returns a new StandardGame.
// FIX: no longer sends commencing, that's done by input "start"
func (sg *StandardStart) Start(immediate bool) (ret *StandardGame, err error) {
	ret = &StandardGame{StandardCore: sg.StandardCore}
	if immediate {
		ret.endTurn("commencing")
	}
	return ret, err
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
// (automatically ends the turn )
func (sg *StandardGame) Input(s string) (err error) {
	if !sg.IsFinished() {
		in := parser.NormalizeInput(s)
		if in == "q" || in == "quit" {
			sg.quit = true
		} else {
			if in == "start" && !sg.started {
				sg.started = true
				sg.endTurn("commencing")
			} else {
				if in == "commence" {
					in = sg.lastInput
				} else {
					sg.lastInput = in
				}
				if matcher, e := sg.Parser.ParseInput(in); e != nil {
					err = e
				} else if e := matcher.OnMatched(); e != nil {
					err = e
				} else {
					sg.EndTurn()
				}
			}
		}
	}
	return err
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
