package runtime

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"

	"math/rand"

	"log"
)

type RuntimeConfig struct {
	core RuntimeCore
}

type RuntimeCore struct {
	calls   Callbacks
	frame   EventFrame
	output  IOutput
	parents ILookupParents
	Log
	rand *rand.Rand // FIX: an interface part of config
}

func NewConfig() *RuntimeConfig {
	return &RuntimeConfig{}
}

func (cfg RuntimeConfig) Finalize() RuntimeCore {
	core := cfg.core
	if core.rand == nil {
		core.rand = rand.New(rand.NewSource(1))
	}
	if core.Log == nil {
		log := log.New(logAdapter{core.output}, "game: ", log.Lshortfile)
		core.Log = LogAdapter{
			func(msg string) {
				log.Output(4, msg)
			}}
	}
	if core.frame == nil {
		core.frame = defaultFrame
	}
	if core.parents == nil {
		core.parents = parentLookup{}
	}
	return core
}

type parentLookup struct{}

func (parentLookup) GetParent(api.Model, api.Instance) (inst api.Instance, rel ident.Id, okay bool) {
	return
}

type logAdapter struct {
	output IOutput
}

func (log logAdapter) Write(p []byte) (n int, err error) {
	log.output.Log(string(p))
	return len(p), nil
}

type ILookupParents interface {
	GetParent(api.Model, api.Instance) (api.Instance, ident.Id, bool)
}

type Callbacks interface {
	// LookupCallback returns nil if not found.
	LookupCallback(ident.Id) (G.Callback, bool)
}

func (cfg *RuntimeConfig) SetCalls(calls Callbacks) *RuntimeConfig {
	cfg.core.calls = calls
	return cfg
}

// StartFrame and EndFrame should be merged into Output
// -- and they should be renamed: BeginEvent() EndEvent()
//*maybe* Target should be mapped into prototype
// Class should be removed from E.Target
// only: how do we know that a thing is a "class" and should get "Class" resource?
// could potentially send target type to startframe
// right now it seems redicoulous that the game decides that.
func (cfg *RuntimeConfig) SetFrame(e EventFrame) *RuntimeConfig {
	cfg.core.frame = e
	return cfg
}
func (cfg *RuntimeConfig) SetOutput(o IOutput) *RuntimeConfig {
	cfg.core.output = o
	return cfg
}
func (cfg *RuntimeConfig) SetParentLookup(l ILookupParents) *RuntimeConfig {
	cfg.core.parents = l
	return cfg
}
func (cfg *RuntimeConfig) SetLog(log Log) *RuntimeConfig {
	cfg.core.Log = log
	return cfg
}
func (cfg *RuntimeConfig) SetRand(rand *rand.Rand) *RuntimeConfig {
	cfg.core.rand = rand
	return cfg
}
