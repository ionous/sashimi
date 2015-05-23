package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	"log"
)

type IOutput interface {
	Print(...interface{})
	Println(...interface{})
	Write(p []byte) (n int, err error)
}

// FIX: standarize member exports by splitting game into smaller classes and interfaces
type Game struct {
	Model          *M.Model
	Objects        GameObjects
	Dispatchers    ClassDispatchers
	console        IOutput
	parser         *ModelParser
	queue          *E.Queue
	nullobj        *NullObject
	defaultActions DefaultActions
	log            *log.Logger
	references     M.References
	parentLookup   ParentLookupStack
	parserSource   ParserSourceStack
}

// each action can have a chain of default actions
type DefaultActions map[*M.ActionInfo][]G.Callback

func NewGame(model *M.Model, output IOutput) (game *Game, err error) {
	log := log.New(output, "game: ", log.Lshortfile)
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
		log,
		M.NewReferences(model.Classes, model.Instances, model.Tables),
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
		log.Printf("creating listener %s", listener.String())

		if inst := listener.Instance(); inst != nil {
			obj := objects[inst.Id()]
			obj.dispatcher.Listen(act.Event(), callback, listener.UseCapture())
		} else if cls := listener.Class(); cls != nil {
			dispatch := dispatchers.CreateDispatcher(cls)
			dispatch.Listen(act.Event(), callback, listener.UseCapture())
		} else {
			err = fmt.Errorf("couldnt find action class %r", cls)
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
		err = fmt.Errorf("too few nouns specified for '%s', %d", action.Name(), diff)
	case diff > 0:
		err = fmt.Errorf("too many nouns specified for '%s', +%d", action.Name(), diff)
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
			obj := this.Objects[inst.Id()]
			values[key] = obj.values.data
			objs[i] = obj
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
func (this *Game) runCommand(action *M.ActionInfo, instances []string) (err error) {
	// make sure the source class matches
	sourceObj := this.parserSource.FindSource()
	//
	if sourceObj == nil {
		err = fmt.Errorf("couldnt find command source for %s", action.Name())
	} else {
		source := sourceObj.Info()
		sourceClass := source.Class()
		if action.Source() != sourceClass && !sourceClass.HasParent(action.Source()) {
			err = fmt.Errorf("source class for %s doesnt match", action.Name())
		} else {
			// inject the source object along with the other nouns
			types := action.NounSlice()
			instances = append([]string{source.Id().String()}, instances...)
			keys := []string{"Source", "Target", "Context"}
			values := make(map[string]TemplateValues)
			objs := make([]*GameObject, len(types))

			for i, id := range instances {
				obj, key := this.Objects[M.StringId(id)], keys[i]
				values[key] = obj.values.data
				objs[i] = obj
			}

			tgt := ObjectTarget{this, objs[0]}
			act := &RuntimeAction{this, action, objs, values, nil}
			this.queue.QueueEvent(tgt, action.Event(), act)
		}
	}
	return err
}
