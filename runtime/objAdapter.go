package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	"strings"
)

//
// Adapts GameObjects for user script callbacks.
// WARNING: for users to test object equality, the ObjectAdapter must be comparable;
// it can't implement the interface as a pointer, and it cant have any cached values.
//
type ObjectAdapter struct {
	game *Game // for console, Go(), and relations
	gobj *GameObject
}

//
// Public for testing.
//
func NewObjectAdapter(game *Game, obj *GameObject) ObjectAdapter {
	return ObjectAdapter{game, obj}
}

//
// Helper for debugging.
//
func (adapt ObjectAdapter) String() string {
	return adapt.Name()
}

//
// Name of the object.
//
func (adapt ObjectAdapter) Name() string {
	return adapt.gobj.inst.Name()
}

//
// Is adapt a valid object?
//
func (adapt ObjectAdapter) Exists() bool {
	return true
}

//
// Is adapt object based on the passed class in any fashion. ( parent or other ancestor )
//
func (adapt ObjectAdapter) Class(class string) (okay bool) {
	if cls, ok := adapt.game.Model.Classes.FindClassBySingular(class); ok {
		okay = adapt.gobj.inst.Class().CompatibleWith(cls.Id())
	}
	return okay
}

//
// Is adapt object in the passed state?
//
func (adapt ObjectAdapter) Is(state string) (ret bool) {
	if prop, index, ok := adapt.gobj.inst.Class().PropertyByChoice(state); !ok {
		adapt.logError(fmt.Errorf("is: no such choice '%s'.'%s'", adapt, state))
	} else {
		testChoice, _ := prop.IndexToChoice(index)
		currChoice := adapt.gobj.GetValue(prop.Id())
		ret = currChoice == testChoice
	}
	return ret
}

//
// Change the state of an object.
//
func (adapt ObjectAdapter) SetIs(state string) {
	if prop, index, ok := adapt.gobj.inst.Class().PropertyByChoice(state); !ok {
		adapt.logError(fmt.Errorf("SetIs: no such choice '%s'.'%s'", adapt, state))
	} else {
		// get the current choice from the implied property slot
		if currChoice, ok := adapt.gobj.GetValue(prop.Id()).(M.StringId); !ok {
			err := TypeMismatch(adapt.gobj.Id().String(), prop.Id().String())
			adapt.logError(err)
		} else {
			newChoice, _ := prop.IndexToChoice(index)
			if currChoice != newChoice {
				adapt.gobj.removeDirect(currChoice)        // delete the old choice's boolean,
				adapt.gobj.setDirect(newChoice, true)      // and set the new
				adapt.gobj.setDirect(prop.Id(), newChoice) // // set the property slot to the new choice
				adapt.game.Properties.Notify(adapt.gobj.Id(), prop.Id(), currChoice, newChoice)
			}
		}
	}
}

//
// Return the value of the passed number property.
//
func (adapt ObjectAdapter) Num(prop string) (ret float32) {
	id := M.MakeStringId(prop)
	if val, ok := adapt.gobj.GetValue(id).(float32); !ok {
		adapt.logError(TypeMismatch(prop, "get num"))
	} else {
		ret = val
	}
	return ret
}

//
// Change the value of an existing number property.
//
func (adapt ObjectAdapter) SetNum(prop string, value float32) {
	id := M.MakeStringId(prop)
	if old, ok := adapt.gobj.SetValue(id, value); !ok {
		adapt.logError(TypeMismatch(prop, "set num"))
	} else {
		adapt.game.Properties.Notify(adapt.gobj.Id(), id, old, value)
	}
}

//
// Return the value of the passed text property ( expanding any templated text. )
// ( interestingly, inform seems to error when trying to store or manipulate templated text. )
//
func (adapt ObjectAdapter) Text(prop string) (ret string) {
	id := M.MakeStringId(prop)
	// is adapt text stored as a template?
	if temp, ok := adapt.gobj.temps[id.String()]; ok {
		if s, e := runTemplate(temp, adapt.gobj.data); e != nil {
			adapt.logError(e)
		} else {
			ret = s
		}
	} else if val, ok := adapt.gobj.GetValue(id).(string); !ok {
		adapt.logError(TypeMismatch(prop, "get text"))
	} else {
		ret = val
	}
	return ret
}

//
// Change the value of an existing text property.
//
func (adapt ObjectAdapter) SetText(prop string, text string) {
	id := M.MakeStringId(prop)
	if e := adapt.gobj.temps.New(id.String(), text); e != nil {
		adapt.logError(e)
	} else if old, ok := adapt.gobj.SetValue(id, text); !ok {
		adapt.logError(TypeMismatch(prop, "set text"))
	} else {
		adapt.game.Properties.Notify(adapt.gobj.Id(), id, old, text)
	}
}

