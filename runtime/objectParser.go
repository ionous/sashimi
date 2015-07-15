package runtime

import (
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	P "github.com/ionous/sashimi/parser"
	"github.com/ionous/sashimi/util/errutil"
)

//
// ObjectParser adapts the parser to the runtime,
// queuing actions for game objects based upon user input.
//
type ObjectParser struct {
	game *Game // for access to all available objects
	*P.Parser
	source *ParserSourceStack
}

//
// NormalizeInput wraps parser.NormalizeInput for convenience's sake.
//
func (op *ObjectParser) NormalizeInput(s string) string {
	return P.NormalizeInput(s)
}

//
// NewParser and add all commands and patterns.
//
func NewParser(game *Game,
) (ret *ObjectParser, err error,
) {
	model := game.Model
	parser := P.NewParser()
	op := &ObjectParser{game, parser, &ParserSourceStack{}}
	// pre-compile the parser statements ( ex. to catch errors. )
	for _, p := range model.ParserActions {
		act := p.Action()
		name := act.Action()
		if comp, e := parser.NewComprehension(name,
			func() (P.IMatch, error) {
				return op.NewObjectMatcher(act)
			}); e != nil {
			err = errutil.Append(err, e)
		} else {
			for _, learn := range p.Commands() {
				if _, e := comp.LearnPattern(learn); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
	}
	if err == nil {
		ret = op
	}
	return ret, err
}

func (op *ObjectParser) NewObjectMatcher(act *M.ActionInfo) (*ObjectMatcher, error) {
	return NewObjectMatcher(op.game, op.source.FindSource(), act)
}

func (op *ObjectParser) PushParserSource(userSource G.SourceLookup) {
	op.source.PushSource(func() (ret *GameObject) {
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
