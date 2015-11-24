package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ionous/sashimi/_examples/fishgen"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/compiler/extract"
	"github.com/ionous/sashimi/compiler/metal"
	M "github.com/ionous/sashimi/compiler/model"
	G "github.com/ionous/sashimi/game"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/util/ident"
	"io"
	"os"
	"path"
)

// 5. try it out
// 6. replace internals with data store that occasionally flushes itself out.
//go:generate go extract %GOPATH/src/github.com/ionous/sashimi/_examples/stories/fishy.go
func main() {
	extract := flag.Bool("extract", false, "extract the game's data.")
	script.AddScript(stories.A_Day_For_Fresh_Sushi)
	script := script.InitScripts()
	opt := standard.ParseCommandLine()

	writer := standard.GetWriter(opt)
	if *extract {
		//go run gofish.go -extract
		if e := extractCalls("fishgen", script, writer); e != nil {
			panic(e)
		}
	} else {
		var model *M.Model
		if e := json.Unmarshal([]byte(fishgen.Data), &model); e != nil {
			panic(e)
		}
		calls := CodeCalls(fishgen.Callbacks)
		//
		cons := standard.GetConsole(opt)
		defer cons.Close()
		//
		out := standard.NewStandardOutput(cons, writer)
		parents := standard.ParentLookup{}
		//
		cfg := R.NewConfig().SetCalls(calls).SetOutput(out).SetParentLookup(parents)
		modelApi := metal.NewMetal(model, make(metal.ObjectValueMap))
		//
		if g, e := cfg.NewGame(modelApi); e != nil {
			panic(e)
		} else if e := standard.PlayGame(cons, g); e != nil {
			panic(e)
		}
	}
}

type CodeCalls map[ident.Id]G.Callback

func (m CodeCalls) LookupCallback(id ident.Id) (ret G.Callback, okay bool) {
	if r, ok := m[id]; !ok {
		panic(fmt.Sprintf("couldnt find callback %s", id))
	} else {
		ret, okay = r, ok
	}
	return
}

// 1. run go generate on this file
// 2. take the data, etc. and write that too.
// - in order for this to be the same, we need to include the files from there to start with. maybe itd be possible that the compile would fail the first time to generate the placeholder files, for now, creating them manually.

func extractCalls(name string, s *script.Script, trace io.Writer) (err error) {
	cx := extract.NewCallExtractor(name, "github.com/ionous/sashimi", trace)
	if model, e := s.CompileCalls(trace, cx); e != nil {
		err = e
	} else if f, e := os.Create(path.Join(name, "code.go")); e != nil {
		err = e
	} else {
		defer f.Close()
		io.WriteString(trace, fmt.Sprintf("writing %d snippets...", cx.Count()))
		if e := extract.WriteSnippets(f, cx); e != nil {
			err = e
		} else if f, e := os.Create(path.Join(name, "data.go")); e != nil {
			err = e
		} else {
			defer f.Close()
			extract.WriteJsonModel(f, name, model)
		}
	}
	return
}
