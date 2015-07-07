package runtime

import (
	P "github.com/ionous/sashimi/parser"
	"github.com/ionous/sashimi/util/errutil"
)

//
// adapt the parser to the model
//
type ModelParser struct {
	*P.Parser
}

func (this *ModelParser) NormalizeInput(s string) string {
	return P.NormalizeInput(s)
}

//
// create the parser, add all commands and patterns
//
func NewParser(game *Game,
) (ret *ModelParser, err error,
) {
	model := game.Model
	parser := P.NewParser()
	// pre-compile the parser statements ( ex. to catch errors. )
	for _, p := range model.ParserActions {
		act := p.Action()
		if comp, e := parser.NewComprehension(
			act.Action(),
			func() (P.IMatch, error) {
				return NewObjectMatcher(game, act)
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
		ret = &ModelParser{parser}
	}
	return ret, err
}
