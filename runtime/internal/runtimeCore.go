package internal

import (
	"github.com/ionous/sashimi/runtime/api"
	"math/rand"
)

type RuntimeCore struct {
	api.LookupCallbacks
	api.LookupParents
	api.Log
	api.SaveLoad
	//
	Frame  api.EventFrame
	Output api.Output
	Rand   *rand.Rand // FIX: an interface part of config
}
