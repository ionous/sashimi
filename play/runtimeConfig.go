package play

import (
	"github.com/ionous/mars/rtm"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/play/api"
	"github.com/ionous/sashimi/play/internal"
	"io"
	"math/rand"
)

type Game struct {
	*internal.Game
}

type PlayConfig struct {
	core internal.PlayCore
}

func NewConfig() *PlayConfig {
	return &PlayConfig{}
}

func (cfg *PlayConfig) MakeGame(model meta.Model) Game {
	// copy
	core := cfg.core
	core.Model = model
	// defaults
	if core.Rand == nil {
		core.Rand = rand.New(rand.NewSource(1))
	}
	if core.Logger == nil {
		core.Logger = core.Writer
	}
	if core.Frame == nil {
		core.Frame = &noFrame{core.Logger, nil}
	}
	if core.LookupParents == nil {
		core.LookupParents = noParents{}
	}
	if core.SaveLoad == nil {
		core.SaveLoad = noSaveLoad{}
	}
	if core.Rtm == nil {
		core.Rtm = rtm.NewRtm(model)

		// func (run *Mars) LookupParent(inst meta.Instance) (ret meta.Instance, rel meta.Property, okay bool) {
		// 	return run.ga.Game.LookupParent(inst)
		// }

		core.Rtm.PushOutput(core.Writer)
	}
	return Game{internal.NewGame(core)}
}

func (cfg *PlayConfig) SetFrame(e api.EventFrame) *PlayConfig {
	cfg.core.Frame = e
	return cfg
}
func (cfg *PlayConfig) SetWriter(out io.Writer) *PlayConfig {
	cfg.core.Writer = out
	return cfg
}
func (cfg *PlayConfig) SetParentLookup(l api.LookupParents) *PlayConfig {
	cfg.core.LookupParents = l
	return cfg
}
func (cfg *PlayConfig) SetLogger(log io.Writer) *PlayConfig {
	cfg.core.Logger = log
	return cfg
}
func (cfg *PlayConfig) SetRand(rand *rand.Rand) *PlayConfig {
	cfg.core.Rand = rand
	return cfg
}
func (cfg *PlayConfig) SetSaveLoad(s api.SaveLoad) *PlayConfig {
	cfg.core.SaveLoad = s
	return cfg
}
