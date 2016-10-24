package internal

import (
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// Implements E.ITarget for classes.
// The objects and classes form a uniform chain of targets
type ClassTarget struct {
	from     ObjectTarget
	class    ident.Id
	upObject meta.Instance
}

//
func (ct ClassTarget) Id() ident.Id {
	return ct.class
}

//
func (ct ClassTarget) Class() ident.Id {
	return ident.MakeId("class")
}

//
func (ct ClassTarget) String() string {
	return ct.class.String()
}

// Walk up the class hierarchy; when we reach the end, move to the next instance.
// (from E.ITarget)
func (ct ClassTarget) Parent() (ret E.ITarget, okay bool) {
	if cls, ok := ct.from.GetClass(ct.class); ok {
		ret = ClassTarget{ct.from, cls.GetParentClass(), ct.upObject}
		okay = true
	} else if next := ct.upObject; next != nil {
		ret = ObjectTarget{ct.from.Dispatch, next}
		okay = true
	}
	return
}

// Send an event to ct target.
// (from E.ITarget)
func (ct ClassTarget) TargetDispatch(evt E.IEvent) (err error) {
	return ct.from.DispatchEvent(evt, ct.class)
}
