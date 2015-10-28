package model

import (
	"fmt"
)

type Callback struct {
	File      string
	Line      int
	Iteration int
}

func (m Callback) String() string {
	return fmt.Sprintf("%s:%d#%d", m.File, m.Line, m.Iteration)
}
