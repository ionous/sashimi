package model

import "github.com/ionous/sashimi/util/ident"

// Aliases maps names typed by the player to ids.
type Aliases map[string]RankedStringIds

// RankedStringIds allows a single name to map to multiple resources.
// Earlier resources are preferred in cases of conflict.
type RankedStringIds []ident.Id

// Try the passed function for every noun matching the passed name.
// Exits if the passed function returns true; returns the number of tries.
func (a Aliases) Try(name string, try func(ident.Id) bool) (tries int, okay bool) {
	if names, ok := a[name]; ok {
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
