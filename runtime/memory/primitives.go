package memory

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
)

type numValue struct{ panicValue }

func (p numValue) SetNum(f float32) error {
	return p.set(f)
}
func (p numValue) GetNum() float32 {
	return p.get().(float32)
}

type numValues struct{ panicValue }

func (p numValues) SetNum(f float32) error {
	return p.set(f)
}
func (p numValues) GetNum() float32 {
	return p.get().(float32)
}

type textValue struct{ panicValue }

func (p textValue) GetText() string {
	return p.get().(string)
}
func (p textValue) SetText(t string) error {
	return p.set(t)
}

type enumValue struct{ panicValue }

func (p enumValue) GetState() (ret ident.Id) {
	v := p.get()
	if idx, ok := v.(int); !ok {
		panic(fmt.Sprintf("internal error, couldnt convert state to int '%s.%s' %v(%T)", p.src, p.GetId(), v, v))
	} else {
		enum := p.prop.(M.EnumProperty)
		c, _ := enum.IndexToChoice(idx)
		ret = c
	}
	return
}
func (p enumValue) SetState(c ident.Id) (err error) {
	enum := p.prop.(M.EnumProperty)
	if idx, e := enum.ChoiceToIndex(c); e != nil {
		err = e
	} else {
		p.set(idx)
	}
	return
}

type pointerValue struct {
	panicValue
}

func (p pointerValue) GetObject() ident.Id {
	return p.get().(ident.Id)
}
func (p pointerValue) SetObject(n ident.Id) error {
	return p.set(n)
}
