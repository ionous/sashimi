package console

import (
	"fmt"
	"strings"
)

//
// creates a console consisting only of the passed strings
func NewBufCon(strs []string) *BufCon {
	return &BufCon{strs, 0, nil, nil}
}

type BufCon struct {
	strs  []string
	index int
	accum []string
	line  []string
}

//
func (this *BufCon) Print(args ...interface{}) {
	str := fmt.Sprint(args...)
	this.line = append(this.line, str)
}

//
func (this *BufCon) Println(args ...interface{}) {
	this.Print(args...)
	this.flush()
}

//
func (this *BufCon) flush() {
	line := strings.Join(this.line, " ")
	fmt.Println(line)
	this.accum = append(this.accum, line)
	this.line = nil
}

//
func (this *BufCon) Results() []string {
	if len(this.line) > 0 {
		this.flush()
	}
	return this.accum
}

//
func (this *BufCon) Readln() (ret string, okay bool) {
	okay = this.index < len(this.strs)
	if okay {
		ret = this.strs[this.index]
		this.index++
	}
	return ret, okay
}
