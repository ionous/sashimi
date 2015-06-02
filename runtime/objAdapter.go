package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	"strings"
)

//
// adapts GameObjects for user code
// WARNING: for users to test object equality, the ObjectAdapter must be comparable;
// it can't implement the interface as a pointer, and it cant have any cached values.
//
type ObjectAdapter struct {
	game *Game // for console, Go(), and relations
	gobj *GameObject
}

//
// public for testing
//
func NewObjectAdapter(game *Game, obj *GameObject) ObjectAdapter {
	return ObjectAdapter{game, obj}
}

//
//
//
func (this ObjectAdapter) String() string {
	return this.Name()
}

//
// IObject
//
func (this ObjectAdapter) Name() string {
	return this.gobj.info.Name()
}

//
//
//
func (this ObjectAdapter) Exists() bool {
	return true
}

//
//
//
func (this ObjectAdapter) Class(class string) (okay bool) {
	if cls, err := this.game.Model.Classes.FindClassBySingular(class); err == nil {
		okay = this.gobj.info.CompatibleWith(cls)
	}
	return okay
}

//
//
//
func (this ObjectAdapter) Is(prop string) (ret bool) {
	if prop, index, ok := this.gobj.info.Class().PropertyByChoice(prop); !ok {
		this.logError(fmt.Errorf("is: no such choice '%s'.'%s'", this, prop))
	} else {
		testChoice, _ := prop.IndexToChoice(index)
		currChoice, _ := this.gobj.values.Choice(prop.Id())
		ret = currChoice == testChoice
	}
	return ret
}

//
//
//
func (this ObjectAdapter) SetIs(prop string) {
	if prop, index, ok := this.gobj.info.Class().PropertyByChoice(prop); !ok {
		this.logError(fmt.Errorf("SetIs: no such choice '%s'.'%s'", this, prop))
	} else {
		// get the current choice from the implied property slot
		if currChoice, existed := this.gobj.values.Choice(prop.Id()); !existed {
			err := fmt.Errorf("internal error: choice mismatch via %s for %s %v", prop, prop.Id(), this.gobj.values)
			this.logError(err)
		} else {
			newChoice, _ := prop.IndexToChoice(index)
			if currChoice != newChoice {
				this.gobj.values.removeDirect(currChoice)        // delete the old choice's boolean,
				this.gobj.values.setDirect(newChoice, true)      // and set the new
				this.gobj.values.setDirect(prop.Id(), newChoice) // // set the property slot to the new choice
				this.game.Properties.Notify(this.gobj.Id().String(), prop.Id().String(), currChoice.String(), newChoice.String())
			}
		}
	}
}

//
//
//
func (this ObjectAdapter) Num(prop string) (ret float32) {
	id := M.MakeStringId(prop)
	if v, ok := this.gobj.values.Num(id); ok {
		ret = v
	} else {
		this.logError(TypeMismatch{prop, "get num"})
	}
	return ret
}

//
//
//
func (this ObjectAdapter) SetNum(prop string, value float32) {
	if old, ok := this.gobj.values.SetValue(prop, value); !ok {
		this.logError(TypeMismatch{prop, "set num"})
	} else {
		this.game.Properties.Notify(this.gobj.Id().String(), prop, old, value)
	}
}

//
// returns the evaluated template
// ( note: inform seems to error when trying to store or manipulate templated text )
//
func (this ObjectAdapter) Text(prop string) (ret string) {
	id := M.MakeStringId(prop)
	// is this text stored as a template?
	if temp, ok := this.gobj.temps[id.String()]; ok {
		if s, e := runTemplate(temp, this.gobj.values.data); e != nil {
			this.logError(e)
		} else {
			ret = s
		}
	} else {
		if v, ok := this.gobj.values.Text(id); ok {
			ret = v
		} else {
			this.logError(TypeMismatch{prop, fmt.Sprint(this.gobj.values.data)})
		}
	}
	return ret
}

//
//
//
func (this ObjectAdapter) SetText(prop string, text string) {
	id := M.MakeStringId(prop)
	if e := this.gobj.temps.New(id.String(), text); e != nil {
		this.logError(e)
	} else if old, ok := this.gobj.values.SetValue(prop, text); !ok {
		this.logError(TypeMismatch{prop, "set text"})
	} else {
		this.game.Properties.Notify(this.gobj.Id().String(), prop, old, text)
	}
}

