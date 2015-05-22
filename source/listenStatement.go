package source

import (
	G "github.com/ionous/sashimi/game"
)

//
// holds event callbacks
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
	Owner     string
	EventName string
	Callback  G.Callback
	Options   ListenOptions
}

func (this ListenStatement) Fields() ListenFields {
	return this.fields
}

func (this ListenStatement) Source() Code {
	return this.source
}

func (this ListenFields) Capture() bool {
	return this.Options&ListenCapture != 0
}

func (this ListenFields) TargetOnly() bool {
	return this.Options&ListenTargetOnly != 0
}

func (this ListenFields) RunAfter() bool {
	return this.Options&ListenRunAfter != 0
}
