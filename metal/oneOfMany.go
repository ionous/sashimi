package metal

import (
	"github.com/ionous/mars/rt"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/sbuf"
)

// oneOfManyProp uses ids to record the one side of a one-to-many, many-to-one relation:
// the many side stores its id in the "one-of" side.
type oneOfManyProp struct{ memProp }

func (p *oneOfManyProp) GetType() meta.PropertyType {
	return meta.ObjectProperty
}

// GetGeneric returns ObjectEval.
func (p *oneOfManyProp) GetGeneric() meta.Generic {
	id := p.getId()
	return rt.Reference(id)
}

// getId, oneOfManyProp stores ids
func (p *oneOfManyProp) getId() (ret ident.Id) {
	if v, ok := p.value.getValue(p.prop.Id); ok {
		// relations store ids
		if id, ok := v.(ident.Id); !ok {
			panic(errutil.New("get oneOfMany", p, "expected id, got", sbuf.Type{v}))
		} else {
			ret = id
		}
	}
	return
}

// SetGeneric oneOfMany expects rt.Object, noting rt.Object can be empty.
func (p *oneOfManyProp) SetGeneric(value meta.Generic) (err error) {
	if obj, ok := value.(rt.Object); !ok {
		err = errutil.New("set oneOfMany", p, "expected object, got", sbuf.Type{value})
	} else if ok := (obj.Instance == nil) || p.mdl.AreCompatible(obj.GetParentClass(), p.prop.Relates); !ok {
		err = errutil.New("set oneOfMany", p, "object", obj, "not compatible with", p.prop.Relates)
	} else {
		// relations store ids
		err = p.setValue(obj.GetId())
	}
	return
}
