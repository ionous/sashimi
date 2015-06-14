package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	"log"
)

// FIX: standarize member exports by splitting game into smaller classes and interfaces
type Game struct {
	Model          *M.Model
	Objects        GameObjects
	Dispatchers    ClassDispatchers
	output         IOutput
	parser         *ModelParser
	queue          *E.Queue
	nullobj        *NullObject
	defaultActions DefaultActions
	SystemActions  SystemActions
	log            *log.Logger
	references     M.References
	Properties     PropertyWatchers
	parentLookup   ParentLookupStack
	parserSource   ParserSourceStack
}

type logAdapter struct {
	output IOutput
}

func (this logAdapter) Write(p []byte) (n int, err error) {
	this.output.Log(string(p))
	return len(p), nil
}

// each action can have a chain of default actions
type DefaultActions map[*M.ActionInfo][]G.Callback

func NewGame(model *M.Model, output IOutput) (game *Game, err error) {
	log := log.New(logAdapter{output}, "game: ", log.Lshortfile)
	dispatchers := NewDispatchers(log)
	objects, e := CreateGameObjects(model.Instances)
	if e != nil {
		return nil, e
	}

	game = &Game{
		model,
		objects,
		dispatchers,
		output, nil,
		E.NewQueue(),
		&NullObject{log},
		make(DefaultActions),
		NewSystemActions(),
		log,
		M.NewReferences(model.Classes, model.Instances, model.Tables),
		PropertyWatchers{},
		ParentLookupStack{},
		ParserSourceStack{},
	}

	parser, e := NewParser(game)
	if e != nil {
		return nil, e
	}
	game.parser = parser

	for _, handler := range model.ActionHandlers {
		act, cb := handler.Action(), handler.Callback()
		arr := game.defaultActions[act]
		// FIX: for now treating target as bubble,
		// really the compiler should hand off a sorted flat list based on three separate groups
		// target growing in the same direction as after, but distinctly in the middle of things.
		if !handler.UseCapture() {
			arr = append(arr, cb)
		} else {
			// prepend:
			arr = append([]G.Callback{cb}, arr...)
		}
		game.defaultActions[act] = arr
	}

	// FUTURE: move into scenes, with a handle for removal
	for _, listener := range model.EventListeners {
		act := listener.Action()
		callback := GameCallback{game, listener}
		//log.Printf("creating listener %s", listener.String())

		if inst := listener.Instance(); inst != nil {
			obj := objects[inst.Id()]
			obj.dispatcher.Listen(act.Event(), callback, listener.UseCapture())
		} else if cls := listener.Class(); cls != nil {
			dispatch := dispatchers.CreateDispatcher(cls)
			dispatch.Listen(act.Event(), callback, listener.UseCapture())
		} else {
			err = fmt.Errorf("couldnt find action class %v", cls)
			break
		}
	}
	if err != nil {
		return nil, err
	}
	return game, err
}

//
func (this *Game) PushParserSource(userSource G.SourceLookup) {
	this.parserSource.PushSource(func() (ret *GameObject) {
		// setup callback context:
		play := &GameEventAdapter{Game: this}
		// call the user function
		res := userSource(play)
		// unpack the result
		if par, ok := res.(ObjectAdapter); ok {
			ret = par.gobj
		}
		return ret
	})
}

//
func (this *Game) PopParserSource() {
	this.parserSource.PopSource()
}

//
// change the user's parent lookup (IObject -> name) into
// the runtime's parent lookup (GameObject->GameObject)
func (this *Game) PushParentLookup(userLookup G.TargetLookup) {
	this.parentLookup.PushLookup(func(gobj *GameObject) (ret *GameObject) {
		// setup callback context:
		play := &GameEventAdapter{Game: this}
		obj := ObjectAdapter{this, gobj}
		// call the user function
		res := userLookup(play, obj)
		// unpack the result
		if par, ok := res.(ObjectAdapter); ok {
			ret = par.gobj
		}
		return ret
	})
}

//
func (this *Game) PopParentLookup() {
	this.parentLookup.PopLookup()
}

//
// For testing:
//
func (this *Game) RunCommand(in string) (err error) {
	if _, res, e := this.parser.Parse(in); e != nil {
		err = e
	} else if e := res.Run(); e != nil {
		err = e
	}
	return err

}

