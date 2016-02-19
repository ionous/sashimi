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
	title,
	author,
	playerInput,
	complete,
	turnCount,
	statusLeft,
	statusRight meta.Value
}

func (sg *StandardCore) IsComplete() bool {
	return sg.complete.GetState() == ident.MakeId("completed")
}

func (sg *StandardCore) Started() bool {
	return sg.complete.GetState() != ident.MakeId("starting")
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
	//
	if story, ok := meta.FindFirstOf(game.Model, ident.MakeId("stories")); !ok {
		err = fmt.Errorf("couldn't find story")
	} else if status, ok := meta.FindFirstOf(game.Model, ident.MakeId("status bar instances")); !ok {
		err = fmt.Errorf("couldn't find status bar")
	} else if author, ok := story.FindProperty("author"); !ok {
		err = fmt.Errorf("couldn't find author")
	} else if title, ok := story.FindProperty("name"); !ok {
		err = fmt.Errorf("couldn't find title")
	} else if turnCount, ok := story.FindProperty("turn count"); !ok {
		err = fmt.Errorf("couldn't find turn count")
	} else if playerInput, ok := story.FindProperty("player input"); !ok {
		err = fmt.Errorf("couldn't find completed status")
	} else if completed, ok := story.GetPropertyByChoice(ident.MakeId("completed")); !ok {
		err = fmt.Errorf("couldn't find completed status")
	} else if left, ok := status.FindProperty("left"); !ok {
		err = fmt.Errorf("couldn't find left status")
	} else if right, ok := status.FindProperty("right"); !ok {
		err = fmt.Errorf("couldn't find right status")
	} else {
		core := &StandardCore{
			Game:        game,
			parser:      nil,
			story:       story.GetId(),
			title:       title.GetValue(),
			author:      author.GetValue(),
			playerInput: playerInput.GetValue(),
			complete:    completed.GetValue(),
			turnCount:   turnCount.GetValue(),
			statusLeft:  left.GetValue(),
			statusRight: right.GetValue()}
		core.SetLeft(title.GetValue().GetText())
		core.SetRight(fmt.Sprintf(`"%s" by %s`, title.GetValue().GetText(), author.GetValue().GetText()))
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
