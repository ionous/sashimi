package play

import (
	"github.com/ionous/mars/rtm"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/play/api"
	"github.com/ionous/sashimi/play/internal"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
)

type Game struct {
	*internal.Game
}

type PlayConfig struct {
	core    internal.PlayCore
	writer  io.Writer
	parents api.LookupParents
}

func NewConfig() *PlayConfig {
	return &PlayConfig{}
}

func (cfg *PlayConfig) MakeGame(model meta.Model) Game {
	// copy
	core := cfg.core
	writer := cfg.writer
	// defaults
	if core.Rand == nil {
		core.Rand = rand.New(rand.NewSource(1))
	}
	if writer == nil {
		writer = os.Stdout
	}
	if core.Logger == nil {
		core.Logger = ioutil.Discard
	}
	if core.Frame == nil {
		core.Frame = &noFrame{core.Logger, nil}
	}
	if core.Parents == nil {
		core.Parents = &api.ParentHolder{}
	}
	if core.SaveLoad == nil {
		core.SaveLoad = noSaveLoad{}
	}
	if core.Runtime == nil {
		run := rtm.NewRtm(model)
		core.Runtime = run.Runtime()
		core.Runtime.PushOutput(writer)
	}
	return Game{internal.NewGame(core)}
}

func (cfg *PlayConfig) SetFrame(e api.EventFrame) *PlayConfig {
	cfg.core.Frame = e
	return cfg
}

func (cfg *PlayConfig) SetWriter(out io.Writer) *PlayConfig {
	cfg.writer = out
	return cfg
}

func (cfg *PlayConfig) SetParentLookup(l api.LookupParents) *PlayConfig {
	cfg.parents = l
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
