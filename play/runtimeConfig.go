package play

import (
	"github.com/ionous/mars/rtm"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/play/api"
	"github.com/ionous/sashimi/play/internal"
	"io"
	"math/rand"
	"os"
)

type Game struct {
	*internal.Game
}

type PlayConfig struct {
	core           internal.PlayCore
	writer, logger io.Writer
	parents        api.LookupParents
}

func NewConfig() *PlayConfig {
	return &PlayConfig{}
}

func (cfg *PlayConfig) MakeGame(model meta.Model) Game {
	// copy
	core := cfg.core
	writer := cfg.writer
	parents := cfg.parents
	logger := cfg.logger
	// defaults
	if core.Rand == nil {
		core.Rand = rand.New(rand.NewSource(1))
	}
	if writer == nil {
		writer = os.Stdout
	}
	if logger == nil {
		logger = writer
	}
	if core.Frame == nil {
		core.Frame = &noFrame{logger, nil}
	}
	if parents == nil {
		parents = api.NoParents{}
	}
	if core.SaveLoad == nil {
		core.SaveLoad = noSaveLoad{}
	}
	if core.Runtime == nil {
		run := rtm.NewRtmParents(model, parents)
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
	cfg.logger = log
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
