package main

import (
	"appengine/aetest"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ionous/sashimi/_examples/fishgen"
	"github.com/ionous/sashimi/_examples/stories"
	D "github.com/ionous/sashimi/appengine/datastore"
	"github.com/ionous/sashimi/compiler/extract"
	M "github.com/ionous/sashimi/compiler/model"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/metal"
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/script/backend"
	"github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/util/ident"
	"io"
	"os"
	"path"
)

func main() {
	script.AddScript(stories.A_Day_For_Fresh_Sushi)
	script := script.InitScripts()
	//
	extract := flag.Bool("extract", false, "extract the game's data.")
	dstest := flag.Bool("ds", false, "use datastore to run the game.")

	opt := standard.ParseCommandLine()
	//
	writer := standard.GetWriter(opt)
	if *extract {
		if e := extractCalls("fishgen", script, writer); e != nil {
			panic(e)
		}
	} else {
		if model, calls, e := getModelCalls(); e != nil {
			panic(e)
		} else {
			var modelApi meta.Model
			var update func()
			if !*dstest {
				modelApi = metal.NewMetal(model, make(metal.ObjectValueMap))
			} else {
				if ctx, e := aetest.NewContext(nil); e != nil {
					panic(e)
				} else {
					defer ctx.Close()

					ds := D.NewModelStore(ctx, model, nil)
					modelApi = ds.Model()

					// every frame flush ( save ) the cache, and empty it.
					update = func() {
						if e := ds.Flush(); e != nil {
							panic(e)
						}
						ds.Drop()
					}
				}
			}

			//
			cons := standard.GetConsole(opt)
			defer cons.Close()
			//
			out := standard.NewStandardOutput(cons, writer)
			parents := standard.NewParentLookup(modelApi)
			//
			cfg := R.NewConfig().SetCalls(calls).SetOutput(out).SetParentLookup(parents)
			//
			g := cfg.MakeGame(modelApi)
			if e := standard.PlayGameUpdate(cons, g, update); e != nil {
				panic(e)
			}
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

func getModelCalls() (*M.Model, api.LookupCallbacks, error) {
	var model *M.Model
	e := json.Unmarshal([]byte(fishgen.Data), &model)
	return model, CodeCalls(fishgen.Callbacks), e
}

func extractCalls(name string, s *backend.Script, trace io.Writer) (err error) {
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
