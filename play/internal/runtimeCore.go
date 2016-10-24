package internal

import (
	"github.com/ionous/mars/rtm"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/play/api"
	"io"
	"math/rand"
)

type PlayCore struct {
	Model  meta.Model
	Logger io.Writer
	Writer io.Writer
	api.LookupParents
	api.SaveLoad
	Frame api.EventFrame
	Rand  *rand.Rand // FIX: an interface part of config
	Rtm   *rtm.Rtm
}
