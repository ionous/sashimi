package console

import (
	"fmt"
	"strings"
)

type BufferedOutput struct {
	accum []string
	line  []string
}

//
// Accumulate the passed args as text for Results().
//
func (this *BufferedOutput) Print(args ...interface{}) {
	str := fmt.Sprint(args...)
	this.line = append(this.line, str)
}

//
// Accumulate the passed args as text for Results().
//
func (this *BufferedOutput) Println(args ...interface{}) {
	this.Print(args...)
	this.flush()
}

//
func (this *BufferedOutput) flush() {
	line := strings.Join(this.line, " ")
	this.accum = append(this.accum, line)
	this.line = nil
}

//
// Returns an array of all printed text; resets the buffer.
//
func (this *BufferedOutput) Flush() (lines []string) {
	if len(this.line) > 0 {
		this.flush()
	}
	if a := this.accum; a != nil {
		lines = a
	} else {
		lines = []string{}
	}
	this.accum = nil
	return lines
}
