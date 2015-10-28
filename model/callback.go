package model

import (
	"fmt"
)

type Callback struct {
	File      string
	Line      int
	Iteration int
}

func (m Callback) String() (ret string) {
	if m.Iteration > 0 {
		ret = fmt.Sprintf("%s:%d#%d", m.File, m.Line, m.Iteration)
	} else {
		ret = fmt.Sprintf("%s:%d", m.File, m.Line)
	}
	return ret
}
