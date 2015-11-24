package main

import (
	"flag"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/compiler/extract"
	//"github.com/ionous/sashimi/compiler/metal"
	//R "github.com/ionous/sashimi/runtime"
	"fmt"
	"github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard"
	"io"
	"os"
	"path"
)

// 1. extract the story callbacks and data into something like
// 	_gen/gofish/.....go
// the data requires a compiled script.
// until you come up with another plan, its not terrible if the current file has to compile, and its responsible, in a given mode, for extracting itself.

// 2. include them here
// 3. replace the model with a json decode of the data
// 4. replace calls with the extracted literal
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
		// cons := standard.GetConsole(opt)
		// defer cons.Close()
		// //
		// out := standard.NewStandardOutput(cons, writer)
		// parents := standard.ParentLookup{}
		// //
		// cfg := R.NewConfig().SetCalls(model.Calls).SetOutput(out).SetParentLookup(parents)
		// modelApi := metal.NewMetal(model.Model, make(metal.ObjectValueMap))
		// //
		// if g, e := cfg.NewGame(modelApi); e != nil {
		// 	panic(e)
		// } else if e := standard.PlayGame(cons, g); e != nil {
		// 	panic(e)
		// }
	}
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
