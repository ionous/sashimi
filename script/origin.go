package script

import (
	"fmt"
	S "github.com/ionous/sashimi/source"
	"runtime" // go runtime, for file line info
)

// location of that the fragment was declared
type Origin struct {
	pc uintptr
}

func NewOrigin(skip int) Origin {
	pc := []uintptr{0}
	// 0 is callers itself, 1 is this code
	runtime.Callers(skip+1, pc)
	return Origin{pc[0]}
}

// FIX: would be nice to change code from a string to an interface
// so that we can delay the expansion
func (this Origin) Code() S.Code {
	return S.Code(this.String())
}

func (this Origin) String() (str string) {
	f := runtime.FuncForPC(this.pc - 1)
	if f != nil {
		file, line := f.FileLine(this.pc - 1)
		str = fmt.Sprintf("%s:%d", file, line)
	}
	return str
}
