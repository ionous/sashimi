package metal

import (
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

// uses pointers re: gopherjs
type numValue struct{ PanicValue }

func (p *numValue) SetNum(f float64) error {
	return p.SetGeneric(f)
}
func (p *numValue) GetNum() float64 {
	return p.getNum()
}

// uses pointers re: gopherjs
type textValue struct{ PanicValue }

func (p *textValue) GetText() string {
	return p.getString()
}
func (p *textValue) SetText(t string) error {
	return p.SetGeneric(t)
}

// uses pointers re: gopherjs
type pointerValue struct {
	PanicValue
}

func (p *pointerValue) GetObject() ident.Id {
	return p.getId()
}

func (p *pointerValue) SetObject(id ident.Id) (err error) {
	if !id.Equals(p.GetObject()) {
		if id.Empty() {
			err = p.SetGeneric(id)
		} else {
			if e := p.mdl.match(id, p.prop.Relates); e != nil {
				err = e
			} else if e := p.SetGeneric(id); e != nil {
				err = e
			}
		}
	}
	return
}

func (mdl *Metal) match(id ident.Id, relates ident.Id) (err error) {
	if target, ok := mdl.GetInstance(id); !ok {
		err = errutil.New("match, no such instance", id)
	} else if ok := mdl.AreCompatible(target.GetParentClass(), relates); !ok {
		err = errutil.New("match", target, "not compatible with", relates)
	}
	return
}
