package internal

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// FUTURE: check mutation?
type ClassQuery struct {
	ga   *GameEventAdapter
	cls  ident.Id
	idx  int
	next meta.Instance
}

func (q *ClassQuery) HasNext() bool {
	return q.next != nil
}

func (q *ClassQuery) Next() (ret G.IObject) {
	if n := q.next; n != nil {
		q.idx, q.next = q.Advance()
		ret = q.ga.NewGameObject(n)
	} else {
		ret = NullObject("class query finished")
	}
	return
}

// Advance does not modify q.
func (q *ClassQuery) Advance() (ret int, inst meta.Instance) {
	ga, idx, clsid := q.ga, q.idx, q.cls
	for l := ga.Model.NumInstance(); idx < l; idx++ {
		gobj := ga.Model.InstanceNum(idx)
		if id := gobj.GetParentClass(); ga.Model.AreCompatible(id, clsid) {
			inst = gobj
			break
		}
	}
	return idx + 1, inst // explicit return to handle idx renaming
}
