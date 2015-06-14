package commands

import (
	"encoding/json"
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/net/session"
	"github.com/ionous/sashimi/standard"
	"io"
	"log"
	"sync"
)

//
// Create a session.
//
func NewCommandSession(id string, model *M.Model) (ret *CommandSession, err error) {
	output := NewCommandOutput(id)
	// after creating the game, but vefore running it --
	if game, e := standard.NewStandardGame(model, output); e != nil {
		err = e
	} else {
		// setup system event callbacks --
		game.SystemActions.
			Capture("setting initial position", output.changedLocation)
		// add watchers for property changes --
		game.Properties.AddWatcher(PropertyChangeHandler{game.Game, output})
		// now start the game, and start receiving changes --
		if game, e := game.Start(); e != nil {
			err = e
		} else {
			ret = &CommandSession{game, output, 1, &sync.RWMutex{}}
		}
	}
	return ret, err
}

//
// Game data associated with a client's particular sessions.
//
type CommandSession struct {
	game       standard.StandardGame
	output     *CommandOutput
	frameCount int
	*sync.RWMutex
}

//
// Return the named game resource
//
func (this *CommandSession) Find(name string) (ret resource.IResource, okay bool) {
	switch name {
	// by default, objects are grouped by their class:
	default:
		if cls, plural := this.game.Model.Classes.FindClass(name); plural {
			ret, okay = ObjectResource(this.game.Game, cls, this.output.serial.ObjectSerializer), true
		}
	// a request for information about a class?
	case "class":
		ret, okay = ClassResource(this.game.Model), true
	}
	return ret, okay
}

//
// The session iteslf doesnt return any data on query.
//
func (this *CommandSession) Query() (ret resource.Document) {
	return
}

//
// Post input to the game.
//
func (this *CommandSession) Post(reader io.Reader) (ret resource.Document, err error) {
	// hrmm.... by sending the io.Reader, we can "decode" different types of json...
	// but, we have to do that wherever we need... [ add a Blank struct maker ? ]
	decoder, input := json.NewDecoder(reader), CommandInput{}
	if e := decoder.Decode(&input); e != nil {
		err = e
	} else {
		if e := this._handleInput(input); e != nil {
			err = e
		} else {
			doc := resource.NewDocumentBuilder(&ret)
			this.output.FlushDocument(doc)
		}
	}
	return ret, err
}

//
// Send the passed input to the game.
//
func (this *CommandSession) _handleInput(input CommandInput) (err error) {
	if this.game.IsQuit() {
		err = session.SessionClosed{"player quit"}
	} else if this.game.IsFinished() {
		err = session.SessionClosed{"game finished"}
	} else {
		// Run raw input:
		if input.Input != "" {
			log.Println("Input", input)
			this.game.Input(input.Input)
		} else {
			// Run json'd clicky action:
			if act, ok := this.game.Model.Actions[M.MakeStringId(input.Action)]; !ok {
				err = fmt.Errorf("unknown action %s", input.Action)
				//FIX? RunActions injects the player, that works out well, but is a little strange.
			} else if e := this.game.RunAction(act, input.Nouns()); e != nil {
				err = e
			} else {
				this.game.EndTurn() // game.Input() does this automatically (dont ask)
			}
		}
	}
	return err
}
