package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/model/table"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/runtime/memory"
	"github.com/ionous/sashimi/util/ident"
	"log"
	"math/rand"
	"reflect"
)

// FIX: remove
func (oa ObjectAdapter) Remove() {
	panic("not implemented")
}

func (null _Null) Remove() {
	panic("not implemented")
}
func (ga *GameEventAdapter) NewFrom(class string) G.IObject {
	panic("not supported")
}

// FIX: standarize member exports by splitting game into smaller classes and interfaces; focus on injecting what game needs, and allowing providers to decorate/instrument what they need.
type Game struct {
	ModelApi     api.Model
	output       IOutput
	queue        EventQueue
	frame        EventFrame
	calls        Callbacks
	log          *log.Logger
	parentLookup ParentLookupStack
	Globals      Globals
	rand         *rand.Rand
	Tables       table.Tables
}

type logAdapter struct {
	output IOutput
}

func (log logAdapter) Write(p []byte) (n int, err error) {
	log.output.Log(string(p))
	return len(p), nil
}

type Globals map[ident.Id]reflect.Value

func (g *Game) NewPlay(data interface{}, hint ident.Id) G.Play {
	adapter := NewGameAdapter(g)
	adapter.data = data.(*RuntimeAction)
	adapter.hint = hint
	return adapter
}

func (cfg RuntimeConfig) NewGame(model *M.Model) (_ *Game, err error) {
	log := log.New(logAdapter{cfg.Output}, "game: ", log.Lshortfile)

	tables := model.Tables.Clone()
	modelApi := memory.NewMemoryModel(model, make(memory.ObjectValueMap), tables)

	globals := make(Globals)
	// DISABLED:
	// for k, gen := range model.Generators {
	// 	globals[k] = reflect.New(gen)
	// }
	frame := cfg.Frame
	if frame == nil {
		frame = DefaultEventFrame
	}

	return &Game{
		modelApi,
		cfg.Output,
		EventQueue{E.NewQueue()},
		frame,
		cfg.Calls,
		log,
		ParentLookupStack{},
		globals,
		rand.New(rand.NewSource(1)),
		tables,
	}, nil
}

// class or instance id
func (g *Game) DispatchEvent(evt E.IEvent, target ident.Id) (err error) {
	if src, ok := g.ModelApi.GetEvent(evt.Id()); ok {
		if ls, ok := src.GetListeners(true); ok {
			err = E.Capture(evt, NewGameListeners(g, evt, target, ls))
		}
		if err == nil {
			if ls, ok := src.GetListeners(false); ok {
				err = E.Bubble(evt, NewGameListeners(g, evt, target, ls))
			}
		}
	}
	return
}

func (g *Game) Random(n int) int {
	return g.rand.Intn(n)
}

// PushParentLookup function into the game's determination of which object is which object's container.
// Changes the user's parent lookup (IObject -> name) into
// the runtime's parent lookup (Instance -> Instance).
// FIX: inject an interface, via the constructor, which makes this possible
// possibly inject the game/object adapter factory even.
// then the caller can have the handle which does the push
// and game can remain ignorant of the push (or not) process.
func (g *Game) PushParentLookup(userLookup G.TargetLookup) {
	g.parentLookup.PushLookup(func(gobj api.Instance) (ret api.Instance) {
		// setup callback context:
		play, obj := NewGameAdapter(g), NewObjectAdapter(g, gobj)
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

func (g *Game) QueueAction(act api.Action, objects []api.Instance) {
	tgt := ObjectTarget{g, objects[0]}
	data := &RuntimeAction{game: g, action: act, objs: objects}
	g.queue.QueueEvent(tgt, act.GetEvent().GetId(), data)
}

// mainly for testing; manual send of an event
func (g *Game) QueueEvent(event string, nouns ...ident.Id,
) (err error,
) {
	eventId := MakeStringId(event)
	if event, ok := g.ModelApi.GetEvent(eventId); !ok {
		err = fmt.Errorf("couldnt find event %s", event)
	} else if act, e := g.newRuntimeAction(event.GetAction(), nouns...); e != nil {
		err = e
	} else {
		tgt := ObjectTarget{g, act.objs[0]}
		g.queue.QueueEvent(tgt, event.GetId(), act)
	}
	return err
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

// TODO: unwind this.
func (g *Game) newRuntimeAction(action api.Action, nouns ...ident.Id,
) (ret *RuntimeAction, err error,
) {
	types := action.GetNouns()
	switch diff := len(nouns) - len(types); {
	case diff < 0:
		err = fmt.Errorf("too few nouns specified for '%s', %d", action, diff)
	case diff > 0:
		err = fmt.Errorf("too many nouns specified for '%s', +%d", action, diff)
	default:
		objs := make([]api.Instance, len(types))
		for i, class := range types {
			noun := nouns[i]
			if gobj, ok := g.ModelApi.GetInstance(noun); !ok {
				err = M.InstanceNotFound(noun.String())
				break
			} else if !g.ModelApi.AreCompatible(gobj.GetParentClass().GetId(), class) {
				err = TypeMismatch(noun.String(), class.String())
				break
			} else {
				objs[i] = gobj
			}
		}
		if err == nil {
			ret = &RuntimeAction{game: g, action: action, objs: objs}
		}
	}
	return ret, err
}
