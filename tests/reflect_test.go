package tests

import (
	"github.com/ionous/sashimi/compiler/call"
	"github.com/ionous/sashimi/compiler/extract"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path"
	"runtime"
	"strings"
	"testing"
)

func TestReflect(t *testing.T) {
	cb := func(G.Play) {
		t.Fatal("this is never called")
	}
	//cfg := call.Config{originalDir()}
	m := call.MakeMarker(cb)
	require.True(t, strings.HasSuffix(m.File, "reflect_test.go"), m.File)
}

func TestExtract(t *testing.T) {
	s := InitScripts()
	ReflectScript(s)

	// 1. parse an ast of a bunch of statements, gank out the callbacks
	//
	// 2. compile a script, the compilation process should replace callback functions with a lookup
	// [ compiler exposes the interface, others implement it ]
	// marker:= Callbacks.Compile(func)
	//
	// compiled?
	// UnifiedCallbacks.Compiler() -> compiler.Callbacks.Compile( g G.Func ) interface -> maybe you can use reflect to get file and line.
	// UnifiedCallbacks.Lookup() -> game.Callbacks.Lookup
	// ExtractedCallbacks.Compiler() -> uses compiler.extract.Extract
	// ExtractedCallbacks.GenerateSource(outputfile)
	// into that generate source package you could dump the JSON, GOB, or whatever of the Model so it can be reconstructed.
	// then you go compile that standalone package -- or, go import it into an application.
	//
	// an interface maybe? given a function, we store a name
	// to reconstitute we will register the callbacks via name
	// and so we will have a runtime interface to fetch them back.
	// the unified compilation and runtime is given an object which just maps
	// the separated does things like generate the callbacks based on the file given to the compiler
	// so step one, is given ... a file -- an inpt
	pcs := make([]uintptr, 1)
	if cnt := runtime.Callers(1, pcs); assert.True(t, cnt > 0, "couldnt get current counter") {
		pc := pcs[0]
		if fp := runtime.FuncForPC(pc - 1); assert.NotNil(t, fp) {
			filename, _ := fp.FileLine(pc - 1)

			if bytes, e := ioutil.ReadFile(filename); assert.NoError(t, e, "couldnt load file") {
				file := path.Base(filename)
				e := extract.Extract(file, bytes, func(f string, l int, sub []byte) {
					t.Log(f, l, string(sub))
				})
				require.NoError(t, e, "extract failed")
			}
		}
	}
}
func ReflectScript(s *Script) {
	s.The("stories", When("commencing").Always(func(g G.Play) {
		g.Say("some stuff")
	}))
}
