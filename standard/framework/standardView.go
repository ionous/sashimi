package framework

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// implement View for the standard rules
// memory games can keep this around, appengine can rebuild on the fly
type StandardView struct {
	mdl          meta.Model
	player, room meta.Instance
	ParentLookup
	visible map[ident.Id]StandardVisibilty
}

func NewStandardView(mdl meta.Model) (ret *StandardView) {
	if player, ok := mdl.GetInstance(ident.MakeId("player")); !ok {
		panic("couldnt find player")
	} else {
		view := &StandardView{mdl: mdl, ParentLookup: NewParentLookup(mdl)}
		view.ResetView(player, view.LookupRoot(player))
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
func (v *StandardView) ResetView(player, room meta.Instance) {
	v.visible = map[ident.Id]StandardVisibilty{player.GetId(): Visible}
	v.player = player
	v.room = room
}

func (v *StandardView) ChangedView(gobj meta.Instance, prop ident.Id, next meta.Instance) (changed bool) {
	// did the player property change
	if gobj == v.player {
		// and is it contaiment
		if _, ok := Containment[prop]; ok {
			if next != nil {
				next = v.LookupRoot(next)
			}
			if next != v.room {
				v.ResetView(gobj, next)
				changed = true
			}
		}
	}
	return
}

func (v *StandardView) EnteredView(gobj meta.Instance, prop ident.Id, next meta.Instance) (entered bool) {
	if v.room != nil {
		if which := Containment[prop]; !which.Empty() {
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
		if v.mdl.AreCompatible(i.GetParentClass().GetId(), "stories") ||
			v.mdl.AreCompatible(i.GetParentClass().GetId(), "status-bar-instances") {
			r = Visible
		} else {
			root := v.LookupRoot(i)
			// NOTE: player inv will have a root, so far as i know --
			// of the room -- because their parent is player, etc.
			// hrmm....
			if root == v.room {
				r = Visible
			} else {
				r = Invsible
			}
		}
		v.visible[id] = r
	}
	return r == Visible
}
