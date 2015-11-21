package runtime

import (
	"fmt"
	"github.com/ionous/sashimi/compiler/metal"
	M "github.com/ionous/sashimi/compiler/model"
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/runtime/internal"
	"log"
	"math/rand"
)

type RuntimeConfig struct {
	core internal.RuntimeCore
}

func NewConfig() *RuntimeConfig {
	return &RuntimeConfig{}
}

func (cfg RuntimeConfig) NewGame(model *M.Model) (Game, error) {
	core := cfg.Finalize()
	modelApi := metal.NewMetal(model, make(metal.ObjectValueMap))
	return Game{modelApi, internal.NewGame(core, modelApi)}, nil
}

func (cfg RuntimeConfig) Finalize() internal.RuntimeCore {
	core := cfg.core
	if core.Rand == nil {
		core.Rand = rand.New(rand.NewSource(1))
	}
	if core.Log == nil {
		log := log.New(logAdapter{core.Output}, "game: ", log.Lshortfile)
		core.Log = LogAdapter{
			func(msg string) {
				log.Output(4, msg)
			}}
	}
	if core.Frame == nil {
		core.Frame = defaultFrame
	}
	if core.LookupParents == nil {
		core.LookupParents = parentLookup{}
	}
	return core
}

type parentLookup struct{}

func (parentLookup) LookupParent(meta.Model, meta.Instance) (inst meta.Instance, rel meta.Property, okay bool) {
	return
}

type logAdapter struct {
	output api.Output
}

func (log logAdapter) Write(p []byte) (n int, err error) {
	log.output.Log(string(p))
	return len(p), nil
}

func (cfg *RuntimeConfig) SetCalls(calls api.LookupCallbacks) *RuntimeConfig {
	cfg.core.LookupCallbacks = calls
	return cfg
}

// StartFrame and EndFrame should be merged into Output
// -- and they should be renamed: BeginEvent() EndEvent()
//*maybe* Target should be mapped into prototype
// Class should be removed from E.Target
// only: how do we know that a thing is a "class" and should get "Class" resource?
// could potentially send target type to startframe
// right now it seems redicoulous that the game decides that.
func (cfg *RuntimeConfig) SetFrame(e api.EventFrame) *RuntimeConfig {
	cfg.core.Frame = e
	return cfg
}
func (cfg *RuntimeConfig) SetOutput(o api.Output) *RuntimeConfig {
	cfg.core.Output = o
	return cfg
}
func (cfg *RuntimeConfig) SetParentLookup(l api.LookupParents) *RuntimeConfig {
	cfg.core.LookupParents = l
	return cfg
}
func (cfg *RuntimeConfig) SetLog(log api.Log) *RuntimeConfig {
	cfg.core.Log = log
	return cfg
}
func (cfg *RuntimeConfig) SetRand(rand *rand.Rand) *RuntimeConfig {
	cfg.core.Rand = rand
	return cfg
}

func defaultFrame(E.ITarget, *E.Message) func() {
	return func() {}
}

type LogAdapter struct {
	print func(s string)
}

func (log LogAdapter) Printf(format string, v ...interface{}) {
	log.print(fmt.Sprintf(format, v...))
}

func (log LogAdapter) Println(v ...interface{}) {
	log.print(fmt.Sprintln(v...))
}
