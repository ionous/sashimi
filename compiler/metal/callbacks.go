package metal

import (
	"github.com/ionous/sashimi/util/ident"
)

type CallbackList struct {
	callbacks []ident.Id
}

func (cl CallbackList) NumCallback() int {
	return len(cl.callbacks)
}

func (cl CallbackList) CallbackNum(i int) ident.Id {
	p := cl.callbacks[i]
	return p // CallbackWrapper(p)
}
