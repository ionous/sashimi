package live

import (
	G "github.com/ionous/sashimi/game"
	R "github.com/ionous/sashimi/runtime"
)

//
func CarriedNotWorn(obj G.IObject) (carrier G.IObject) {
	carried := false
	for _, wob := range []string{"owner"} {
		if p := obj.Object(wob); p.Exists() {
			carrier, carried = p, true
			break
		}
	}
	// so we arent pointing to nil, which cant easily be tested for. thanks go. :(
	if !carried {
		carrier = R.NullObject("CarriedNotWorn")
	}
	return
}

// FIX: note: this wouldnt work for something in a container
func Carrier(obj G.IObject) (carrier G.IObject) {
	carried := false
	for _, wob := range []string{"wearer", "owner"} {
		if p := obj.Object(wob); p.Exists() {
			carrier, carried = p, true
			break
		}
	}
	// so we arent pointing to nil, which cant easily be tested for. thanks go. :(
	if !carried {
		carrier = R.NullObject("Carrier")
	}
	return
}

// find the location ( the outermost room ) of the passed object
func Locate(obj G.IObject) (where G.IObject) {
	if p, ok := _location(obj); ok {
		where = p
	} else {
		where = R.NullObject("location")
	}
	return where
}
func _location(obj G.IObject) (parent G.IObject, okay bool) {
	if room := obj.Object("whereabouts"); room.Exists() {
		parent, okay = room, true
	} else if carrier := Carrier(obj); carrier.Exists() {
		parent, okay = _location(carrier)
	} else if supporter := obj.Object("support"); supporter.Exists() {
		parent, okay = _location(supporter)
	} else if container := obj.Object("enclosure"); container.Exists() {
		parent, okay = _location(container)
	}
	return
}

// find the outermost room or closed container containing the passed object
func Enclosure(obj G.IObject) (parent G.IObject) {
	if room := obj.Object("whereabouts"); room.Exists() {
		parent = room
	} else if carrier := Carrier(obj); carrier.Exists() {
		parent = Enclosure(carrier)
	} else if supporter := obj.Object("support"); supporter.Exists() {
		parent = Enclosure(supporter)
	} else if container := obj.Object("enclosure"); container.Exists() {
		if container.Is("open") {
			parent = Enclosure(container)
		} else {
			parent = container
		}
	} else {
		parent = R.NullObject("Enclosure")
	}
	return
}
