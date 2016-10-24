package tests

import (
	"bytes"
	"github.com/ionous/mars/script"
	"github.com/ionous/sashimi/compiler"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/metal"
	"github.com/ionous/sashimi/parser"
	"github.com/ionous/sashimi/play"
	"github.com/ionous/sashimi/play/api"
	"github.com/ionous/sashimi/play/parse"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/sbuf"
	"strings"
	"testing"
)

type TestLogger struct {
	t *testing.T
}

func (out TestLogger) Write(p []byte) (int, error) {
	out.t.Log(string(p))
	return len(p), nil
}

//
type TestWriter struct {
	t   *testing.T
	buf *bytes.Buffer
}

// Standard output.
func (out TestWriter) Write(p []byte) (n int, err error) {
	out.t.Log(string(p))
	return out.buf.Write(p)
}

func (out TestWriter) Flush() string {
	ret := out.buf.String()
	out.buf.Reset()
	return ret
}

type ParentCreator func(meta.Model) api.LookupParents

func NewTestGameSource(t *testing.T, s *script.Script, src string, pc ParentCreator) (ret TestGame, err error) {
	if statements, e := s.BuildStatements(); e != nil {
		err = e
	} else {
		logger := TestLogger{t}
		if model, e := compiler.Compile(logger, statements); e != nil {
			err = e
		} else {
			storage := make(metal.ObjectValueMap)
			//saver := &TestSaver{}
			writer := TestWriter{t, bytes.NewBuffer(nil)}
			values := TestValueMap{storage}
			modelApi := metal.NewMetal(model.Model, values)
			var parents api.LookupParents
			if pc != nil {
				parents = pc(modelApi)
			}
			cfg := play.NewConfig().SetWriter(writer).SetLogger(logger).SetParentLookup(parents)
			//.SetSaveLoad(mem.NewSaveHelper("testing", storage, saver))
			//
			game := cfg.MakeGame(modelApi)
			if parser, e := parse.NewObjectParser(modelApi, ident.MakeId(src)); e != nil {
				err = e
			} else {
				ret = TestGame{t, game, modelApi, writer, parser, storage}
			}
		}
	}
	return
}

// MARS - FIX the test game -- any game -- shouldnt require a parser.
// that should be on the front end, wrapping the game.
// ditto the "player"
// the understandings used by the parser can just sit there
// in the future, maybe we could put the understanding in an outer layer
func NewTestGame(t *testing.T, s *script.Script) (ret TestGame, err error) {
	return NewTestGameSource(t, s, "player", nil)
}

type TestGame struct {
	t      *testing.T
	Game   play.Game
	Model  meta.Model
	out    TestWriter
	Parser parser.P
	//saver  *TestSaver
	values metal.ObjectValueMap
}

func (test *TestGame) Commence() (ret []string, err error) {
	if story, ok := meta.FindFirstOf(test.Model, ident.MakeId("stories")); !ok {
		err = errutil.New("should have test story")
	} else if e := test.Game.RunAction("commence", story); e != nil {
		err = e
	} else {
		ret, err = test.FlushOutput()
	}
	return
}

func (test *TestGame) RunInput(s string) (ret []string, err error) {
	in := parser.NormalizeInput(s)
	if p, m, e := test.Parser.ParseInput(in); e != nil {
		test.out.buf.WriteString(sbuf.New("RunInput: failed parse:", sbuf.Value{p}, "orig:", s, "in:", in, "e:", e).Line())
		err = e
	} else if act, objs, e := m.(*parse.ObjectMatcher).GetMatch(); e != nil {
		test.out.buf.WriteString(sbuf.New("RunInput: no match:", s, e).Line())
		err = e
	} else {
		test.Game.RunAction(act.GetId(), objs)
		// the standard rules send an "ending the turn", we do not have to.
		if r, e := test.FlushOutput(); e != nil {
			err = e
		} else {
			ret = r
		}
	}
	return
}

func (test *TestGame) FlushOutput() ([]string, error) {
	s := test.out.Flush()
	return strings.Split(s, "\n"), nil
}

// FIX: diabled for mars testing
// // TestSaver implements mem.MemSaver
// type TestSaver struct {
// 	blob mem.SaveGameBlob
// }

// func (t *TestSaver) SaveBlob(slot string, blob mem.SaveGameBlob) (string, error) {
// 	t.blob = blob
// 	return slot, nil
// }
// func (t *TestSaver) LoadBlob(slot string) (mem.SaveGameBlob, error) {
// 	return t.blob, nil
// }

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
	return m.values.SetValue(obj, field, value)
}
