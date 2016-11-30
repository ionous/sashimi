package framework

import (
	"errors"
	"github.com/ionous/mars/rt"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/parser"
	"github.com/ionous/sashimi/play"
	"github.com/ionous/sashimi/play/api"
	"github.com/ionous/sashimi/play/parse"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

// StandardCore assists the transformation of a StandardStart into a StandardGame.
type StandardCore struct {
	play.Game
	parser *parser.P // P is a struct, cached via getParser()
	story  rt.Object
	playerInput,
	complete,
	turnCount meta.Property
}

func (sg *StandardCore) IsComplete() bool {
	return sg.complete.GetGeneric().(rt.State) == rt.State("completed")
}

func (sg *StandardCore) Started() bool {
	return sg.complete.GetGeneric().(rt.State) != rt.State("starting")
}

// frame is the turn count + 1 ( so that it's never zero while playing )
func (sc *StandardCore) Frame() (ret int) {
	if sc.Started() {
		// FIX? this is a little fragile: the frame count should be 1 for the data sent by the first frame; 0 before.
		// commencing sets us to story(sc) to started, end turn increments the turn count, and finally the session samples the Frame() just before sending the response.
		i := sc.turnCount.GetGeneric().(rt.Number).Int()
		ret = i + 1
	}
	return
}

// NewStandardGame creates a game which is based on the standard rules.
func NewStandardCore(game play.Game) (ret *StandardCore, err error) {
	if story, ok := meta.FindFirstOf(game, ident.MakeId("stories")); !ok {
		err = errors.New("couldn't find story object")
	} else if turnCount, ok := story.FindProperty("turn count"); !ok {
		err = errors.New("couldn't find turn count property")
	} else if playerInput, ok := story.FindProperty("player input"); !ok {
		err = errors.New("couldn't find player input property")
	} else if completed, ok := story.GetPropertyByChoice(ident.MakeId("completed")); !ok {
		err = errors.New("couldn't find completed property")
	} else {
		core := &StandardCore{
			Game:        game,
			parser:      nil,
			story:       rt.Object{story},
			playerInput: playerInput,
			complete:    completed,
			turnCount:   turnCount,
		}
		ret = core
	}
	return
}

// NOTE: input should be normalized!
func (sg *StandardCore) HandleInput(in string) (err error) {
	if sg.IsComplete() {
		err = errutil.New("handle input", in, "game is finished.")
	} else {
		if in == "start" && !sg.Started() {
			sg.Game.Log("starting game")
			err = sg.EndTurn("commence")
		} else {
			if in == "commence" {
				in = sg.playerInput.GetGeneric().(rt.Text).String()
			}
			//
			if e := sg.playerInput.SetGeneric(rt.Text{in}); e != nil {
				err = e
			} else if e := sg.Game.RunAction(ident.MakeId("parse player input"), sg.Game, sg.story); e != nil {
				err = sg.EndTurn("end turn")
			} else if parser, e := sg.getParser(); e != nil {
				err = e
			} else {
				if _, matcher, e := parser.ParseInput(in); e != nil {
					err = e
				} else if act, insts, e := matcher.(*parse.ObjectMatcher).GetMatch(); e != nil {
					err = e
				} else {
					objs := make([]meta.Generic, len(insts))
					for i, inst := range insts {
						objs[i] = rt.Object{inst}
					}
					sg.Game.Log("running action", act.GetId(), objs)
					if e := sg.Game.RunAction(act.GetId(), sg.Game, objs...); e != nil {
						err = e
					} else {
						sg.Game.Log("ending turn")
						err = sg.EndTurn("end turn")
					}
				}
			}

		}
	}
	return err
}

// STORE-FIX: it takes a bit of work create a parser from the model.
// for the command sessions ( local server, app engine, ) we only need the parser with raw input
// so we leave it nil till its needed.
// FUTURE: implement zero-copy parser by re-formating model nto something more amenable for the parser.
func (sg *StandardCore) getParser() (ret parser.P, err error) {
	if sg.parser != nil {
		ret = *sg.parser
	} else {
		// cache!
		if parser, e := parse.NewObjectParser(sg, "player"); e != nil {
			err = e
		} else {
			ret, sg.parser = parser, &parser
		}
	}
	return
}

func (sg *StandardCore) EndTurn(action string) error {
	return sg.Game.RunAction(ident.MakeId(action), sg.Game, sg.story)
}
