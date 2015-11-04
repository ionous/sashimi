package runtime

import (
	"github.com/ionous/sashimi/util/ident"
)

// PropertyChange for hearing when an object's data has changed.
type PropertyChange interface {
	NumChange(obj *GameObject, prop ident.Id, prev, next float32)
	TextChange(obj *GameObject, prop ident.Id, prev, next string)
	StateChange(obj *GameObject, prop ident.Id, prev, next ident.Id)
	ReferenceChange(obj *GameObject, prop, other ident.Id, prev, next *GameObject)
}

// PropertyWatchers provides a collection of PropertyChange interfaces
type PropertyWatchers struct {
	arr []PropertyChange
}

// VisitWatchers
func (w PropertyWatchers) VisitWatchers(cb func(PropertyChange)) {
	for _, el := range w.arr {
		cb(el)
	}
}

// AddWatcher starts listening for property changes.
func (w *PropertyWatchers) AddWatcher(p PropertyChange) {
	w.arr = append(w.arr, p)
}

// RemoveWatcher stops listening for property changes.
func (w *PropertyWatchers) RemoveWatcher(p PropertyChange) (found bool) {
	a := w.arr
	if l := len(a); l > 0 {
		for i, el := range a {
			if el == p {
				// slice and allow the swapped item to garbage collect
				a[i], a[l-1], a = a[l-1], nil, a[:l-1]
				found = true
				break
			}
		}
		w.arr = a
	}
	return found
}