//
//
// Run the event queue till there's an error
//
func (this *Game) ProcessEvents() (err error) {
	for err == nil && !this.queue.Empty() {
		tgt, msg := this.queue.Pop()
		// see also: Go()
		path := E.NewPathTo(tgt)
		this.log.Printf("sending `%s` to: %s", msg.Name, path)
		if runDefault, e := msg.Send(path); e != nil {
			this.log.Println("error", e)
			err = e
		} else {
			if runDefault {
				if act, ok := msg.Data.(*RuntimeAction); !ok {
					err = fmt.Errorf("unknown action data %T", msg.Data)
				} else {
					act.runDefaultActions()
				}
			}
		}
	}
	return err
}

// FIX: TEMP(ish)
// it might be better to add a name search (interface) to the model
// and then use the id in the runtime.
func (this *Game) FindObject(name string) (ret *GameObject, okay bool) {
	if info, err := this.Model.Instances.FindInstance(name); err == nil {
		ret = this.Objects[info.Id()]
		okay = true
	}
	return ret, okay
}

// FIX: TEMP(ish)
func (this *Game) FindFirstOf(cls *M.ClassInfo, _ ...bool) (ret *GameObject) {
	for _, o := range this.Objects {
		if o.info.Class() == cls {
			ret = o
			break
		}
	}
	return ret
}

//
// mainly for testing; manual send of an event
// FIX: merge this with runCommand()
func (this *Game) SendEvent(event string, nouns ...string,
) (err error,
) {
	if action, e := this.Model.Events.FindEventByName(event); e != nil {
		err = e
	} else {
		if act, e := this.newRuntimeAction(action, nouns...); e != nil {
			err = e
		} else {
			tgt := ObjectTarget{this, act.objs[0]}
			this.queue.QueueEvent(tgt, action.Event(), act)
		}
	}
	return err
}

//
// FIX: merge with runCommand()
//
func (this *Game) newRuntimeAction(action *M.ActionInfo, nouns ...string,
) (ret *RuntimeAction, err error,
) {
	types := action.NounSlice()
	switch diff := len(nouns) - len(types); {
	case diff < 0:
		err = fmt.Errorf("too few nouns specified for '%s', %d", action, diff)
	case diff > 0:
		err = fmt.Errorf("too many nouns specified for '%s', +%d", action, diff)
	default:
		objs := make([]*GameObject, len(types))
		keys := []string{"Source", "Target", "Context"}
		values := make(map[string]TemplateValues)

		for i, class := range types {
			noun, key := nouns[i], keys[i]
			inst, e := this.Model.Instances.FindInstanceWithClass(noun, class)
			if e != nil {
				err = e
				break
			}
			gobj := this.Objects[inst.Id()]
			values[key] = gobj.data
			objs[i] = gobj
		}
		if err == nil {
			ret = &RuntimeAction{this, action, objs, values, nil}
		}
	}
	return ret, err
}

//
// Called from the parser after it has succesfully found the command and nouns
//
func (this *Game) RunAction(action *M.ActionInfo, instances []string) (err error) {
	// make sure the source class matches
	sourceObj := this.parserSource.FindSource()
	if sourceObj == nil {
		err = fmt.Errorf("couldnt find command source for %s", action)
	} else {
		sourceClass := sourceObj.info.Class()
		if action.Source() != sourceClass && !sourceClass.HasParent(action.Source()) {
			err = fmt.Errorf("source class for %s doesnt match", action)
		} else {
			// inject the source object along with the other nouns
			types := action.NounSlice()
			instances = append([]string{sourceObj.Id().String()}, instances...)
			if len(types) != len(instances) {
				err = fmt.Errorf("mismatched nouns %d!=%d", len(types), len(instances))
			} else {
				keys := []string{"Source", "Target", "Context"}
				if len(instances) > len(keys) {
					err = fmt.Errorf("too many nouns %v", instances)
				} else {
					values := make(map[string]TemplateValues)
					objs := make([]*GameObject, len(types))

					for i, id := range instances {
						// convert to string id for net sake
						gobj, key := this.Objects[M.MakeStringId(id)], keys[i]
						if gobj == nil {
							return fmt.Errorf("unknown object %s", id)
						}
						values[key] = gobj.data
						objs[i] = gobj
					}

					tgt := ObjectTarget{this, objs[0]}
					act := &RuntimeAction{this, action, objs, values, nil}
					//log.Println("!!!", tgt, action)

					this.queue.QueueEvent(tgt, action.Event(), act)
				}
			}
		}
	}
	return err
}
