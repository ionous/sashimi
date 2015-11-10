package parser

import "github.com/ionous/sashimi/util/ident"

type Comprehensions map[ident.Id]*Comprehension

// NewComprehension adds a pattern set.
// Name must be unique ( used to help with error-handling and auto-documentation. )
func (comps Comprehensions) NewComprehension(id ident.Id) (
	ret *Comprehension,
	err error,
) {
	if id.Empty() {
		err = InvalidComprehension(id)
	} else if _, exists := comps[id]; exists {
		err = DuplicateComprehension(id)
	} else {
		comp := &Comprehension{id, nil}
		comps[id] = comp
		ret = comp
	}
	return ret, err
}
