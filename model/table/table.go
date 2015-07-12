package table

import "github.com/ionous/sashimi/util/ident"

type Table struct {
	table []Pair
}

func NewTable() *Table {
	return &Table{}
}

// might actually implement this in sqlite at some point
// ( along with the rest of the model data )
// not quite there yet.
type Pair struct {
	x, y ident.Id
}

func (this *Table) Add(x, y ident.Id) (index int) {
	index = this.Find(x, y)
	if index == 0 {
		this.table = append(this.table, Pair{x, y})
		index = len(this.table)
	}
	return index
}

func (this *Table) Find(x, y ident.Id) (index int) {
	for i, pair := range this.table {
		if pair.x == x && pair.y == y {
			index = i + 1
			break
		}
	}
	return index
}

func (this *Table) list(x ident.Id) (ret []ident.Id) {
	for _, pair := range this.table {
		if pair.x == x {
			ret = append(ret, pair.y)
		}
	}
	return ret
}

func (this *Table) listRev(y ident.Id) (ret []ident.Id) {
	for _, pair := range this.table {
		if pair.y == y {
			ret = append(ret, pair.x)
		}
	}
	return ret
}

func (this *Table) List(x ident.Id, rev bool) (ret []ident.Id) {
	if !rev {
		ret = this.list(x)
	} else {
		ret = this.listRev(x)
	}
	return ret
}

type pairTest func(x, y ident.Id) bool

func (this *Table) Remove(pairTest pairTest) (removed int) {
	t := this.table
	for i := 0; i < len(t); i++ {
		pair := &t[i]
		if pairTest(pair.x, pair.y) {
			end := len(t) - 1
			*pair = t[end]
			t = t[:end]
			removed = removed + 1
		}
	}
	this.table = t
	return removed
}
