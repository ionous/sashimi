package runtime

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/runtime/api"
)

type propList struct {
	oa     ObjectAdapter
	prop   api.Property
	values api.Values
}

func (n propList) Len() (ret int) {
	return n.values.NumValue()
}

func (n propList) Get(i int) (ret G.IValue) {
	if cnt := n.values.NumValue(); i < 0 || i >= cnt {
		ret = NewNullValue(n.oa, fmt.Sprintf("(%s.%s).Get(%d) is out of range.", n.oa, n.prop, i))
	} else {
		ret = propValue{n.oa, n.prop, n.values.ValueNum(i)}
	}
	return
}

func (n propList) Contains(interface{}) bool {
	panic("not implemented")
}
