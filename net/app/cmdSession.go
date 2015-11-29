package app

import (
	"encoding/json"
	"fmt"
	"github.com/ionous/sashimi/change"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	// FIX? can we carve out metal from the cmds?
	"github.com/ionous/sashimi/metal"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/net/session"
	"github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/runtime/parse"
	"github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/util/ident"
	"io"
	"log"
	"sync"
)

var playerId = ident.MakeId("player")

// NewCommandSession create a web session which uses jsonapi style commands.
func NewCommandSession(id string, model *M.Model, calls api.LookupCallbacks) (ret ICommandSession, err error) {
	metal := metal.NewMetal(model, make(metal.ObjectValueMap))
	out := NewCommandOutput(id, NewObjectSerializer(make(KnownObjectMap)))
	if s, e := NewPartialSession(out, metal, calls); e != nil {
		err = e
	} else {
		ret = &CommandSession{s, &sync.RWMutex{}}
	}
	return ret, err
}

type ICommandSession interface {
	resource.IResource

	RLock()
	RUnlock()
	Lock()
	Unlock()

	FrameCount() int
}

type PartialSession struct {
	game *standard.StandardCore
	out  *CommandOutput
}

func NewPartialSession(out *CommandOutput, m meta.Model, calls api.LookupCallbacks) (ret *PartialSession, err error) {

	cfg := R.NewConfig().SetCalls(calls).SetOutput(out).SetFrame(out).SetParentLookup(standard.ParentLookup{})
	watched := change.NewModelWatcher(out, m)
	if game, e := cfg.NewGame(watched); e != nil {
	} else {
		// after creating the game, but before running it --
		if game, e := standard.NewStandardCore(game); e != nil {
			err = e
		} else {
			ret = &PartialSession{game, out}
		}
	}
	return ret, err
}

// CommandSession holds game data associated with a particular client.
// the parent reource locks and unlocks it.
// the parent resource needs an interface we can provide.
// we have
type CommandSession struct {
	*PartialSession
	*sync.RWMutex
	// type RWMutex provides:
	//    func (rw *RWMutex) RLock()
	//    func (rw *RWMutex) RUnlock()
	//    func (rw *RWMutex) Lock()
	//    func (rw *RWMutex) Unlock()
}

func (s *CommandSession) FrameCount() int {
	return s.game.Frame()
}

func (s *CommandSession) Post(reader io.Reader) (ret resource.Document, err error) {
	if r, e := s.PartialSession.Post(reader); e != nil {
		err = e
	} else {
		ret = r
	}
	return
}

// Find the named game resource.
func (s *PartialSession) Find(name string) (ret resource.IResource, okay bool) {
	mdl := s.game
	switch name {
	// by default, objects are grouped by their class:
	default:
		if cls, ok := mdl.GetClass(ident.MakeId(name)); ok {
			ret, okay = ObjectResource(mdl, cls.GetId(), s.out.serial), true
		}
	// a request for information about a class:
	case "class":
		ret, okay = ClassResource(mdl), true
		// a request for information about a parser input action:
	case "action":
		ret, okay = ParserResource(mdl), true
	}
	return ret, okay
}

// Query on the session itself always returns a blank document. See: post.
func (s *PartialSession) Query() (ret resource.Document) {
	return
}

// Post input to the game.
func (s *PartialSession) Post(reader io.Reader) (ret resource.Document, err error) {
	// hrmm.... by sending the io.Reader, we can "decode" different types of json...
	// but, we have to do that wherever we need... [ add a Blank struct maker ? ]
	decoder, input := json.NewDecoder(reader), CommandInput{}
	if e := decoder.Decode(&input); e != nil {
		err = e
	} else {
		if e := s._handleInput(input); e != nil {
			err = e
		} else {
			doc := resource.NewDocumentBuilder(&ret)
			s.out.FlushDocument(doc)
		}
	}
	return ret, err
}

// Send the passed input to the game.
func (s *PartialSession) _handleInput(input CommandInput) (err error) {
	// NOTE: cmd session doesnt support quit
	if s.game.IsComplete() {
		err = session.SessionClosed{"game finished"}
	} else if in := input.Input; in != "" {
		log.Println("Input", in)
		s.game.HandleInput(parser.NormalizeInput(in))
	} else {
		// Run json'd clicky action:
		id := ident.MakeId(input.Action)
		if act, ok := s.game.GetAction(id); !ok {
			err = fmt.Errorf("unknown action %s", input.Action)
			//FIX? RunActions injects the player, that works out well, but is a little strange.
		} else if om, e := parse.NewObjectMatcher(act, playerId, s.game.Model); e != nil {
			err = e
		} else {
			for _, n := range input.Nouns() {
				id := ident.MakeId(n)
				if e := om.AddObject(id); e != nil {
					err = e
					break
				}
			}
			if err == nil {
				if act, objs, e := om.GetMatch(); e != nil {
					err = e
				} else {
					s.game.QueueActionInstances(act, objs)
					s.game.EndTurn("end turn")
				}
			}
		}
	}
	return err
}
