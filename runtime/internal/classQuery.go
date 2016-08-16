package internal

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// FUTURE: check mutation?
type ClassQuery struct {
	ga    *GameEventAdapter
	cls   ident.Id
	exact bool
	idx   int
	next  meta.Instance
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
func (q *ClassQuery) Advance() (int, meta.Instance) {
	m, idx, clsid := q.ga.Model, q.idx, q.cls
	l := m.NumInstance()
	if q.exact {
		for ; idx < l; idx++ {
			gobj := m.InstanceNum(idx)
			if id := gobj.GetParentClass(); id == clsid {
				return idx + 1, gobj // explicit return to handle idx renaming
			}
		}
	} else {
		for ; idx < l; idx++ {
			gobj := m.InstanceNum(idx)
			if id := gobj.GetParentClass(); m.AreCompatible(id, clsid) {
				return idx + 1, gobj // explicit return to handle idx renaming
			}
		}
	}
	return l, nil
}
