package model

//
// Given a name find the associated resource.
//
type NounNames map[string]RankedStringIds

//
// A single name can map to multiple resources.
// Earlier resources are preferred in cases of conflict.
//
type RankedStringIds []StringId

//
// Add a name to resource mapping.
//
func (this NounNames) AddNameForId(name string, id StringId) {
	arr := this[name]
	arr = append(arr, id)
	this[name] = arr
}
