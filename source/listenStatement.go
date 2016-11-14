package source

import (
	"github.com/ionous/mars/rt"
	. "github.com/ionous/sashimi/source/types"
)

// ListenStatement holds event callbacks.
type ListenStatement struct {
	fields ListenFields
	source Code
}

type ListenOptions int

const (
	ListenBubble  ListenOptions = 0
	ListenCapture               = 1 << iota
	ListenTargetOnly
	ListenRunAfter
)

type ListenFields struct {
	Owner   string
	Event   string
	Calls   []rt.Execute
	Options ListenOptions
}

func (ts ListenStatement) Fields() ListenFields {
	return ts.fields
}

func (ts ListenStatement) Source() Code {
	return ts.source
}

func (ts ListenFields) Captures() bool {
	return ts.Options&ListenCapture != 0
}

func (ts ListenFields) OnlyTargets() bool {
	return ts.Options&ListenTargetOnly != 0
}

func (ts ListenFields) RunsAfter() bool {
	return ts.Options&ListenRunAfter != 0
}
