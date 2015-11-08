package compiler

import (
	"github.com/ionous/sashimi/compiler/call"
	i "github.com/ionous/sashimi/compiler/internal"
	//	X "github.com/ionous/sashimi/model"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"io"
	"log"
)

type Config struct {
	Calls  call.Compiler
	Output io.Writer
}

// Compile script statements into a "model", a form usable by the runtime.
func (cfg Config) Compile(src S.Statements) (ret *M.Model, err error) {
	names := i.NewNameSource()
	rel := i.NewRelativeFactory(names.NewScope(nil))
	log := log.New(cfg.Output, "compling: ", log.Lshortfile)
	ctx := &i.Compiler{
		src, names.NewScope(nil),
		i.NewClassFactory(names, rel),
		i.NewInstanceFactory(names, log),
		rel,
		log,
		cfg.Calls,
	}
	if x, e := ctx.Compile(); e != nil {
		err = e
	} else {
		// m := &M.Model{}
		// for k, v := range x.Classes {
		// 	_, _ = k, v
		// }
		// for k, v := range x.Relations {
		// 	_, _ = k, v
		// }
		// for k, v := range x.Actions {
		// 	_, _ = k, v
		// }
		// for k, v := range x.Events {
		// 	_, _ = k, v
		// }
		// for i, v := range x.ParserActions {
		// 	_, _ = i, v
		// }
		// for k, v := range x.Instances {
		// 	_, _ = k, v
		// }
		// for i, v := range x.ActionHandlers {
		// 	_, _ = i, v
		// }
		// for i, v := range x.EventListeners {
		// 	_, _ = i, v
		// }
		// for k, v := range x.Tables.Tables {
		// 	_, _ = k, v
		// }
		// for _, v := range x.Generators {
		// }
		ret = x
	}
	return
}

type MemoryResult struct {
	Model *M.Model
	Calls call.MemoryStorage
}

func Compile(out io.Writer, src S.Statements) (res MemoryResult, err error) {
	calls := call.MakeMemoryStorage()
	cfg := Config{calls, out}
	if m, e := cfg.Compile(src); e != nil {
		err = e
	} else {
		res = MemoryResult{m, calls}
	}
	return
}
