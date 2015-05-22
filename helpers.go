package sashimi

import (
	"fmt"
	"github.com/ionous/sashimi/console"
	M "github.com/ionous/sashimi/model"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
)

//
func CompileGameWithConsole(s *Script, c console.IConsole, verbose bool) (game *R.Game, model *M.Model, err error) {
	logger := console.NewLogger(verbose)
	if m, e := s.Compile(logger); e != nil {
		err = e
	} else {
		if g, e := R.NewGame(m, c, logger); e != nil {
			err = e
		} else {
			model = m
			game = g
		}
	}
	return game, model, err
}

//
// REQUIRES: a story
func CompileAndRun(s *Script, verbose bool) {
	_, _, err := CompileAndRunWithConsole(s, console.NewConsole(), verbose)
	if err != nil {
		panic(err)
	}
}

// REQUIRES: a story
func CompileAndRunWithConsole(s *Script, c console.IConsole, verbose bool) (*R.Game, *M.Model, error) {
	game, model, err := CompileGameWithConsole(s, c, verbose)
	if err == nil {
		story := game.FindFirstOf(model.Classes.FindClass("stories"))
		if story == nil {
			err = fmt.Errorf("couldnt find story")
		} else {
			if e := game.SendEvent("starting to play", story.String()); e != nil {
				err = e
			} else {
				game.RunForever()
			}
		}
	}
	return game, model, err
}
