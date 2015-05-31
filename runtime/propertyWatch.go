package runtime

//
// Callback interface for hearing when an object's property has changed.
// FIX? I think better would be a property proivder interface
// and then the user/implementer code could implement the watch inside its own set.
//
type IPropertyChanged interface {
	PropertyChanged(objectId, propertyId string, prev, value interface{})
}

//
// Set of property changed callbacks.
//
type PropertyWatchers struct {
	arr []IPropertyChanged
}

//
// Sends a notification that an object property has changed.
//
func (this *PropertyWatchers) Notify(
	object string,
	property string,
	prev interface{},
	value interface{}) {
	for _, el := range this.arr {
		el.PropertyChanged(object, property, prev, value)
	}
}

//
// Start listening for property changes.
//
func (this *PropertyWatchers) AddWatcher(p IPropertyChanged) {
	this.arr = append(this.arr, p)
}

//
// Stop listening for property changes.
//
func (this *PropertyWatchers) RemoveWatcher(p IPropertyChanged) (found bool) {
	a := this.arr
	if l := len(a); l > 0 {
		for i, el := range a {
			if el == p {
				// slice and allow the swapped item to garbage collect
				a[i], a[l-1], a = a[l-1], nil, a[:l-1]
				found = true
				break
			}
		}
		this.arr = a
	}
	return found
}
