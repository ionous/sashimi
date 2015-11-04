package runtime

import (
	"fmt"
	"github.com/ionous/sashimi/parser"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

// ObjectMatcher returns nouns which matches an instance's string id
type ObjectMatcher struct {
	game    *Game
	act     api.Action
	objects []*GameObject
}

// make sure the source class matches
func NewObjectMatcher(game *Game, source *GameObject, act api.Action,
) (
	ret *ObjectMatcher,
	err error,
) {
	if source == nil {
		err = fmt.Errorf("couldnt find command source for %s", act)
	} else {
		nouns := act.GetNouns()
		if !source.Class().CompatibleWith(nouns.Get(api.SourceNoun)) {
			err = fmt.Errorf("source class for %s doesnt match", act)
		} else {
			om := &ObjectMatcher{game, act, nil}
			om.addObject(source)
			ret = om
		}
	}
	return ret, err
}

//
// MatchNOun to relate the passed name and article to internal objects.
//
func (om *ObjectMatcher) MatchNoun(name string, _ string) (err error) {
	nouns := om.act.GetNouns()
	if cnt, max := len(om.objects), nouns.GetNounCount(); cnt >= max {
		err = fmt.Errorf("You've told me more than I've understood.")
	} else {
		tried, ok := om.game.Model.NounNames.Try(name, func(id ident.Id) (okay bool) {
			return om.MatchId(id)
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

//
// Matches gets called by the parser after succesfully found the command and nouns.
//
func (om *ObjectMatcher) OnMatch() (err error) {
	nouns := om.act.GetNouns()
	if cnt, max := len(om.objects), nouns.GetNounCount(); cnt != max {
		err = parser.MismatchedNouns("I", max, cnt)
	} else {
		tgt := ObjectTarget{om.game, om.objects[0]}
		act := &RuntimeAction{om.game, om.act.GetId(), om.objects, nil}
		om.game.queue.QueueEvent(tgt, om.act.GetEventName(), act)
	}
	return err
}

//
// MatchId is usually called by MatchNoun, public for net sessions :(
//
func (om *ObjectMatcher) MatchId(id ident.Id) (okay bool) {
	if gobj, ok := om.game.Objects[id]; ok {
		nouns := om.act.GetNouns()
		if cnt, max := len(om.objects), nouns.GetNounCount(); cnt < max {
			if gobj.Class().CompatibleWith(nouns[cnt]) {
				om.addObject(gobj)
				okay = true
			}
		}
	}
	return okay
}

func (om *ObjectMatcher) addObject(gobj *GameObject) {
	om.objects = append(om.objects, gobj)
	// cnt := len(om.objects)-1
	// keys := []string{"Source", "Target", "Context"}
	// key := keys[cnt]
	// om.values[key] = gobj.vals
}
