package commands

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/web/session"
	"io"
)

//
// game session implementation
//
func NewSession(id string, model *M.Model) (ret session.ISession, err error) {
	out := &CommandOutput{id: id}
	// after creating the game, but vefore running it --
	if game, e := standard.NewStandardGame(model, out); e != nil {
		err = e
	} else if player, ok := game.FindObject("Player"); !ok {
		err = fmt.Errorf("unknown player")
	} else {
		// find the player object --
		player := R.NewObjectAdapter(game.Game, player)
		// setup system event callbacks --
		game.SystemActions.
			Finish("setting initial position", func() {
			// view := SerializeView(game.Model, "StatusBar")
			// out.NewCommand("present", view)
			present(game.Game, player, out)
		}).
			Finish("ending the turn", func() { endTurn(out) }).
			Finish("ending the story", func() { endStory(out) })
		// add watchers for property changes --
		game.Properties.AddWatcher(PropertyChangeHandler{game.Game, out})
		// now start the game, and start receiving changes --
		if game, e := game.Start(); e != nil {
			err = e
		} else {
			ret = &CommandSession{game, id, out, nil}
		}
	}
	return ret, err
}

type CommandSession struct {
	game      standard.StandardGame
	state     string
	output    *CommandOutput
	lastError error
}

func (this *CommandSession) Game() *standard.StandardGame {
	return &this.game
}

//
// ISession implementation
//
func (this *CommandSession) Read(in string) session.ISession {
	if this.lastError == nil {
		if this.game.IsQuit() {
			this.lastError = session.SessionClosed{"player quit"}
		} else if this.game.IsFinished() {
			this.lastError = session.SessionClosed{"game finished"}
		} else {
			this.game.Input(in)
		}
	}
	return this
}

//
// ISession implementation which writes json via CommandOutput.
//
func (this *CommandSession) Write(w io.Writer) (err error) {
	if e := this.lastError; e != nil {
		err, this.lastError = e, nil
	} else {
		this.output.Write(w)
	}
	return err
}
