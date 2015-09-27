package console

import (
	"fmt"
	"strings"
)

// BufferedOutput acts similar to go's fmt.Print and fmt.Println;
// accumulating output until Flush()
type BufferedOutput struct {
	accum []string
	line  []string
}

// Print into the current (pending) line.
func (buf *BufferedOutput) Print(args ...interface{}) {
	if len(args) > 0 {
		str := fmt.Sprint(args...)
		buf.line = append(buf.line, str)
	}
}

// Accumulate the passed args as text for Results().
func (buf *BufferedOutput) Println(args ...interface{}) {
	if len(args) > 0 {
		buf.Print(args...)
	}
	buf.flush()
}

// flush generates any recent print statements into a single line.
func (buf *BufferedOutput) flush() {
	if len(buf.line) > 0 {
		line := strings.Join(buf.line, " ")
		buf.accum = append(buf.accum, line)
	}
	buf.line = nil
}

//
// Returns an array of all printed text; resets the buffer.
//
func (buf *BufferedOutput) Flush() (lines []string) {
	if len(buf.line) > 0 {
		buf.flush()
	}
	if a := buf.accum; len(a) != 0 {
		lines = a
	} else {
		lines = []string{}
	}
	buf.accum = nil
	return lines
}
