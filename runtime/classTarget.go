package runtime

import (
	E "github.com/ionous/sashimi/event"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
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
func (ct ClassTarget) Id() ident.Id {
	return ct.class.Id()
}

//
func (ct ClassTarget) Class() ident.Id {
	return ident.MakeId("class")
}

//
func (ct ClassTarget) String() string {
	return ct.class.String()
}

//
// Walk up the class hierarchy; when we reach the end, move to the next instance.
// (from E.ITarget)
//
func (ct ClassTarget) Parent() (ret E.ITarget, ok bool) {
	parent := ct.class.Parent()
	if parent != nil {
		ret = ClassTarget{ct.host, parent, ct.upObject}
		ok = true
	} else if next := ct.upObject; next != nil {
		ret = ObjectTarget{ct.host.game, next}
		ok = true
	}
	return ret, ok
}

//
// Send an event to ct target.
// (from E.ITarget)
//
func (ct ClassTarget) Dispatch(evt E.IEvent) (err error) {
	if d, ok := ct.host.game.Dispatchers.GetDispatcher(ct.class.Id()); ok {
		err = d.Dispatch(evt)
	}
	return err
}
