package runtime

import (
	G "github.com/ionous/sashimi/game"
)

type nullValue struct {
	prop string
}

func (_ nullValue) Num() (ret float32)  { return }
func (_ nullValue) SetNum(float32)      {}
func (n nullValue) Object() G.IObject   { return NullObject(n.prop) }
func (_ nullValue) SetObject(G.IObject) {}
func (_ nullValue) Text() (ret string)  { return }
func (_ nullValue) SetText(string)      {}

func NewNullValue(oa ObjectAdapter, str string) nullValue {
	oa.log(str)
	return nullValue{str}
}
