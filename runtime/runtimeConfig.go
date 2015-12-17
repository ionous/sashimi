package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/runtime/internal"
	"io"
	"log"
	"math/rand"
	"strings"
)

type RuntimeConfig struct {
	core internal.RuntimeCore
}

func NewConfig() *RuntimeConfig {
	return &RuntimeConfig{}
}

//modelApi := metal.NewMetal(model, make(metal.ObjectValueMap))
func (cfg RuntimeConfig) NewGame(m meta.Model) (Game, error) {
	core := cfg.Finalize()
	return Game{m, internal.NewGame(core, m)}, nil
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
		core.Frame = &defaultFrame{core.Log, nil}
	}
	if core.LookupParents == nil {
		core.LookupParents = noParents{}
	}
	if core.SaveLoad == nil {
		core.SaveLoad = noSaveLoad{}
	}
	return core
}

type noSaveLoad struct{}

func (noSaveLoad) Save(io.Writer) error {
	return fmt.Errorf("not implemented")
}
func (noSaveLoad) Load(io.Reader) error {
	return fmt.Errorf("not implemented")
}

type noParents struct{}

func (noParents) LookupParent(meta.Instance) (inst meta.Instance, rel meta.Property, okay bool) {
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
func (cfg *RuntimeConfig) SetSaveLoad(s api.SaveLoad) *RuntimeConfig {
	cfg.core.SaveLoad = s
	return cfg
}

type defaultFrame struct {
	log   api.Log
	parts []string
}

func (d *defaultFrame) BeginEvent(_ E.ITarget, path E.PathList, msg *E.Message) api.IEndEvent {
	d.parts = append(d.parts, msg.String())
	fullName := strings.Join(d.parts, "/")
	d.log.Printf("sending `%s` to: %s.", fullName, path)
	return d
}

func (d *defaultFrame) FlushFrame() {
}

func (d *defaultFrame) EndEvent() {
	d.parts = d.parts[0 : len(d.parts)-1]
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
