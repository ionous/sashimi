package metal

import (
	"fmt"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

var _ = fmt.Println

type objectList struct {
	panicValue
	objs []ident.Id
}

func manyValue(p *propBase) meta.Values {
	objs := p.mdl.getObjects(p.src, p.prop.Relation, p.prop.IsRev)
	return &objectList{panicValue{p}, objs}
}

func (p objectList) NumValue() int {
	return len(p.objs)
}

func (p objectList) ValueNum(i int) meta.Value {
	return objectReadValue{p.panicValue, p.objs[i]}
}

func (p *objectList) ClearValues() (err error) {
	p.mdl.clearValues(p.src, p.prop)
	p.objs = nil
	return
}

func (objectList) AppendNum(float32) error {
	panic("not implemented")
}
func (p objectList) AppendText(string) error {
	panic("not implemented")
}

func (p *objectList) AppendObject(id ident.Id) (err error) {
	if e := p.mdl.canAppend(id, p.src, p.prop); e != nil {
		err = e
	} else {
		p.mdl.appendObject(id, p.src, p.prop)
		p.objs = append(p.objs, id)
	}
	return
}

func (mdl Metal) clearValues(src ident.Id, rel *M.PropertyModel) {
	table := mdl.getTable(rel.Relation)
	isRev := rel.IsRev
	table.Remove(func(x, y ident.Id) bool {
		var test, other ident.Id
		if isRev {
			test, other = y, x
		} else {
			test, other = x, y
		}
		_ = other
		return src == test
	})
}
