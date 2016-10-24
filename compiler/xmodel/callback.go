package xmodel

import (
	"github.com/ionous/sashimi/util/sbuf"
)

type Callback struct {
	File      string
	Line      int
	Iteration int
}

func (m Callback) String() (ret string) {
	buf := sbuf.New(m.File, ":", m.Line)
	if m.Iteration > 0 {
		buf = buf.Append("#").Append(m.Iteration)
	}
	return buf.String()
}
