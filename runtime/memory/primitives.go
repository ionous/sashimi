package memory

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
)

type numValue struct{ panicValue }

func (p numValue) SetNum(f float32) error {
	if p.setValue == nil {
		p.panicValue.SetNum(f)
	}
	return p.setValue(p.prop, f)
}
func (p numValue) GetNum() float32 {
	return p.getValue(p.prop).(float32)
}

type textValue struct{ panicValue }

func (p textValue) GetText() string {
	return p.getValue(p.prop).(string)
}
func (p textValue) SetText(t string) error {
	if p.setValue == nil {
		p.panicValue.SetText(t)
	}
	return p.setValue(p.prop, t)
}

type enumValue struct{ panicValue }

func (p enumValue) GetState() (ret ident.Id) {
	v := p.getValue(p.prop)
	if idx, ok := v.(int); !ok {
		panic(fmt.Sprintf("internal error, couldnt convert state to int '%s.%s' %v(%T)", p.src, p.prop.GetId(), v, v))
	} else {
		enum := p.prop.(M.EnumProperty)
		c, _ := enum.IndexToChoice(idx)
		ret = c
	}
	return
}
func (p enumValue) SetState(c ident.Id) (err error) {
	if p.setValue == nil {
		p.panicValue.SetState(c)
	}
	enum := p.prop.(M.EnumProperty)
	if idx, e := enum.ChoiceToIndex(c); e != nil {
		err = e
	} else {
		p.setValue(p.prop, idx)
	}
	return
}

type pointerValue struct {
	panicValue
}

func (p pointerValue) GetObject() ident.Id {
	return p.getValue(p.prop).(ident.Id)
}
func (p pointerValue) SetObject(o ident.Id) error {
	if p.setValue == nil {
		p.panicValue.SetObject(o)
	}
	return p.setValue(p.prop, o)
}
