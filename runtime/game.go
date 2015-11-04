package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/model/table"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"log"
	"math/rand"
	"reflect"
)

// FIX: standarize member exports by splitting game into smaller classes and interfaces; focus on injecting what game needs, and allowing providers to decorate/instrument what they need.
type Game struct {
	Model          *M.Model
	Objects        GameObjects
	Dispatchers    Dispatchers
	output         IOutput
	queue          EventQueue
	frame          EventFrame
	defaultActions DefaultActions
	SystemActions  SystemActions
	log            *log.Logger
	Properties     PropertyWatchers
	parentLookup   ParentLookupStack
	Globals        Globals
	rand           *rand.Rand
	Tables         table.Tables
}

type logAdapter struct {
	output IOutput
}

func (log logAdapter) Write(p []byte) (n int, err error) {
	log.output.Log(string(p))
	return len(p), nil
}

// each action can have a chain of default actions
type CallbackPair struct {
	src  ident.Id
	call G.Callback
}

// FIX: change callback structure to contain info on location
func (f CallbackPair) String() string {
	return fmt.Sprint(f.call)
}

type DefaultActions map[*M.ActionInfo][]CallbackPair

type Globals map[ident.Id]reflect.Value

func (cfg Config) NewGame(model *M.Model) (_ *Game, err error) {
	log := log.New(logAdapter{cfg.Output}, "game: ", log.Lshortfile)
	dispatchers := NewDispatchers(log)
	tables := model.Tables.Clone()
	objects, e := CreateGameObjects(model.Instances, tables)
	if e != nil {
		return nil, e
	}

	globals := make(Globals)
	// DISABLED:
	// for k, gen := range model.Generators {
	// 	globals[k] = reflect.New(gen)
	// }
	frame := cfg.Frame
	if frame == nil {
		frame = DefaultEventFrame
	}

	game := &Game{
		model,
		objects,
		dispatchers,
		cfg.Output,
		EventQueue{E.NewQueue()},
		frame,
		make(DefaultActions),
		NewSystemActions(),
		log,
		PropertyWatchers{},
		ParentLookupStack{},
		globals,
		rand.New(rand.NewSource(1)),
		tables,
	}
	for _, handler := range model.ActionHandlers {
		src := handler.Callback
		if cb := cfg.Calls.Lookup(src); cb == nil {
			err = errutil.Append(err, fmt.Errorf("couldnt find callback for action", handler))
		} else {
			cb := CallbackPair{src, cb}
			act := handler.Action
			arr := game.defaultActions[act]
			// FIX: for now treating target as bubble,
			// really the compiler should hand off a sorted flat list based on three separate groups
			// target growing in the same direction as after, but distinctly in the middle of things.
			if !handler.UseCapture() {
				arr = append(arr, cb)
			} else {
				// prepend:
				arr = append([]CallbackPair{cb}, arr...)
			}
			game.defaultActions[act] = arr
		}
	}

	if err != nil {
		return nil, err
	}

	// FUTURE: move into scenes, with a handle for removal
	for _, listener := range model.EventListeners {
		src := listener.Callback
		if cb := cfg.Calls.Lookup(src); cb == nil {
			err = errutil.Append(err, fmt.Errorf("couldnt find callback for listener", listener))
		} else {
			act := listener.Action
			callback := GameCallback{game, listener, CallbackPair{src, cb}}

			var id ident.Id
			if inst := listener.Instance; inst != nil {
				id = inst.Id
			} else if cls := listener.Class; cls != nil {
				id = cls.Id
			} else {
				e := fmt.Errorf("couldnt create listener %v", listener)
				err = errutil.Append(err, e)
				continue
			}
			dispatch := dispatchers.CreateDispatcher(id)
			dispatch.Listen(act.EventName, callback, listener.UseCapture())
		}
	}
	if err != nil {
		return nil, err
	}
	return game, err
}

func (g *Game) Random(n int) int {
	return g.rand.Intn(n)
}

// PushParentLookup function into the game's determination of which object is which object's container.
// Changes the user's parent lookup (IObject -> name) into
// the runtime's parent lookup (GameObject->GameObject).
// FIX: inject an interface, via the constructor, which makes this possible
// possibly inject the game/object adapter factory even.
// then the caller can have the handle which does the push
// and game can remain ignorant of the push (or not) process.
func (g *Game) PushParentLookup(userLookup G.TargetLookup) {
	g.parentLookup.PushLookup(func(gobj *GameObject) (ret *GameObject) {
		// setup callback context:
		play := NewGameAdapter(g)
		obj := NewObjectAdapter(g, gobj)
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
func (g *Game) PopParentLookup() {
	g.parentLookup.PopLookup()
}

// ProcessEvents in the event queue till empty, or errored.
func (g *Game) ProcessEvents() (err error) {
	for !g.queue.Empty() {
		tgt, msg := g.queue.Pop()
		if e := g.frame.SendMessage(tgt, msg); e != nil {
			g.log.Println("error", e)
			err = e
			break
		}
	}
	return err
}

// FIX: TEMP(ish)
// it might be better to add a name search (interface) to the model
// and then use the id in the runtime.
func (g *Game) FindObject(name string) (ret *GameObject, okay bool) {
	if info, ok := g.Model.Instances.FindInstance(name); ok {
		ret = g.Objects[info.Id]
		okay = true
	}
	return ret, okay
}

// FIX: TEMP(ish)
func (g *Game) FindFirstOf(cls *M.ClassInfo, _ ...bool) (ret *GameObject) {
	for _, o := range g.Objects {
		if o.Class() == cls {
			ret = o
			break
		}
	}
	return ret
}

//
// mainly for testing; manual send of an event
// FIX: merge game with runCommand()
func (g *Game) SendEvent(event string, nouns ...ident.Id,
) (err error,
) {
	if action, e := g.Model.Events.FindEventByName(event); e != nil {
		err = e
	} else {
		if act, e := g.newRuntimeAction(action, nouns...); e != nil {
			err = e
		} else {
			tgt := ObjectTarget{g, act.objs[0]}
			g.queue.QueueEvent(tgt, action.EventName, act)
		}
	}
	return err
}

//
// FIX: merge with runCommand()
//
func (g *Game) newRuntimeAction(action *M.ActionInfo, nouns ...ident.Id,
) (ret *RuntimeAction, err error,
) {
	types := action.NounTypes
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
			if gobj, ok := g.Objects[noun]; !ok {
				err = M.InstanceNotFound(noun.String())
				break
			} else if !gobj.Class().CompatibleWith(class.Id) {
				err = TypeMismatch(noun.String(), class.String())
				break
			} else {
				values[key] = gobj.vals
				objs[i] = gobj
			}
		}
		if err == nil {
			ret = &RuntimeAction{g, action, objs, values, nil}
		}
	}
	return ret, err
}
