package framework

import (
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// implement View for the standard rules
// memory games can keep this around, appengine can rebuild on the fly
// NOTE! this has the unwatched model, making it dangerous to compare instances by value.
type StandardView struct {
	mdl          meta.Model
	player, room ident.Id
	ParentLookup
	visible map[ident.Id]StandardVisibilty
}

func NewStandardView(mdl meta.Model) (ret *StandardView) {
	if player, ok := mdl.GetInstance(ident.MakeId("player")); !ok {
		panic("couldnt find player")
	} else {
		view := &StandardView{mdl: mdl, ParentLookup: NewParentLookup(mdl)}
		view.ResetView(player.GetId(), view.LookupRoot(player).GetId())
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

func (v *StandardView) View() ident.Id {
	return v.room
}

func (v *StandardView) Viewer() ident.Id {
	return v.player
}

func (v *StandardView) String() (ret string) {
	return v.room.String()
}

func (v *StandardView) ResetView(player, room ident.Id) {
	v.visible = map[ident.Id]StandardVisibilty{player: Visible}
	v.player = player
	v.room = room
}

func (v *StandardView) ChangedView(gobj meta.Instance, prop ident.Id, next meta.Instance) (changed bool) {
	// only allows changes to the player containment
	// kind of silly, but there you are.
	// FUTURE: viewpoints...
	if gobj == nil || gobj.GetId() != v.player {
		panic("changing viewer not supported")
	}
	// supports hearing all viewer relations, so just limit to containment.
	if _, ok := Containment[prop]; ok {
		var newRoom ident.Id
		if next != nil {
			newRoom = v.LookupRoot(next).GetId()
		}
		if v.room != newRoom {
			v.ResetView(v.player, newRoom)
			changed = true
		}
	}
	return
}

func (v *StandardView) EnteredView(gobj meta.Instance, prop ident.Id, next meta.Instance) (entered bool) {
	if !v.room.Empty() {
		if which := Containment[prop]; !which.Empty() {
			nowInRoom := next != nil && v.room == v.LookupRoot(next).GetId()
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

func (v *StandardView) InView(i meta.Instance) (ret bool) {
	if !v.room.Empty() {
		id := i.GetId()
		r := v.visible[id]
		if r != VisibilityUnknown {
			/**/ //fmt.Println("StandardView: cached vis", id, r)
		} else {
			if v.mdl.AreCompatible(i.GetParentClass(), "globals") {
				r = Visible
			} else {
				root := v.LookupRoot(i).GetId()
				if root == v.room {
					r = Visible
				} else {
					r = Invsible
				}
				/**/ //fmt.Println("StandardView: stored vis", id, r == Visible, v.room, root)
			}
			v.visible[id] = r
		}
		ret = r == Visible
	}
	return
}
