package model

import "github.com/ionous/sashimi/util/ident"

// NounNames maps names typed by the player to ids.
type NounNames map[string]RankedStringIds

// RankedStringIds allows a single name to map to multiple resources.
// Earlier resources are preferred in cases of conflict.
type RankedStringIds []ident.Id

// AddNameForId
func (nouns NounNames) AddNameForId(name string, id ident.Id) {
	arr := nouns[name]
	arr = append(arr, id)
	nouns[name] = arr
}

// Try the passed function for every noun matching the passed name.
// Exits if the passed function returns true; returns the number of tries.
func (nouns NounNames) Try(name string, try func(ident.Id) bool) (tries int, okay bool) {
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
