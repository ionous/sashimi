package framework

import (
	"github.com/ionous/sashimi/parser"
	"github.com/ionous/sashimi/play"
)

// NewStardardGame assists the creation of a standard game.
// returns StandardStart which supports the interface Start()
func NewStandardGame(game play.Game) (ret StandardStart, err error) {
	if core, e := NewStandardCore(game); e != nil {
		err = e
	} else {
		ret = StandardStart{core}
	}
	return
}

// see: NewStandardGame()
type StandardStart struct {
	*StandardCore
}

// Start sends commencing, and returns a new StandardGame.
// FIX: no longer sends commencing, that's done by input "start"
func (sg StandardStart) Start() (ret *StandardGame, err error) {
	ret = &StandardGame{StandardCore: sg.StandardCore}
	err = ret.EndTurn("commence")
	return
}

// StandardGame wraps the core rules with "quit"
// FIX: move this to boilerplate, etc.
type StandardGame struct {
	*StandardCore
	quit bool
}

func (sg *StandardGame) IsQuit() bool {
	return sg.quit
}

// Input turns the passed user input to a game command.
// (automatically ends the turn )
func (sg *StandardGame) Input(s string) (err error) {
	if !sg.quit {
		in := parser.NormalizeInput(s)
		if in == "q" || in == "quit" {
			sg.quit = true
		} else {
			err = sg.HandleInput(in)
		}
	}
	return
}
