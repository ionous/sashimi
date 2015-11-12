package runtime

import (
	G "github.com/ionous/sashimi/game"
)

type nullList struct {
	prop string
}

func (_ nullList) Len() int {
	return 0
}

func (n nullList) Get(int) G.IValue {
	return nullValue{n.prop}
}

func (n nullList) Contains(interface{}) bool {
	return false
}

func NewNullList(oa ObjectAdapter, str string) nullList {
	oa.log(str)
	return nullList{str}
}
