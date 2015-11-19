package table

import "github.com/ionous/sashimi/util/ident"

type Table struct {
	Pairs []Pair
}

func NewTable() *Table {
	return &Table{}
}

// might actually implement t in sqlite at some point
// ( along with the rest of the model data )
// not quite there yet.
type Pair struct {
	X, Y ident.Id
}

func (t *Table) Add(x, y ident.Id) (index int) {
	index = t.Find(x, y)
	if index == 0 {
		t.Pairs = append(t.Pairs, Pair{x, y})
		index = len(t.Pairs)
	}
	return index
}

func (t *Table) Find(x, y ident.Id) (index int) {
	for i, pair := range t.Pairs {
		if ident.Compare(pair.X, x) == 0 && ident.Compare(pair.Y, y) == 0 {
			index = i + 1
			break
		}
	}
	return index
}

func (t *Table) list(x ident.Id) (ret []ident.Id) {
	for _, pair := range t.Pairs {
		if ident.Compare(pair.X, x) == 0 {
			ret = append(ret, pair.Y)
		}
	}
	return ret
}

func (t *Table) listRev(y ident.Id) (ret []ident.Id) {
	for _, pair := range t.Pairs {
		if ident.Compare(pair.Y, y) == 0 {
			ret = append(ret, pair.X)
		}
	}
	return ret
}

func (t *Table) List(x ident.Id, rev bool) (ret []ident.Id) {
	if !rev {
		ret = t.list(x)
	} else {
		ret = t.listRev(x)
	}
	return ret
}

type pairTest func(x, y ident.Id) bool

func (t *Table) Remove(pairTest pairTest) (removed int) {
	pairs := t.Pairs
	for i := 0; i < len(pairs); i++ {
		pair := &pairs[i]
		if pairTest(pair.X, pair.Y) {
			end := len(pairs) - 1
			*pair = pairs[end]
			pairs = pairs[:end]
			removed = removed + 1
		}
	}
	t.Pairs = pairs
	return removed
}