//
//
//
func (this ObjectAdapter) Object(prop string) (ret G.IObject) {
	var res *ObjectAdapter
	if val, ok := this.gobj.info.ValueByName(prop); !ok {
		// TBD: should this be logged? its sure nice to have be able to test objects generically for properties
		// this.logError(fmt.Errorf("object requested, but no such property %s", prop))
	} else {
		if rel, ok := val.(*M.RelativeValue); !ok {
			this.logError(fmt.Errorf("object requested, but property is %T", val))
		} else {
			if rel.RelativeProperty().ToMany() {
				this.logError(fmt.Errorf("object requested, but relation is list"))
			} else {
				list := rel.List()
				if len(list) != 0 {
					if gobj, ok := this.game.FindObject(list[0]); ok {
						res = &ObjectAdapter{this.game, gobj}
					}
				}
			}
		}
	}
	if res != nil {
		ret = *res
	} else {
		ret = this.game.nullobj
	}
	return ret
}

//
//
//
func (this ObjectAdapter) ObjectList(prop string) (ret []G.IObject) {
	if val, ok := this.gobj.info.ValueByName(prop); !ok {
		this.logError(fmt.Errorf("object list requested, but no such property"))
	} else {
		if rel, ok := val.(*M.RelativeValue); !ok {
			this.logError(fmt.Errorf("object list requested, but property is %T", val))
		} else {
			if !rel.RelativeProperty().ToMany() {
				this.logError(fmt.Errorf("object list requested, but relation is singular"))
			} else {
				list := rel.List()
				ret = make([]G.IObject, len(list))
				for i, objName := range list {
					if gobj, ok := this.game.FindObject(objName); ok {
						ret[i] = ObjectAdapter{this.game, gobj}
					} else {
						ret[i] = this.game.nullobj
					}
				}
			}
		}
	}
	return ret
}

//
//
//
func (this ObjectAdapter) SetObject(prop string, other G.IObject) {
	if val, ok := this.gobj.info.ValueByName(prop); !ok {
		this.logError(fmt.Errorf("SetObject: no such relation '%s'.'%s'", this, prop))
	} else {
		if rel, ok := val.(*M.RelativeValue); !ok {
			this.logError(TypeMismatch{prop, "SetObject"})
		} else {
			// if the object doesnt exist, we take it to mean they are clearing the reference.
			if other, ok := other.(ObjectAdapter); !ok {
				this.game.log.Println("clearing", this.Name(), prop)
				if removed, e := rel.ClearReference(); e != nil {
					this.logError(e)
				} else {
					this.game.Properties.Notify(this.gobj.Id().String(), prop, removed, "")
				}
			} else {
				// FIX? the impedence b/t IObject and Reference is annoying.
				other := other.gobj.Id().String()
				if ref, e := this.game.references.FindByName(other); e != nil {
					this.logError(e)
				} else if removed, e := rel.SetReference(ref); e != nil {
					this.logError(e)
				} else {
					// removed is probably a single object
					this.game.Properties.Notify(this.gobj.Id().String(), prop, removed, other)
				}
			}
		}
	}
}

//
//
//
func (this ObjectAdapter) Says(text string) {
	// FIX: share some template love with GameEventAdapter.Say()
	lines := strings.Split(text, "\n")
	this.game.output.ActorSays(this.Name(), lines)
}

//
// send all the events associated with the propd action; and,
// run the default action if appropriate
// @see also: Game.ProcessEventQueue
//
func (this ObjectAdapter) Go(prop string, objects ...G.IObject) {
	if action, e := this.game.Model.Actions.FindActionByName(prop); e != nil {
		this.logError(e)
	} else {
		// ugly: we need the props, even tho we already have the objects...
		nouns := make([]string, len(objects)+1)
		nouns[0] = this.Name()
		for i, o := range objects {
			nouns[i+1] = o.Name()
		}
		if act, e := this.game.newRuntimeAction(action, nouns...); e != nil {
			this.logError(e)
		} else {
			tgt := ObjectTarget{this.game, this.gobj}
			msg := E.Message{Name: action.Event(), Data: act}
			// see ProcessEventQueue()
			path := E.NewPathTo(tgt)
			//this.game.log.Output(3, fmt.Sprintf("go %s %s", prop, path))
			if runDefault, err := msg.Send(path); err != nil {
				this.logError(err)
			} else if runDefault {
				act.runDefaultActions()
			}
		}
	}
}

//
//
//
func (this ObjectAdapter) logError(err error) (hadError bool) {
	if err != nil {
		this.game.log.Output(4, fmt.Sprint("!!!Error:", err.Error()))
		hadError = true
	}
	return hadError
}
