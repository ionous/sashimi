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
	Parser         *ModelParser
	queue          *E.Queue
	nullobj        *NullObject
	defaultActions DefaultActions
	SystemActions  SystemActions
	log            *log.Logger
	Properties     PropertyWatchers
	parentLookup   ParentLookupStack
	parserSource   ParserSourceStack
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
		PropertyWatchers{},
		ParentLookupStack{},
		ParserSourceStack{},
	}

	parser, e := NewParser(game)
	if e != nil {
		return nil, e
	}
	game.Parser = parser

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
func (game *Game) PushParserSource(userSource G.SourceLookup) {
	game.parserSource.PushSource(func() (ret *GameObject) {
		// setup callback context:
		play := &GameEventAdapter{Game: game}
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
func (game *Game) PopParserSource() {
	game.parserSource.PopSource()
}

//
// change the user's parent lookup (IObject -> name) into
// the runtime's parent lookup (GameObject->GameObject)
func (game *Game) PushParentLookup(userLookup G.TargetLookup) {
	game.parentLookup.PushLookup(func(gobj *GameObject) (ret *GameObject) {
		// setup callback context:
		play := &GameEventAdapter{Game: game}
		obj := ObjectAdapter{game, gobj}
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
		if o.inst.Class() == cls {
			ret = o
			break
		}
	}
	return ret
}

//
// mainly for testing; manual send of an event
// FIX: merge game with runCommand()
func (game *Game) SendEvent(event string, nouns ...string,
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
func (game *Game) newRuntimeAction(action *M.ActionInfo, nouns ...string,
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
			inst, e := game.Model.Instances.FindInstanceWithClass(noun, class)
			if e != nil {
				err = e
				break
			}
			gobj := game.Objects[inst.Id()]
			values[key] = gobj.data
			objs[i] = gobj
		}
		if err == nil {
			ret = &RuntimeAction{game, action, objs, values, nil}
		}
	}
	return ret, err
}
