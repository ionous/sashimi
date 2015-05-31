package table

type Table struct {
	table []Pair
}

func NewTable() Table {
	return Table{}
}

// might actually implement this in sqlite at some point
// ( along with the rest of the model data )
// not quite there yet.
type Pair struct {
	x, y string
}

func (this *Table) Add(x, y string) (index int) {
	index = this.Find(x, y)
	if index == 0 {
		this.table = append(this.table, Pair{x, y})
		index = len(this.table)
	}
	return index
}

func (this *Table) Find(x, y string) (index int) {
	for i, pair := range this.table {
		if pair.x == x && pair.y == y {
			index = i + 1
			break
		}
	}
	return index
}

func (this *Table) list(x string) (ret []string) {
	for _, pair := range this.table {
		if pair.x == x {
			ret = append(ret, pair.y)
		}
	}
	return ret
}

func (this *Table) listRev(y string) (ret []string) {
	for _, pair := range this.table {
		if pair.y == y {
			ret = append(ret, pair.x)
		}
	}
	return ret
}

func (this *Table) List(x string, rev bool) (ret []string) {
	if !rev {
		ret = this.list(x)
	} else {
		ret = this.listRev(x)
	}
	return ret
}

type pairTest func(x, y string) bool

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
