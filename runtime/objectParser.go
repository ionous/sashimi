package runtime

import (
	"fmt"
	"github.com/ionous/sashimi/parser"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/errutil"
)

type IQueueAction interface {
	QueueAction(api.Action, []api.Instance)
}

// NewObjectParser and add all commands and patterns.
// STORE: can this code generated, lifted into some higher level api, or expanded on use?????
// OTHERWISE, we are going to be generating this each and everytime we process code
func NewObjectParser(queue IQueueAction, mdl api.Model, source string) (p parser.P, err error) {
	// pre-compile the parser statements ( ex. to catch errors. )
	sourceId := MakeStringId(source)
	p = make(parser.P)
	for i := 0; i < mdl.NumParserAction(); i++ {
		pa := mdl.ParserActionNum(i)
		if action, ok := mdl.GetAction(pa.Action); !ok {
			err = errutil.Append(err, fmt.Errorf("couldnt find action", pa.Action))
		} else {
			// PROBABLY WANT TO GENERATE COMPS AHEAD OF TIME, and paramterize the callback with something.
			if comp, e := p.NewComprehension(pa.Action,
				func() (ret parser.IMatch, err error) {
					om := NewObjectMatcher(mdl, action.GetNouns(), func(objects []api.Instance) {
						queue.QueueAction(action, objects)
					})
					if e := om.AddObject(sourceId); e != nil {
						err = fmt.Errorf("couldnt add source, error: %s", e)
					} else {
						ret = om
					}
					return
				}); e != nil {
				err = errutil.Append(err, e)
			} else {
				for _, learn := range pa.Commands {
					if _, e := comp.LearnPattern(learn); e != nil {
						err = errutil.Append(err, e)
					}
				}
			}
		}
	}
	return
}
