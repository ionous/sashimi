package internal

import (
	"github.com/ionous/mars/rt"
	// "github.com/ionous/mars/rtm"
	"github.com/ionous/sashimi/play/api"
	//	"io"
	"math/rand"
)

type PlayCore struct {
	rt.Runtime
	api.SaveLoad
	Frame api.EventFrame
	Rand  *rand.Rand // FIX: an interface part of config
}