package metal

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
)

type numValue struct{ panicValue }

func (p numValue) SetNum(f float32) error {
	return p.set(f)
}
func (p numValue) GetNum() float32 {
	return p.getNum()
}

type numValues struct{ panicValue }

func (p numValues) SetNum(f float32) error {
	return p.set(f)
}
func (p numValues) GetNum() float32 {
	return p.getNum()
}

type textValue struct{ panicValue }

func (p textValue) GetText() string {
	return p.getString()
}
func (p textValue) SetText(t string) error {
	return p.set(t)
}

type pointerValue struct {
	panicValue
}

func (p pointerValue) GetObject() ident.Id {
	return p.getId()
}

func (p pointerValue) SetObject(id ident.Id) (err error) {
	if p.GetObject() != id {
		if id.Empty() {
			err = p.set(id)
		} else {
			if e := p.mdl.match(id, p.prop.Relates); e != nil {
				err = e
			} else if e := p.set(id); e != nil {
				err = e
			}
		}
	}
	return
}

func (mdl Metal) match(id ident.Id, relates ident.Id) (err error) {
	if target, ok := mdl.GetInstance(id); !ok {
		err = fmt.Errorf("no such instance '%s'", id)
	} else if ok := mdl.AreCompatible(target.GetParentClass(), relates); !ok {
		err = fmt.Errorf("%s not compatible with %v", target, relates)
	}
	return
}