//
// Return a related object.
//
func (adapt ObjectAdapter) Object(prop string) (ret G.IObject) {
	// TBD: should these be logged? its sure nice to have be able to test objects generically for properties
	var res M.StringId
	if p, ok := adapt.gobj.inst.Class().FindProperty(prop); ok {
		switch p := p.(type) {
		case *M.PointerProperty:
			if val, ok := adapt.gobj.GetValue(p.Id()).(M.StringId); ok {
				res = val
			}
		case *M.RelativeProperty:
			// TBD: can the relative property changes automatically reflect into the value table
			// ex. on event?
			rel := adapt.gobj.inst.GetRelativeValue(p)
			if rel.GetRelativeProperty().ToMany() {
				adapt.logError(fmt.Errorf("object requested, but relation is list"))
			} else {
				list := rel.List()
				if len(list) != 0 {
					res = M.StringId(list[0])
				}
			}

		}
	}
	if gobj, ok := adapt.game.Objects[res]; ok {
		ret = ObjectAdapter{adapt.game, gobj}
	} else {
		ret = adapt.game.nullobj
	}
	return ret
}

//
// Changes a relationship.
//
func (adapt ObjectAdapter) SetObject(prop string, object G.IObject) {
	if p, ok := adapt.gobj.inst.Class().FindProperty(prop); ok {
		switch p := p.(type) {
		default:
			adapt.logError(TypeMismatch(adapt.Name(), prop))
		case *M.PointerProperty:
			set := false
			if other, ok := object.(ObjectAdapter); !ok {
				adapt.gobj.SetValue(p.Id(), M.StringId(""))
				set = true
			} else {
				adapt.gobj.SetValue(p.Id(), other.gobj.Id())
				set = true
			}
			if !set {
				adapt.logError(fmt.Errorf("couldnt set value for prop %s", prop))
			}
		case *M.RelativeProperty:
			rel := adapt.gobj.inst.GetRelativeValue(p)

			// if the referenced object doesnt exist, we take it to mean they are clearing the reference.
			if other, ok := object.(ObjectAdapter); !ok {
				if removed, e := rel.ClearReference(); e != nil {
					adapt.logError(e)
				} else {
					adapt.game.Properties.Notify(adapt.gobj.Id(), rel.Property().Id(), removed, "")
				}
			} else {
				// FIX? the impedence b/t IObject and Reference is annoying.
				other := other.gobj.Id()
				if ref, ok := adapt.game.Model.Instances[other]; !ok {
					adapt.logError(fmt.Errorf("SetObject: couldnt find object names %s", other))
				} else if removed, e := rel.SetReference(ref); e != nil {
					adapt.logError(e)
				} else {
					// removed is probably a single object
					adapt.game.Properties.Notify(adapt.gobj.Id(), rel.Property().Id(), removed, other.String())
				}
			}
		}
	}
}

//
// Return a list of related objects.
//
func (adapt ObjectAdapter) ObjectList(prop string) (ret []G.IObject) {
	if rel, ok := adapt.gobj.inst.FindRelativeValue(prop); !ok {
		adapt.logError(fmt.Errorf("object list requested, but no such property"))
	} else {
		if !rel.GetRelativeProperty().ToMany() {
			adapt.logError(fmt.Errorf("object list requested, but relation is singular"))
		} else {
			list := rel.List()
			ret = make([]G.IObject, len(list))
			for i, objName := range list {
				if gobj, ok := adapt.game.FindObject(objName); ok {
					ret[i] = ObjectAdapter{adapt.game, gobj}
				} else {
					ret[i] = adapt.game.nullobj
				}
			}
		}
	}
	return ret
}

//
// This actor has something to say.
//
func (adapt ObjectAdapter) Says(text string) {
	// FIX: share some template love with GameEventAdapter.Say()
	lines := strings.Split(text, "\n")
	adapt.game.output.ActorSays(adapt.gobj, lines)
}

//
// Send all the events associated with the action; and,
// run the default action if appropriate
// @see also: Game.ProcessEventQueue
//
func (adapt ObjectAdapter) Go(act string, objects ...G.IObject) {
	if action, ok := adapt.game.Model.Actions.FindActionByName(act); !ok {
		adapt.logError(fmt.Errorf("unknown action for Go %s", act))
	} else {
		// ugly: we need the props, even tho we already have the objects...
		nouns := make([]string, len(objects)+1)
		nouns[0] = adapt.Name()
		for i, o := range objects {
			nouns[i+1] = o.Name()
		}
		if act, e := adapt.game.newRuntimeAction(action, nouns...); e != nil {
			adapt.logError(e)
		} else {
			tgt := ObjectTarget{adapt.game, adapt.gobj}
			msg := E.Message{Name: action.Event(), Data: act}
			// see ProcessEventQueue()
			path := E.NewPathTo(tgt)
			//adapt.game.log.Output(3, fmt.Sprintf("go %s %s", prop, path))
			if runDefault, err := msg.Send(path); err != nil {
				adapt.logError(err)
			} else if runDefault {
				act.runDefaultActions()
			}
		}
	}
}

//
//
//
func (adapt ObjectAdapter) logError(err error) (hadError bool) {
	if err != nil {
		adapt.game.log.Output(4, fmt.Sprint("!!!Error:", err.Error()))
		hadError = true
	}
	return hadError
}
