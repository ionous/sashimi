package runtime

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	P "github.com/ionous/sashimi/parser"
	"github.com/ionous/sashimi/util/ident"
)

//
// Return a noun which matches an instance's string id
//
type ObjectMatcher struct {
	game    *Game
	act     *M.ActionInfo
	objects []*GameObject
}

// make sure the source class matches
func NewObjectMatcher(game *Game, source *GameObject, act *M.ActionInfo,
) (
	ret *ObjectMatcher,
	err error,
) {
	if source == nil {
		err = fmt.Errorf("couldnt find command source for %s", act)
	} else if !source.Class().CompatibleWith(act.Source().Id) {
		err = fmt.Errorf("source class for %s doesnt match", act)
	} else {
		om := &ObjectMatcher{game, act, nil}
		om.addObject(source)
		ret = om
	}
	return ret, err
}

//
// MatchNOun to relate the passed name and article to internal objects.
//
func (om *ObjectMatcher) MatchNoun(name string, _ string) (err error) {
	nouns := om.act.NounTypes
	if cnt, max := len(om.objects), len(nouns); cnt >= max {
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
	nouns := om.act.NounTypes
	if cnt, max := len(om.objects), len(nouns); cnt != max {
		err = P.MismatchedNouns("I", max, cnt)
	} else {
		tgt := ObjectTarget{om.game, om.objects[0]}
		act := &RuntimeAction{om.game, om.act, om.objects, nil}
		om.game.queue.QueueEvent(tgt, om.act.EventName, act)
	}
	return err
}

//
// MatchId is usually called by MatchNoun, public for net sessions :(
//
func (om *ObjectMatcher) MatchId(id ident.Id) (okay bool) {
	if gobj, ok := om.game.Objects[id]; ok {
		nouns := om.act.NounTypes
		if cnt, max := len(om.objects), len(nouns); cnt < max {
			class := nouns[cnt]
			if gobj.Class().CompatibleWith(class.Id) {
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
