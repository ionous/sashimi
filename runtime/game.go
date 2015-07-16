package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"log"
	"reflect"
)

// FIX: standarize member exports by splitting game into smaller classes and interfaces
type Game struct {
	Model          *M.Model
	Objects        GameObjects
	Dispatchers    Dispatchers
	output         IOutput
	queue          *E.Queue
	defaultActions DefaultActions
	SystemActions  SystemActions
	log            *log.Logger
	Properties     PropertyWatchers
	parentLookup   ParentLookupStack
	Globals        Globals
}

type logAdapter struct {
	output IOutput
}

func (log logAdapter) Write(p []byte) (n int, err error) {
	log.output.Log(string(p))
	return len(p), nil
}

// each action can have a chain of default actions
type DefaultActions map[*M.ActionInfo][]G.Callback

type Globals map[ident.Id]reflect.Value

func NewGame(model *M.Model, output IOutput) (game *Game, err error) {
	log := log.New(logAdapter{output}, "game: ", log.Lshortfile)
	dispatchers := NewDispatchers(log)
	objects, e := CreateGameObjects(model.Instances, model.Tables)
	if e != nil {
		return nil, e
	}

	globals := make(Globals)
	for k, gen := range model.Generators {
		globals[k] = reflect.New(gen)
	}
	if err != nil {
		return nil, err
	}

	game = &Game{
		model,
		objects,
		dispatchers,
		output,
		E.NewQueue(),
		make(DefaultActions),
		NewSystemActions(),
		log,
		PropertyWatchers{},
		ParentLookupStack{},
		globals,
	}
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

		var id ident.Id
		if inst := listener.Instance(); inst != nil {
			id = inst.Id()
		} else if cls := listener.Class(); cls != nil {
			id = cls.Id()
		} else {
			e := fmt.Errorf("couldnt create listener %v", listener)
			err = errutil.Append(err, e)
			continue
		}
		dispatch := dispatchers.CreateDispatcher(id)
		dispatch.Listen(act.Event(), callback, listener.UseCapture())
	}
	if err != nil {
		return nil, err
	}
	return game, err
}

// PushParentLookup function into the game's determination of which object is which object's container.
// Changes the user's parent lookup (IObject -> name) into
// the runtime's parent lookup (GameObject->GameObject).
// FIX? move the the adapter??? dont think we should be referencing adapter in game...
func (game *Game) PushParentLookup(userLookup G.TargetLookup) {
	game.parentLookup.PushLookup(func(gobj *GameObject) (ret *GameObject) {
		// setup callback context:
		play := NewGameAdapter(game)
		obj := NewObjectAdapter(game, gobj)
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
func (game *Game) PopParentLookup() {
	game.parentLookup.PopLookup()
}

//
//
// Run the event queue till there's an error
//
func (game *Game) ProcessEvents() (err error) {
	for err == nil && !game.queue.Empty() {
		tgt, msg := game.queue.Pop()
		// see also: Go()
		path := E.NewPathTo(tgt)
		// game.log.Printf("sending `%s` to: %s", msg.Name, path)
		if runDefault, e := msg.Send(path); e != nil {
			game.log.Println("error", e)
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
func (game *Game) FindObject(name string) (ret *GameObject, okay bool) {
	if info, ok := game.Model.Instances.FindInstance(name); ok {
		ret = game.Objects[info.Id()]
		okay = true
	}
	return ret, okay
}

// FIX: TEMP(ish)
func (game *Game) FindFirstOf(cls *M.ClassInfo, _ ...bool) (ret *GameObject) {
	for _, o := range game.Objects {
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
func (game *Game) SendEvent(event string, nouns ...ident.Id,
) (err error,
) {
	if action, e := game.Model.Events.FindEventByName(event); e != nil {
		err = e
	} else {
		if act, e := game.newRuntimeAction(action, nouns...); e != nil {
			err = e
		} else {
			tgt := ObjectTarget{game, act.objs[0]}
			game.queue.QueueEvent(tgt, action.Event(), act)
		}
	}
	return err
}

//
// FIX: merge with runCommand()
//
func (game *Game) newRuntimeAction(action *M.ActionInfo, nouns ...ident.Id,
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
			if gobj, ok := game.Objects[noun]; !ok {
				err = M.InstanceNotFound(noun.String())
				break
			} else if !gobj.Class().CompatibleWith(class.Id()) {
				err = TypeMismatch(noun.String(), class.String())
				break
			} else {
				values[key] = gobj.vals
				objs[i] = gobj
			}
		}
		if err == nil {
			ret = &RuntimeAction{game, action, objs, values, nil}
		}
	}
	return ret, err
}
