package model

import (
	"github.com/ionous/mars/rt"
	"github.com/ionous/sashimi/util/sbuf"
)

type CallbackModel struct {
	File      string
	Line      int
	Iteration int
	Executes  rt.Execute
}

func (m CallbackModel) String() (ret string) {
	b := sbuf.New(m.File, ":", m.Line)
	if m.Iteration > 0 {
		b = b.Append(m.Iteration)
	}
	return b.String()
}
