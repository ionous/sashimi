package parse

import (
	"fmt"
	"github.com/ionous/sashimi/parser"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

// ObjectMatcher returns nouns which matches an instance's string id
type ObjectMatcher struct {
	mdl     api.Model
	act     api.Action
	objects []api.Instance
}

// make sure the source class matches
func NewObjectMatcher(act api.Action, src ident.Id, mdl api.Model) (ret *ObjectMatcher, err error) {
	om := &ObjectMatcher{mdl, act, nil}
	if e := om.AddObject(src); e != nil {
		err = e
	} else {
		ret = om
	}
	return
}

// MatchNoun to relate the passed name and article to internal objects.
func (om *ObjectMatcher) MatchNoun(name string, _ string) (err error) {
	nouns := om.act.GetNouns()
	if cnt, max := len(om.objects), nouns.GetNounCount(); cnt >= max {
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

func (om *ObjectMatcher) GetMatch() (act api.Action, ret []api.Instance, err error) {
	objects, nouns := om.objects, om.act.GetNouns()
	if cnt, max := len(objects), nouns.GetNounCount(); cnt != max {
		err = parser.MismatchedNouns("I", max, cnt)
	} else {
		act, ret = om.act, objects
	}
	return
}

// AddObject is usually called by MatchNoun, public for net sessions :(
func (om *ObjectMatcher) AddObject(id ident.Id) (err error) {
	nouns := om.act.GetNouns()
	if gobj, ok := om.mdl.GetInstance(id); !ok {
		err = fmt.Errorf("couldnt find noun %s", id)
		panic(id.Split())
	} else if cnt, max := len(om.objects), nouns.GetNounCount(); !(cnt < max) {
		err = fmt.Errorf("too many nouns %d<%d", cnt, max)
	} else {
		noun := nouns[cnt]
		if cls := gobj.GetParentClass(); !om.mdl.AreCompatible(cls.GetId(), noun) {
			err = fmt.Errorf("noun %d not compatible %s(%s) != %s", cnt, id, cls, noun)
		} else {
			om.objects = append(om.objects, gobj)
		}
	}
	return
}
