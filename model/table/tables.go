package table

import "github.com/ionous/sashimi/util/ident"

type Tables map[ident.Id]*Table

func (t Tables) Clone() Tables {
	clone := make(Tables)
	for k, v := range t {
		pairs := make([]Pair, len(v.Pairs))
		copy(pairs, v.Pairs)
		clone[k] = &Table{pairs}
	}
	return clone
}
