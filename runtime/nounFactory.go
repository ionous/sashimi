package runtime

import (
	"fmt"
	"github.com/ionous/sashimi/parser"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

// ObjectMatcher returns nouns which matches an instance's string id
type ObjectMatcher struct {
	mdl     api.Model
	nouns   api.Nouns
	objects []api.Instance
	onMatch OnMatch
}

type OnMatch func([]api.Instance)

// make sure the source class matches
func NewObjectMatcher(mdl api.Model, nouns api.Nouns, onMatch OnMatch) *ObjectMatcher {
	return &ObjectMatcher{mdl, nouns, nil, onMatch}
}

// MatchNoun to relate the passed name and article to internal objects.
func (om *ObjectMatcher) MatchNoun(name string, _ string) (err error) {
	if cnt, max := len(om.objects), om.nouns.GetNounCount(); cnt >= max {
		err = fmt.Errorf("You've told me more than I've understood.")
	} else {
		tried, ok := om.mdl.MatchNounName(name, func(id ident.Id) bool {
			return om.AddObject(id) == nil
		})
		if !ok {
			if tried > 0 {
				err = fmt.Errorf("I don't know how to use that for this.")
			} else {
				err = fmt.Errorf("I don't see any such thing.")
			}
		}
	}
	return err
}

// OnMatch called by the parser after succesfully found the command and nouns.
func (om *ObjectMatcher) OnMatch() (err error) {
	objects, nouns := om.objects, om.nouns
	if cnt, max := len(objects), nouns.GetNounCount(); cnt != max {
		err = parser.MismatchedNouns("I", max, cnt)
	} else {
		om.onMatch(objects)
	}
	return
}

// AddObject is usually called by MatchNoun, public for net sessions :(
func (om *ObjectMatcher) AddObject(id ident.Id) (err error) {
	if gobj, ok := om.mdl.GetInstance(id); !ok {
		err = fmt.Errorf("couldnt find noun %s", id)
	} else if cnt, max := len(om.objects), om.nouns.GetNounCount(); !(cnt < max) {
		err = fmt.Errorf("too many nouns %d<%d", cnt, max)
	} else {
		noun := om.nouns[cnt]
		if cls := gobj.GetParentClass(); !om.mdl.AreCompatible(cls.GetId(), noun) {
			err = fmt.Errorf("noun %d not compatible %s(%s) != %s", cnt, id, cls, noun)
		} else {
			om.objects = append(om.objects, gobj)
		}
	}
	return
}
