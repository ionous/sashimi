package app

import (
	"encoding/json"
	"fmt"
	E "github.com/ionous/sashimi/event"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/net/session"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/util/ident"
	"io"
	"log"
	"sync"
)

//
// Create a session.
//
func NewCommandSession(id string, model *M.Model, calls R.Callbacks) (ret *CommandSession, err error) {
	output := NewCommandOutput(id)
	// event start/end frame
	frame := func(tgt E.ITarget, msg *E.Message) func() {
		// msg.Data is RuntimeAction -- theres not really parameters for events right now
		// other than tgt, src, ctx right now.
		output.flushPending()
		//msg.Data == RunTimeAction
		output.events.PushEvent(msg.Name, tgt, nil)
		return func() {
			output.flushPending()
			output.events.PopEvent()
		}
	}
	cfg := R.RuntimeConfig{Calls: calls, Frame: frame, Output: output}
	if game, e := cfg.NewGame(model); e != nil {
	} else {
		// after creating the game, but vefore running it --
		if game, e := standard.NewStandardGame(game); e != nil {
			err = e
		} else {
			// setup system event callbacks --
			// STORE-FIX: can this be removed? why is it needed?
			setInitialPos := R.MakeStringId("set initial position")
			if e := game.SystemActions.Capture(setInitialPos, output.changedLocation); e != nil {
				err = e
			} else {
				// add watchers for property changes --
				game.Properties.AddWatcher(output)
				// now start the game, and start receiving changes --
				immediate := false
				if game, e := game.Start(immediate); e != nil {
					err = e
				} else {
					ret = &CommandSession{game, model, output, 1, &sync.RWMutex{}}
				}
			}
		}
	}
	return ret, err
}

//
// Game data associated with a client's particular sessions.
//
type CommandSession struct {
	game       *standard.StandardGame
	Model      *M.Model
	output     *CommandOutput
	frameCount int
	*sync.RWMutex
}

//
// Return the named game resource
//
func (sess *CommandSession) Find(name string) (ret resource.IResource, okay bool) {
	switch name {
	// by default, objects are grouped by their class:
	default:
		if cls, plural := sess.Model.Classes.FindClass(name); plural {
			ret, okay = ObjectResource(sess.game.Game, cls.Id, sess.output.serial.ObjectSerializer), true
		}
	// a request for information about a class:
	case "class":
		ret, okay = ClassResource(sess.Model.Classes), true
		// a request for information about a parser input action:
	case "action":
		ret, okay = ParserResource(sess.game.ModelApi), true
	}
	return ret, okay
}

//
// The session iteslf doesnt return any data on query.
//
func (sess *CommandSession) Query() (ret resource.Document) {
	return
}

//
// Post input to the game.
//
func (sess *CommandSession) Post(reader io.Reader) (ret resource.Document, err error) {
	// hrmm.... by sending the io.Reader, we can "decode" different types of json...
	// but, we have to do that wherever we need... [ add a Blank struct maker ? ]
	decoder, input := json.NewDecoder(reader), CommandInput{}
	if e := decoder.Decode(&input); e != nil {
		err = e
	} else {
		if e := sess._handleInput(input); e != nil {
			err = e
		} else {
			doc := resource.NewDocumentBuilder(&ret)
			sess.output.FlushDocument(doc)
		}
	}
	return ret, err
}

//
// Send the passed input to the game.
//
func (sess *CommandSession) _handleInput(input CommandInput) (err error) {
	if sess.game.IsQuit() {
		err = session.SessionClosed{"player quit"}
	} else if sess.game.IsFinished() {
		err = session.SessionClosed{"game finished"}
	} else {
		// Run raw input:
		if input.Input != "" {
			log.Println("Input", input)
			sess.game.Input(input.Input)
		} else {
			// Run json'd clicky action:
			id := ident.MakeId(input.Action)
			if act, ok := sess.game.ModelApi.GetAction(id); !ok {
				err = fmt.Errorf("unknown action %s", input.Action)
				//FIX? RunActions injects the player, that works out well, but is a little strange.
			} else {
				if om, e := sess.game.StandardParser.ObjectParser.GetObjectMatcher(act); e != nil {
					err = e
				} else {
					for _, n := range input.Nouns() {
						id := R.MakeStringId(n)
						if ok := om.MatchId(id); !ok {
							err = fmt.Errorf("unknown object %s; %v", n, id)
							break
						}
					}
					if err == nil {
						if e := om.OnMatch(); e != nil {
							err = e
						} else {
							sess.game.EndTurn() // game.Input() does sess automatically (dont ask)
						}
					}
				}
			}
		}
	}
	return err
}
