package errutil

import (
	"github.com/ionous/sashimi/util/sbuf"
)

// Join expects each part is suppsrted by sbuf.Swtich
func New(parts ...interface{}) joined {
	return joined(parts)
}

type joined sbuf.Switch

func (j joined) Error() string {
	return sbuf.Switch(j).Join(" ")
}
