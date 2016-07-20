package app

import (
	"encoding/json"
	"fmt"
	"github.com/ionous/sashimi/change"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/runtime/parse"
	"github.com/ionous/sashimi/standard/framework"
	"github.com/ionous/sashimi/util/ident"
	"io"
	"log"
)

// implements IResource for the session and Frame for the
type PartialSession struct {
	game *framework.StandardCore
	out  *CommandOutput
}

func NewPartialSession(m meta.Model, calls api.LookupCallbacks, saver api.SaveLoad, out *CommandOutput) (ret *PartialSession, err error) {
	watched := change.NewModelWatcher(out, m)
	cfg := R.NewConfig().
		SetCalls(calls).
		SetOutput(out).
		SetFrame(out).
		SetSaveLoad(saver).
		SetParentLookup(framework.NewParentLookup(watched))
	game := cfg.MakeGame(watched)
	// after creating the game, but before running it --
	if game, e := framework.NewStandardCore(game); e != nil {
		err = e
	} else {
		ret = &PartialSession{game, out}
	}
	return ret, err
}

func (s *PartialSession) FlushDocument(doc resource.DocumentBuilder) error {
	s.out.FlushDocument(doc)
	return nil
}

func (s *PartialSession) Frame() int {
	return s.game.Frame()
}

// Find the named game resource.
func (s *PartialSession) Find(name string) (ret resource.IResource, okay bool) {
	mdl := s.game
	switch name {
	// by default, objects are grouped by their class:
	default:
		if cls, ok := mdl.GetClass(ident.MakeId(name)); ok {
			ret, okay = ObjectResource(mdl, cls.GetId()), true
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

var playerId = ident.MakeId("player")

// Send the passed input to the game.
func (s *PartialSession) _handleInput(input CommandInput) (err error) {
	log.Println(fmt.Sprintf("processing %+v", input))
	// NOTE: cmd session doesnt support quit
	if s.game.IsComplete() {
		err = fmt.Errorf("game finished.")
	} else if in := input.Input; in != "" {
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
