package internal

import (
	"github.com/ionous/sashimi/runtime/api"
	"io"
	"math/rand"
)

type RuntimeCore struct {
	api.LookupParents
	api.Log
	api.SaveLoad
	//
	Frame  api.EventFrame
	Output api.Output
	Rand   *rand.Rand // FIX: an interface part of config
	Writer io.Writer
}
