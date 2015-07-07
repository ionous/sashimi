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
func (nouns NounNames) AddNameForId(name string, id StringId) {
	arr := nouns[name]
	arr = append(arr, id)
	nouns[name] = arr
}

//
// Try the passed function for every noun matching the passed name.
// Exits if the passed function returns true; returns the number of tries
//
func (nouns NounNames) Try(name string, try func(StringId) bool) (tries int, okay bool) {
	if names, ok := nouns[name]; ok {
		for _, id := range names {
			tries++
			if try(id) {
				okay = true
				break
			}
		}
	}
	return tries, okay
}
