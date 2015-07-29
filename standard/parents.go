package standard

import (
	. "github.com/ionous/sashimi/game"
	R "github.com/ionous/sashimi/runtime"
)

// touchable ceiling, visibility ceiling (visibility level count)
func DirectParent(obj IObject) (parent IObject, where string) {
	for _, wse := range []string{"wearer", "owner", "whereabouts", "support", "enclosure"} {
		if p := obj.Object(wse); p.Exists() {
			parent, where = p, wse
			break
		}
	}
	return
}

//
func CarriedNotWorn(obj IObject) (carrier IObject, okay bool) {
	for _, wob := range []string{"owner"} {
		if p := obj.Object(wob); p.Exists() {
			carrier, okay = p, true
			break
		}
	}
	return
}

//
func Carrier(obj IObject) (carrier IObject, okay bool) {
	for _, wob := range []string{"wearer", "owner"} {
		if p := obj.Object(wob); p.Exists() {
			carrier, okay = p, true
			break
		}
	}
	return
}

// find the location ( the outermost room ) of the passed object
func Locate(obj IObject) (where IObject) {
	if p, ok := _location(obj); ok {
		where = p
	} else {
		where = R.NullObject("location")
	}
	return where
}
func _location(obj IObject) (parent IObject, okay bool) {
	if room := obj.Object("whereabouts"); room.Exists() {
		parent, okay = room, true
	} else if carrier, ok := Carrier(obj); ok {
		parent, okay = _location(carrier)
	} else if supporter := obj.Object("support"); supporter.Exists() {
		parent, okay = _location(supporter)
	} else if container := obj.Object("enclosure"); container.Exists() {
		parent, okay = _location(container)
	}
	return
}

// //
// // find the direct parent ( moving towards the outermost room or closed container ) of the passed object
// func ParentByEnclosure(obj IObject) (parent IObject, where string) {
// 	if room, ok := obj.Object("whereabouts"); ok {
// 		parent, wheres = room, "whereabouts"
// 	} else if carrier, ok := Carrier(obj); ok != "" {
// 		parent, where = carrier, ok
// 	} else if supporter, ok := obj.Object("support"); ok {
// 		parent, where = supporter, "support"
// 	} else if container, ok := obj.Object("enclosure"); ok {
// 		if container.Is("open") {
// 			parent, where = container, "enclosure"
// 		}
// 	}
// 	return
// }

//
// find the outermost room or closed container containing the passed object
func Enclosure(obj IObject) (parent IObject, okay bool) {
	if room := obj.Object("whereabouts"); room.Exists() {
		parent, okay = room, true
	} else if carrier, ok := Carrier(obj); ok {
		parent, okay = Enclosure(carrier)
	} else if supporter := obj.Object("support"); supporter.Exists() {
		parent, okay = Enclosure(supporter)
	} else if container := obj.Object("enclosure"); container.Exists() {
		if container.Is("open") {
			parent, okay = Enclosure(container)
		} else {
			parent, okay = container, true
		}
	}
	return
}
