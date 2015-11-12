package event

import "github.com/ionous/sashimi/util/ident"

// Proc implements IEvent
type Proc struct {
	msg           *Message
	currentTarget ITarget
	phase         Phase
	cancelled     bool
	stopMore      bool
	stopNow       bool
	path          PathList
	target        ITarget
}

//
func (evt *Proc) sendToTarget(loc ITarget) (err error) {
	evt.currentTarget = loc
	return evt.currentTarget.TargetDispatch(evt)
}

//
func (evt *Proc) Id() ident.Id {
	return evt.msg.Id
}

//
func (evt *Proc) Data() interface{} {
	return evt.msg.Data
}

//
func (evt *Proc) Bubbles() bool {
	return !evt.msg.CaptureOnly
}

//
func (evt *Proc) Cancelable() bool {
	return !evt.msg.CantCancel
}

//
func (evt *Proc) DefaultBlocked() bool {
	return evt.cancelled
}

//
func (evt *Proc) Target() ITarget {
	return evt.target
}

//
func (evt *Proc) Path() PathList {
	return evt.path
}

//
func (evt *Proc) Phase() Phase {
	return evt.phase
}

//
func (evt *Proc) CurrentTarget() ITarget {
	return evt.currentTarget
}

//
func (evt *Proc) PreventDefault() bool {
	if !evt.msg.CantCancel {
		evt.cancelled = true
	}
	return evt.cancelled
}

//
func (evt *Proc) StopPropagation() {
	evt.stopMore = true
}

//
func (evt *Proc) StopImmediatePropagation() {
	evt.stopNow = true
	evt.stopMore = true
}
