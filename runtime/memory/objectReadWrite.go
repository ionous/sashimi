package memory

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

// returned by object/relation array b/c we cant mutate individual values
type objectReadValue struct {
	panicValue
	currentVal ident.Id
}

func (p objectReadValue) GetObject() (ret ident.Id) {
	return p.currentVal
}

// the one side of a many-to-one, one-to-one, or one-to-many relation.
func singleValue(p *propBase) api.Value {
	rel := p.prop.(M.RelativeProperty)
	objs := p.mdl.getObjects(p.src, rel.Relation, rel.IsRev)
	var v ident.Id
	if len(objs) > 0 {
		v = objs[0]
	}
	return objectWriteValue{objectReadValue{panicValue{p}, v}}
}

// returned by singleValue
type objectWriteValue struct {
	objectReadValue
}

//
func (p objectWriteValue) SetObject(id ident.Id) (err error) {
	if !id.Empty() {
		err = p.mdl.canAppend(id, p.src, p.prop.(M.RelativeProperty))
	}
	if err == nil {
		p.mdl.clearValues(p.src, p.prop.(M.RelativeProperty))
		p.mdl.appendObject(id, p.src, p.prop.(M.RelativeProperty))
		p.currentVal = id
	}
	return err
}
