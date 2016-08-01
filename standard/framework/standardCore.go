package framework

import (
	"fmt"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/parse"
	"github.com/ionous/sashimi/util/ident"
	"log"
)

// StandardCore assists the transformation of a StandardStart into a StandardGame.
type StandardCore struct {
	R.Game
	parser *parser.P // P is a struct, cached via getParser()
	story  ident.Id
	playerInput,
	complete,
	turnCount meta.Value
}

func (sg *StandardCore) IsComplete() bool {
	return sg.complete.GetState() == ident.MakeId("completed")
}

func (sg *StandardCore) Started() bool {
	return sg.complete.GetState() != ident.MakeId("starting")
}

// frame is the turn count + 1 ( so that it's never zero while playing )
func (sc *StandardCore) Frame() (ret int) {
	if sc.Started() {
		// FIX? this is a little fragile: the frame count should be 1 for the data sent by the first frame; 0 before.
		// commencing sets us to story(sc) to started, end turn increments the turn count, and finally the session samples the Frame() just before sending the response.
		ret = int(sc.turnCount.GetNum()) + 1
	}
	return
}

// NewStandardGame creates a game which is based on the standard rules.
func NewStandardCore(game R.Game) (ret *StandardCore, err error) {
	if story, ok := meta.FindFirstOf(game.Model, ident.MakeId("stories")); !ok {
		err = fmt.Errorf("couldn't find story object")
	} else if turnCount, ok := story.FindProperty("turn count"); !ok {
		err = fmt.Errorf("couldn't find turn count property")
	} else if playerInput, ok := story.FindProperty("player input"); !ok {
		err = fmt.Errorf("couldn't find player input property")
	} else if completed, ok := story.GetPropertyByChoice(ident.MakeId("completed")); !ok {
		err = fmt.Errorf("couldn't find completed property")
	} else {
		core := &StandardCore{
			Game:        game,
			parser:      nil,
			story:       story.GetId(),
			playerInput: playerInput.GetValue(),
			complete:    completed.GetValue(),
			turnCount:   turnCount.GetValue(),
		}
		ret = core
	}
	return
}

// NOTE: input should be normalized!
func (sg *StandardCore) HandleInput(in string) (err error) {
	if sg.IsComplete() {
		log.Println("complete")
	} else {
		if in == "start" && !sg.Started() {
			//log.Println("starting")
			sg.EndTurn("commence")
		} else {
			if in == "commence" {
				in = sg.playerInput.GetText()
			}
			if e := sg.playerInput.SetText(in); e != nil {
				err = e
			} else if act, e := sg.Game.QueueAction("parse player input", sg.story); e != nil {
				err = e
			} else if e := sg.Game.ProcessActions(); e != nil {
				err = e
			} else if act.Cancelled() {
				sg.EndTurn("end turn")
				// NOTE: canceling is not an error
			} else if parser, e := sg.getParser(); e != nil {
				log.Println("error getting parser", e)
				err = e
			} else {
				if _, matcher, e := parser.ParseInput(in); e != nil {
					fmt.Println("error parsing input", in, e)
					err = e
				} else if act, objs, e := matcher.(*parse.ObjectMatcher).GetMatch(); e != nil {
					err = e
					//log.Println("error matching input", err)
				} else {
					sg.Game.QueueActionInstances(act, objs)
					sg.EndTurn("end turn")
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
		if parser, e := parse.NewObjectParser(sg.Model, ident.MakeId("player")); e != nil {
			err = e
		} else {
			ret, sg.parser = parser, &parser
		}
	}
	return
}

func (sg *StandardCore) EndTurn(action string) {
	if _, e := sg.Game.QueueAction(action, sg.story); e != nil {
		log.Println(e)
	} else if e := sg.Game.ProcessActions(); e != nil {
		log.Println(e)
	}
}
