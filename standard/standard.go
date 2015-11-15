package standard

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/runtime/parse"
	"github.com/ionous/sashimi/util/ident"
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
	Parser parser.P
	story  ident.Id
	title,
	author,
	playerInput,
	complete,
	statusLeft,
	statusRight api.Value
}

// Left status bar text.
func (sc *StandardCore) Left() string {
	return sc.statusLeft.GetText()
}

// Right status bar text.
func (sc *StandardCore) Right() string {
	return sc.statusRight.GetText()
}

// SetLeft status bar text.
func (sc *StandardCore) SetLeft(status string) {
	sc.statusLeft.SetText(status)
}

// SetRight status bar text.
func (sc *StandardCore) SetRight(status string) {
	sc.statusRight.SetText(status)
}

// NewStandardGame creates a game which is based on the standard rules.
func NewStandardGame(game *R.Game) (ret StandardStart, err error) {
	if parser, e := parse.NewObjectParser(game.ModelApi, ident.MakeId("player")); e != nil {
		err = e
	} else {
		game.PushParentLookup(func(g G.Play, o G.IObject) (ret G.IObject) {
			if parent, where := DirectParent(o); where != "" {
				ret = parent
			}
			return ret
		})
		//
		if story, ok := api.FindFirstOf(game.ModelApi, ident.MakeId("stories")); !ok {
			err = fmt.Errorf("couldn't find story")
		} else if status, ok := api.FindFirstOf(game.ModelApi, ident.MakeId("status bar instances")); !ok {
			err = fmt.Errorf("couldn't find status bar")
		} else if title, ok := story.GetProperty(ident.MakeId("name")); !ok {
			err = fmt.Errorf("couldn't find title")
		} else if author, ok := story.GetProperty(ident.MakeId("author")); !ok {
			err = fmt.Errorf("couldn't find author")
		} else if playerInput, ok := story.GetProperty(ident.MakeId("player input")); !ok {
			err = fmt.Errorf("couldn't find completed status")
		} else if completed, ok := story.GetPropertyByChoice(ident.MakeId("completed")); !ok {
			err = fmt.Errorf("couldn't find completed status")
		} else if left, ok := status.GetProperty(ident.MakeId("left")); !ok {
			err = fmt.Errorf("couldn't find left status")
		} else if right, ok := status.GetProperty(ident.MakeId("right")); !ok {
			err = fmt.Errorf("couldn't find right status")
		} else {
			core := StandardCore{game, parser,
				story.GetId(),
				title.GetValue(),
				author.GetValue(),
				playerInput.GetValue(),
				completed.GetValue(),
				left.GetValue(),
				right.GetValue()}

			core.SetLeft(title.GetValue().GetText())
			core.SetRight(fmt.Sprintf(`"%s" by %s.`, title.GetValue().GetText(), author.GetValue().GetText()))
			ret = StandardStart{core}
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
				if e := sg.playerInput.SetText(in); e != nil {
					err = e
				} else if act, e := sg.Game.QueueEvent("parsing player input", sg.story); e != nil {
					err = e
				} else if e := sg.Game.ProcessEvents(); e != nil {
					err = e
				} else if act.Cancelled() {
					log.Print("input cancelled")
					err = nil // input cancelling is not an error
				} else if _, matcher, e := sg.Parser.ParseInput(in); e != nil {
					err = e
				} else if act, objs, e := matcher.(*parse.ObjectMatcher).GetMatch(); e != nil {
					err = e
				} else {
					sg.RunTurn(act, objs)
				}
			}
		}
	}
	return err
}

// EndTurn finishes the turn for the player.
// ( This is normally called automatically by Input )
func (sg *StandardGame) RunTurn(act api.Action, objs []api.Instance) {
	sg.Game.QueueAction(act, objs)
	sg.endTurn("ending the turn")
}

func (sg *StandardGame) endTurn(event string) {
	sg.Game.QueueEvent(event, sg.story)
	if e := sg.Game.ProcessEvents(); e != nil {
		log.Println(e)
	} else {
		if sg.complete.GetState() == ident.MakeId("completed") {
			sg.completed = true
		}
	}
}
