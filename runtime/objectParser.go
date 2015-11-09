package runtime

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/parser"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/errutil"
)

// ObjectParser adapts the parser to the runtime,
// queuing actions for game objects based upon user input.
type ObjectParser struct {
	game   *Game // for generating object adapters after a noun is matfched
	Parser parser.Parser
	source *ParserSourceStack
}

//
// NormalizeInput wraps parser.NormalizeInput for convenience's sake.
//
func (op *ObjectParser) NormalizeInput(s string) string {
	return parser.NormalizeInput(s)
}

//
// NewObjectParser and add all commands and patterns.
//
func NewObjectParser(game *Game,
) (ret *ObjectParser, err error,
) {
	model := game.ModelApi
	p := parser.NewParser()
	op := &ObjectParser{game, p, &ParserSourceStack{}}

	// STORE: can this code generated, lifted into some higher level api, or expanded on use?????
	// OTHERWISE, we are going to be generating this each and everytime we process code

	// pre-compile the parser statements ( ex. to catch errors. )
	for i := 0; i < model.NumParserAction(); i++ {
		pa := model.ParserActionNum(i)
		actionId, commands := pa.Action, pa.Commands
		if action, ok := model.GetAction(actionId); !ok {
			err = errutil.Append(err, fmt.Errorf("couldnt find action", actionId))
		} else {
			game.log.Println("adding comprehension", actionId, commands)
			if comp, e := p.NewComprehension(actionId,
				func() (parser.IMatch, error) {
					return op.GetObjectMatcher(action)
				}); e != nil {
				err = errutil.Append(err, e)
			} else {
				for _, learn := range commands {
					if _, e := comp.LearnPattern(learn); e != nil {
						err = errutil.Append(err, e)
					}
				}
			}
		}
	}
	if err == nil {
		ret = op
	}
	return ret, err
}

func (op *ObjectParser) GetObjectMatcher(act api.Action) (*ObjectMatcher, error) {
	return NewObjectMatcher(op.game, op.source.FindSource(), act)
}

func (op *ObjectParser) PushParserSource(userSource G.SourceLookup) {
	op.source.PushSource(func() (ret api.Instance) {
		// setup callback context:
		play := NewGameAdapter(op.game)
		// call the user function
		res := userSource(play)
		// unpack the result
		if par, ok := res.(ObjectAdapter); ok {
			ret = par.gobj
		}
		return ret
	})
}

//
func (op *ObjectParser) PopParserSource() {
	op.source.PopSource()
}
