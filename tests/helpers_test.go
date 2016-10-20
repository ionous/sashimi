package tests

import (
	"fmt"
	"github.com/ionous/mars/script"
	"github.com/ionous/sashimi/compiler"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/metal"
	"github.com/ionous/sashimi/net/mem"
	"github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/runtime/parse"
	"github.com/ionous/sashimi/util"
	"github.com/ionous/sashimi/util/ident"
	"log"
	"strings"
	"testing"
)

var _ = fmt.Print

type LogOutput struct {
	t *testing.T
}

func Log(t *testing.T) LogOutput {
	return LogOutput{t}
}

func (out LogOutput) Write(bytes []byte) (int, error) {
	out.t.Log(strings.TrimSpace(string(bytes)))
	return len(bytes), nil
}

//
type TestOutput struct {
	t *testing.T
	*util.BufferedOutput
}

//
// Standard output.
//
func (out TestOutput) ScriptSays(lines []string) {
	for _, l := range lines {
		out.Println(l)
	}
}

func (out TestOutput) ActorSays(whose meta.Instance, lines []string) {
	var name string
	if prop, ok := whose.FindProperty("name"); ok {
		name = prop.GetValue().GetText()
	}

	for _, l := range lines {
		out.Println(name, ": ", l)
	}
}

func (out TestOutput) Log(s string) {
	out.t.Log(strings.TrimSpace(s))
}

type ParentCreator func(meta.Model) api.LookupParents

func NewTestGameSource(t *testing.T, s *script.Script, src string, pc ParentCreator) (ret TestGame, err error) {
	if statements, e := s.BuildStatements(); e != nil {
		err = e
	} else {
		if model, e := compiler.Compile(Log(t), statements); e != nil {
			err = e
		} else {
			storage := make(metal.ObjectValueMap)
			saver := &TestSaver{}
			cons := TestOutput{t, &util.BufferedOutput{}}
			values := TestValueMap{storage}
			modelApi := metal.NewMetal(model.Model, values)
			var parents api.LookupParents
			if pc != nil {
				parents = pc(modelApi)
			}
			cfg := R.NewConfig().SetCalls(model.Calls).SetOutput(cons).SetSaveLoad(mem.NewSaveHelper("testing", storage, saver)).SetParentLookup(parents)
			//
			game := cfg.MakeGame(modelApi)
			if parser, e := parse.NewObjectParser(game, ident.MakeId(src)); e != nil {
				err = e
			} else {
				ret = TestGame{t, game, model, cons, parser, saver, storage}
			}
		}
	}
	return
}

//
func NewTestGame(t *testing.T, s *script.Script) (ret TestGame, err error) {
	return NewTestGameSource(t, s, "player", nil)
}

type TestGame struct {
	t    *testing.T
	Game R.Game
	compiler.MemoryResult
	out    TestOutput
	Parser parser.P
	saver  *TestSaver
	values metal.ObjectValueMap
}

func (test *TestGame) Commence() (ret []string, err error) {
	if story, ok := meta.FindFirstOf(test.Game, ident.MakeId("stories")); !ok {
		err = fmt.Errorf("should have test story")
	} else if _, e := test.Game.QueueAction("commence", story.GetId()); e != nil {
		err = e
	} else {
		ret, err = test.FlushOutput()
	}
	return
}

func (test *TestGame) RunInput(s string) (ret []string, err error) {
	if e := test.Game.ProcessActions(); e != nil {
		err = e
	} else {
		in := parser.NormalizeInput(s)
		if p, m, e := test.Parser.ParseInput(in); e != nil {
			test.out.Log(fmt.Sprintf("RunInput: failed parse: %v orig: '%s' in: '%s' e: '%s'", p, s, in, e))
			err = e
		} else if act, objs, e := m.(*parse.ObjectMatcher).GetMatch(); e != nil {
			test.out.Log(fmt.Sprint("RunInput: no match: ", s, e))
			err = e
		} else {
			test.Game.QueueActionInstances(act, objs)
			// the standard rules send an "ending the turn", we do not have to.
			if r, e := test.FlushOutput(); e != nil {
				err = e
			} else {
				ret = r
			}
		}
	}
	return
}

func (test *TestGame) FlushOutput() (ret []string, err error) {
	if e := test.Game.ProcessActions(); e != nil {
		test.out.Println(e)
		err = e
	} else {
		ret = test.out.Flush()
	}
	return
}

// TestSaver implements mem.MemSaver
type TestSaver struct {
	blob mem.SaveGameBlob
}

func (t *TestSaver) SaveBlob(slot string, blob mem.SaveGameBlob) (string, error) {
	t.blob = blob
	return slot, nil
}
func (t *TestSaver) LoadBlob(slot string) (mem.SaveGameBlob, error) {
	return t.blob, nil
}

// TestValueMap implements metal.ObjectValue
type TestValueMap struct {
	values metal.ObjectValueMap
}

// GetValue succeeds if SetValue was called on a corresponding obj.field.
func (m TestValueMap) GetValue(obj, field ident.Id) (ret interface{}, okay bool) {
	return m.values.GetValue(obj, field)
}

// SetValue always succeeds, storing the passed value to the map at obj.field.
func (m TestValueMap) SetValue(obj, field ident.Id, value interface{}) (err error) {
	log.Println("SetValue:", obj, field, value)
	return m.values.SetValue(obj, field, value)
}
