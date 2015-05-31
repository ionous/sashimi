package runtime

import (
	E "github.com/ionous/sashimi/event"
	M "github.com/ionous/sashimi/model"
)

//
// Implements E.ITarget for classes.
// The objects and classes form a uniform chain of targets
//
type ClassTarget struct {
	host     ObjectTarget
	class    *M.ClassInfo
	upObject *GameObject
}

//
func (this ClassTarget) String() string {
	return this.class.String()
}

//
// Walk up the class hierarchy; when we reach the end, move to the next instance.
// (from E.ITarget)
//
func (this ClassTarget) Parent() (ret E.ITarget, ok bool) {
	parent := this.class.Parent()
	if parent != nil {
		ret = ClassTarget{this.host, parent, this.upObject}
		ok = true
	} else if next := this.upObject; next != nil {
		ret = ObjectTarget{this.host.game, next}
		ok = true
	}
	return ret, ok
}

//
// Send an event to this target.
// (from E.ITarget)
//
func (this ClassTarget) Dispatch(evt E.IEvent) (err error) {
	if d, ok := this.host.game.Dispatchers.GetDispatcher(this.class); ok {
		err = d.Dispatch(evt)
	}
	return err
}
