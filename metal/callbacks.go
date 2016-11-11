package metal

import (
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
)

type CallbackList struct {
	callbacks []M.CallbackModel
}

func (cl CallbackList) NumCallback() int {
	return len(cl.callbacks)
}

func (cl CallbackList) CallbackNum(i int) meta.Callback {
	return cl.callbacks[i].Calls
}
