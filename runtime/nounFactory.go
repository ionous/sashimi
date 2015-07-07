package runtime

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	P "github.com/ionous/sashimi/parser"
)

//
// Return a noun which matches an instance's string id
//
type ObjectMatcher struct {
	game    *Game
	act     *M.ActionInfo
	objects []*GameObject
	values  map[string]TemplateValues
}

// make sure the source class matches
func NewObjectMatcher(game *Game, source *GameObject, act *M.ActionInfo,
) (
	ret *ObjectMatcher,
	err error,
) {
	if source == nil {
		err = fmt.Errorf("couldnt find command source for %s", act)
	} else if !source.inst.Class().CompatibleWith(act.Source().Id()) {
		err = fmt.Errorf("source class for %s doesnt match", act)
	} else {
		om := &ObjectMatcher{game, act, nil, make(map[string]TemplateValues)}
		om.addObject(source)
		ret = om
	}
	return ret, err
}

//
// MatchNOun to relate the passed name and article to internal objects.
//
func (om *ObjectMatcher) MatchNoun(name string, _ string) (err error) {
	nouns := om.act.NounSlice()
	if cnt, max := len(om.objects), len(nouns); cnt >= max {
		err = fmt.Errorf("You've told me more than I've understood.")
	} else {
		tried, ok := om.game.Model.NounNames.Try(name, func(id M.StringId) (okay bool) {
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
func (om *ObjectMatcher) Matched() (err error) {
	nouns := om.act.NounSlice()
	if cnt, max := len(om.objects), len(nouns); cnt != max {
		err = P.MismatchedNouns("I", max, cnt)
	} else {
		tgt := ObjectTarget{om.game, om.objects[0]}
		act := &RuntimeAction{om.game, om.act, om.objects, om.values, nil}
		om.game.queue.QueueEvent(tgt, om.act.Event(), act)
	}
	return err
}

//
// MatchId is usually called by MatchNoun, public for net sessions :(
//
func (om *ObjectMatcher) MatchId(id M.StringId) (okay bool) {
	if obj, ok := om.game.Objects[id]; ok {
		nouns := om.act.NounSlice()
		if cnt, max := len(om.objects), len(nouns); cnt < max {
			class := nouns[cnt]
			if obj.inst.Class().CompatibleWith(class.Id()) {
				om.addObject(obj)
				okay = true
			}
		}
	}
	return okay
}

func (om *ObjectMatcher) addObject(gobj *GameObject) {
	cnt := len(om.objects)
	keys := []string{"Source", "Target", "Context"}
	om.objects = append(om.objects, gobj)
	key := keys[cnt]
	om.values[key] = gobj.data
}
