package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
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
func (this ObjectAdapter) Is(name string) (ret bool) {
	if prop, index, ok := this.gobj.info.Class().PropertyByChoice(name); !ok {
		this.logError(fmt.Errorf("is: no such choice '%s'.'%s'", this, name))
	} else {
		testChoice, _ := prop.IndexToChoice(index)
		currChoice, _ := this.gobj.values.getChoice(prop.Id())
		ret = currChoice == testChoice
	}
	return ret
}

//
//
//
func (this ObjectAdapter) SetIs(name string) {
	if prop, index, ok := this.gobj.info.Class().PropertyByChoice(name); !ok {
		this.logError(fmt.Errorf("SetIs: no such choice '%s'.'%s'", this, name))
	} else {
		// get the current choice from the implied property slot
		if currChoice, existed := this.gobj.values.getChoice(prop.Id()); !existed {
			err := fmt.Errorf("internal error: choice mismatch via %s for %s %v", name, prop.Id(), this.gobj.values)
			this.logError(err)
		} else {
			newChoice, _ := prop.IndexToChoice(index)
			if currChoice != newChoice {
				this.gobj.values.remove(currChoice)        // delete the old choice's boolean,
				this.gobj.values.set(newChoice, true)      // and set the new
				this.gobj.values.set(prop.Id(), newChoice) // // set the property slot to the new choice
			}
		}
	}
}

//
//
//
func (this ObjectAdapter) Num(name string) (ret float32) {
	id := M.MakeStringId(name)
	if v, okay := this.gobj.values.getNum(id); okay {
		ret = v
	} else {
		this.logError(TypeMismatch{name, "get num"})
	}
	return ret
}

//
//
//
func (this ObjectAdapter) SetNum(name string, value float32) {
	if !this.gobj.values.safeSet(name, value) {
		this.logError(TypeMismatch{name, "set num"})
	}
}

//
// returns the evaluated template
// ( note: inform seems to error when trying to store or manipulate templated text )
//
func (this ObjectAdapter) Text(name string) (ret string) {
	id := M.MakeStringId(name)
	// is this text stored as a template?
	if temp, okay := this.gobj.temps[id.String()]; okay {
		if s, e := runTemplate(temp, this.gobj.values.data); e != nil {
			this.logError(e)
		} else {
			ret = s
		}
	} else {
		if v, okay := this.gobj.values.getText(id); okay {
			ret = v
		} else {
			this.logError(TypeMismatch{name, fmt.Sprint(this.gobj.values.data)})
		}
	}
	return ret
}

//
//
//
func (this ObjectAdapter) SetText(name string, text string) {
	id := M.MakeStringId(name)
	if e := this.gobj.temps.New(id.String(), text); e != nil {
		this.logError(e)
	} else if !this.gobj.values.safeSet(name, text) {
		this.logError(TypeMismatch{name, "set text"})
	}
}

//
//
//
func (this ObjectAdapter) Object(name string) (ret G.IObject) {
	var res *ObjectAdapter
	if val, ok := this.gobj.info.ValueByName(name); !ok {
		// TBD: should this be logged? its sure nice to have be able to test objects generically for properties
		// this.logError(fmt.Errorf("object requested, but no such property %s", name))
	} else {
		if rel, ok := val.(*M.RelativeValue); !ok {
			this.logError(fmt.Errorf("object requested, but property is %T", val))
		} else {
			if rel.IsMany() {
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
func (this ObjectAdapter) ObjectList(name string) (ret []G.IObject) {
	if val, ok := this.gobj.info.ValueByName(name); !ok {
		this.logError(fmt.Errorf("object list requested, but no such property"))
	} else {
		if rel, ok := val.(*M.RelativeValue); !ok {
			this.logError(fmt.Errorf("object list requested, but property is %T", val))
		} else {
			if !rel.IsMany() {
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
func (this ObjectAdapter) SetObject(name string, other G.IObject) {
	if val, ok := this.gobj.info.ValueByName(name); !ok {
		this.logError(fmt.Errorf("setobject: no such choice '%s'.'%s'", this, name))
	} else {
		if rel, ok := val.(*M.RelativeValue); !ok {
			this.logError(TypeMismatch{name, "set object"})
		} else {
			if obj, ok := other.(ObjectAdapter); !ok {
				this.game.log.Println("clearing", this.Name(), name)
				if e := rel.ClearReference(); e != nil {
					this.logError(e)
				}
			} else {
				// FIX? the impedence b/t IObject and Reference is annoying.
				if ref, e := this.game.references.FindByName(obj.Name()); e != nil {
					this.logError(e)
				} else {
					if e := rel.SetReference(ref); e != nil {
						this.logError(e)
					}
				}
			}
		}
	}
}

//
//
//
func (this ObjectAdapter) Says(s string) {
	this.game.console.Println(this.Name(), ": ", s)
	this.game.console.Println()
}

//
// send all the events associated with the named action; and,
// run the default action if appropriate
// @see also: Game.ProcessEventQueue
//
func (this ObjectAdapter) Go(name string, objects ...G.IObject) {
	if action, e := this.game.Model.Actions.FindActionByName(name); e != nil {
		this.logError(e)
	} else {
		// ugly: we need the names, even tho we already have the objects...
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
			this.game.log.Output(3, fmt.Sprintf("go %s %s", name, path))
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
