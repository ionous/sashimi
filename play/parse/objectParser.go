package parse

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/parser"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

type MatchMaker struct {
	mdl meta.Model
	src ident.Id
}

func (m MatchMaker) NewMatcher(id ident.Id) (ret parser.IMatch, err error) {
	if act, ok := m.mdl.GetAction(id); !ok {
		err = errutil.New("couldnt find action", id)
	} else {
		ret, err = NewObjectMatcher(act, m.src, m.mdl)
	}
	return
}

// NewObjectParser and add all commands and patterns.
// FIX-STORE: generate literals for all comprehensions; pre-compile to check for errors.
func NewObjectParser(mdl meta.Model, source ident.Id) (p parser.P, err error) {
	if _, ok := mdl.GetInstance(source); !ok {
		err = errutil.New("couldnt find source", source)
	} else {
		p = parser.NewParser(MatchMaker{mdl, source})
		for i := 0; i < mdl.NumParserAction(); i++ {
			pa := mdl.ParserActionNum(i)
			// lookup the parser actions to catch any strange compiler errors
			if _, ok := mdl.GetAction(pa.Action); !ok {
				err = errutil.Append(err, errutil.New("couldnt find action", pa.Action))
			} else {
				if comp, e := p.NewComprehension(pa.Action); e != nil {
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
	}
	return
}
