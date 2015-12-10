package app

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/util/ident"
)

type View interface {
	Viewpoint() meta.Instance
	InView(meta.Instance) bool
	ChangedView(meta.Instance, ident.Id, meta.Instance) bool
	EnteredView(meta.Instance, ident.Id, meta.Instance) bool
}

// implement View for the standard rules
// memory games can keep this around, appengine can rebuild on the fly
type StandardView struct {
	player, status, room meta.Instance
	standard.ParentLookup
	visible map[ident.Id]StandardVisibilty
}

func NewStandardView(mdl meta.Model) (ret *StandardView) {
	if status, ok := mdl.GetInstance(ident.MakeId("status bar")); !ok {
		panic("couldnt find status")
	} else if player, ok := mdl.GetInstance(ident.MakeId("player")); !ok {
		panic("couldnt find player")
	} else {
		view := &StandardView{player: player, status: status, ParentLookup: standard.NewParentLookup(mdl)}
		view.ResetView(view.LookupRoot(player))
		ret = view
	}
	return ret
}

type StandardVisibilty int

const (
	VisibilityUnknown StandardVisibilty = iota
	Visible
	Invsible
)

func (v *StandardView) Viewpoint() meta.Instance {
	return v.player
}
func (v *StandardView) ResetView(curr meta.Instance) {
	v.visible = make(map[ident.Id]StandardVisibilty)
	v.room = curr
}

func (v *StandardView) ChangedView(gobj meta.Instance, prop ident.Id, next meta.Instance) (changed bool) {
	if which := standard.Containment[prop]; !which.Empty() {
		newRoom := v.LookupRoot(gobj)
		if newRoom != v.room {
			v.ResetView(newRoom)
			changed = true
		}
	}
	return
}

func (v *StandardView) EnteredView(gobj meta.Instance, prop ident.Id, next meta.Instance) (entered bool) {
	if v.room != nil {
		if which := standard.Containment[prop]; !which.Empty() {
			nowInRoom := next != nil && v.room == v.LookupRoot(next)
			id := gobj.GetId()
			if !nowInRoom {
				v.visible[id] = Invsible
			} else {
				if v.visible[id] != Visible {
					v.visible[id] = Visible
					entered = true
				}
			}
		}
	}
	return
}

func (v *StandardView) InView(i meta.Instance) bool {
	id := i.GetId()
	r := v.visible[id]
	if r == VisibilityUnknown {
		root := v.LookupRoot(i)
		// NOTE: player inv will have a root, so far as i know --
		// of the room -- because their parent is player, etc.
		// hrmm....
		if root == v.room || root == v.status {
			r = Visible
		} else {
			r = Invsible
		}
		v.visible[id] = r
	}
	return r == Visible
}
