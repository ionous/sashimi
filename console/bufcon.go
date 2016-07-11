package console

import "github.com/ionous/sashimi/util"

//
// Creates a console whose input comes from the passed strings.
// when the strings are exhausted the Readln() returns false.
//
func NewBufCon(strs []string) *BufCon {
	return &BufCon{strs: strs}
}

//
type BufCon struct {
	strs  []string
	index int
	util.BufferedOutput
}

//
// Returns the next input string, false when input has been exhausted.
//
func (this *BufCon) Readln() (ret string, okay bool) {
	okay = this.index < len(this.strs)
	if okay {
		ret = this.strs[this.index]
		this.index++
	}
	return ret, okay
}
