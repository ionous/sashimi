package event

// implements IEvent
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
func (this *Proc) sendToTarget(loc ITarget) (err error) {
	this.currentTarget = loc
	return this.currentTarget.Dispatch(this)
}

//
func (this *Proc) Name() string {
	return this.msg.Name
}

//
func (this *Proc) Data() interface{} {
	return this.msg.Data
}

//
func (this *Proc) Bubbles() bool {
	return !this.msg.CaptureOnly
}

//
func (this *Proc) Cancelable() bool {
	return !this.msg.CantCancel
}

//
func (this *Proc) DefaultBlocked() bool {
	return this.cancelled
}

//
func (this *Proc) Target() ITarget {
	return this.target
}

//
func (this *Proc) Path() PathList {
	return this.path
}

//
func (this *Proc) Phase() Phase {
	return this.phase
}

//
func (this *Proc) CurrentTarget() ITarget {
	return this.currentTarget
}

//
func (this *Proc) PreventDefault() bool {
	if !this.msg.CantCancel {
		this.cancelled = true
	}
	return this.cancelled
}

//
func (this *Proc) StopPropagation() {
	this.stopMore = true
}

//
func (this *Proc) StopImmediatePropagation() {
	this.stopNow = true
	this.stopMore = true
}
