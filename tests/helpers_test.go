package tests

import (
	"bytes"
	"github.com/ionous/mars/rt"
	. "github.com/ionous/mars/script"
	"github.com/ionous/mars/script/backend"
	"github.com/ionous/mars/std"
	"github.com/ionous/sashimi/compiler"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/metal"
	"github.com/ionous/sashimi/parser"
	"github.com/ionous/sashimi/play"
	"github.com/ionous/sashimi/play/api"
	"github.com/ionous/sashimi/play/parse"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/lang"
	"github.com/ionous/sashimi/util/sbuf"
	"io/ioutil"
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
	s := string(p)
	if strings.TrimSpace(s) != "" {
		out.t.Log(s)
	}
	return out.buf.Write(p)
}

func (out TestWriter) Flush() string {
	ret := out.buf.String()
	out.buf.Reset()
	return ret
}

type ParentCreator func(run rt.Runtime) api.LookupParents

type StandardParents struct {
	run rt.Runtime
}

func (p StandardParents) LookupParent(i meta.Instance) (ret meta.Instance, err error) {
	if i, e := std.Parent(rt.Object{i}).GetObject(p.run); e != nil {
		err = e
	} else {
		ret = i.Instance
	}
	return
}

// NewStandardParents satisfies ParentCreator
func NewStandardParents(run rt.Runtime) api.LookupParents {
	return StandardParents{run}
}

func NewTestGameScript(t *testing.T, s backend.Declaration, gen string, pc ParentCreator) (ret TestGame, err error) {
	src := &S.Statements{}
	if gen == "no parser" {
		noParser := The("kind", Called("no parser"), Exists())
		if e := noParser.Generate(src); e != nil {
			err = e
			return
		}
	}
	if e := s.Generate(src); e != nil {
		err = e
	} else {
		ret, err = NewTestGameSource(t, *src, gen, pc)
	}
	return
}

func NewTestGameSource(t *testing.T, src S.Statements, gen string, pc ParentCreator) (ret TestGame, err error) {
	if model, e := compiler.Compile(src, ioutil.Discard); e != nil {
		err = e
	} else {
		storage := make(metal.ObjectValueMap)
		//saver := &TestSaver{}
		writer := TestWriter{t, bytes.NewBuffer(nil)}
		logger := TestLogger{t}
		values := TestValueMap{storage}
		modelApi := metal.NewMetal(model.Model, values)
		parents := &api.ParentHolder{}
		cfg := play.NewConfig().SetWriter(writer).SetLogger(logger).SetParentLookup(parents)
		//.SetSaveLoad(mem.NewSaveHelper("testing", storage, saver))
		//
		game := cfg.MakeGame(modelApi)
		if pc != nil {
			parents.Parents = pc(game)
		}
		//
		if parser, e := parse.NewObjectParser(modelApi, ident.MakeId(gen)); e != nil {
			err = e
		} else {
			ret = TestGame{t, game, modelApi, writer, parser, storage}
		}
	}
	return
}

// MARS - FIX the test game -- any game -- shouldnt require a parser.
// that should be on the front end, wrapping the game.
// ditto the "player"
// the understandings used by the parser can just sit there
// in the future, maybe we could put the understanding in an outer layer
func NewTestGame(t *testing.T, s backend.Declaration) (ret TestGame, err error) {
	src := &S.Statements{}
	player := The("actor", Called("player"), Exists())
	if e := player.Generate(src); e != nil {
		err = e
	} else if e := s.Generate(src); e != nil {
		err = e
	} else {
		ret, err = NewTestGameSource(t, *src, "player", nil)
	}
	return
}

type TestGame struct {
	t      *testing.T
	Game   play.Game
	Metal  *metal.Metal
	out    TestWriter
	Parser parser.P
	//saver  *TestSaver
	values metal.ObjectValueMap
}

func (test *TestGame) RunNamedAction(name string, params ...meta.Generic) (err error) {
	if id := ident.MakeId(lang.StripArticle(name)); id.Empty() {
		err = errutil.New("unnamed action", name)
	} else {
		err = test.Game.RunAction(id, test.Game, params...)
	}
	return
}

func (test *TestGame) Commence() (ret []string, err error) {
	if story, ok := meta.FindFirstOf(test.Metal, ident.MakeId("stories")); !ok {
		err = errutil.New("should have test story")
	} else if e := test.RunNamedAction("commence", story); e != nil {
		err = e
	} else {
		ret, err = test.FlushOutput()
	}
	return
}

func (test *TestGame) RunInput(s string) (ret []string, err error) {
	in := parser.NormalizeInput(s)
	if _, m, e := test.Parser.ParseInput(in); e != nil {
		err = errutil.New("failed parse", sbuf.Q(s), "=>", sbuf.Q(in), e)
	} else if act, objs, e := m.(*parse.ObjectMatcher).GetMatch(); e != nil {
		err = errutil.New("no match", sbuf.Q(s), "=>", sbuf.Q(in), e)
	} else {
		parms := make([]meta.Generic, len(objs))
		for i, o := range objs {
			parms[i] = rt.Object{o}
		}
		if e := test.Game.RunAction(act.GetId(), test.Game, parms...); e != nil {
			err = e
		} else if r, e := test.FlushOutput(); e != nil {
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
