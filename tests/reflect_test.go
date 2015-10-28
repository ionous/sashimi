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
