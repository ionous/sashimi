package runtime

import (
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

//
// Implements E.ITarget for classes.
// The objects and classes form a uniform chain of targets
//
type ClassTarget struct {
	host     ObjectTarget
	class    api.Class
	upObject api.Instance
}

//
func (ct ClassTarget) Id() ident.Id {
	return ct.class.GetId()
}

//
func (ct ClassTarget) Class() ident.Id {
	return ident.MakeId("class")
}

//
func (ct ClassTarget) String() string {
	return ct.class.GetId().String()
}

// Walk up the class hierarchy; when we reach the end, move to the next instance.
// (from E.ITarget)
func (ct ClassTarget) Parent() (ret E.ITarget, ok bool) {
	if parent := ct.class.GetParentClass(); parent != nil {
		ret = ClassTarget{ct.host, parent, ct.upObject}
		ok = true
	} else if next := ct.upObject; next != nil {
		ret = ObjectTarget{ct.host.game, next}
		ok = true
	}
	return ret, ok
}

// Send an event to ct target.
// (from E.ITarget)
func (ct ClassTarget) TargetDispatch(evt E.IEvent) (err error) {
	return ct.host.game.DispatchEvent(evt, ct.class.GetId())
}
