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
func NewGameSession(model *M.Model) (ret session.ISession, err error) {
	out := &CommandOutput{}
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
			Finish("setting initial position", func() { present(game.Game, player, out) }).
			Finish("ending the turn", func() { endTurn(out) }).
			Finish("ending the story", func() { endStory(out) })
		// add watchers for property changes --
		game.Properties.AddWatcher(PropertyChangeHandler{game.Game, out})
		// now start the game, and start receiving changes --
		if game, e := game.Start(); e != nil {
			err = e
		} else {
			ret = &CommandSession{game, out, nil}
		}
	}
	return ret, err
}

type CommandSession struct {
	game      standard.StandardGame
	output    *CommandOutput
	lastError error
}

//
// ISession implementation
//
func (this *CommandSession) Read(in string) session.ISession {
	//log.Println("cmd session", in)
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

func (this *CommandSession) Write(w io.Writer) (err error) {
	if e := this.lastError; e != nil {
		err, this.lastError = e, nil
	} else {
		//log.Println("cmd session out:", this.output.cmds)
		this.output.Write(w)
	}
	return err
}
