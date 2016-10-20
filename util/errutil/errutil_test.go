package errutil

import (
	"errors"
	//"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPrefix(t *testing.T) {
	err := errors.New("error")
	err = Prefix(err, "prefix")
	assert.EqualError(t, err, "prefix: error")
}

func TestAppend(t *testing.T) {
	one, two := errors.New("1"), errors.New("2")
	err := Append(one, two)
	list := strings.Split(err.Error(), "\n")
	assert.EqualValues(t, []string{"1", "2"}, list)
}

func TestErrorFunc(t *testing.T) {
	err := Func(func() string { return "fun" })
	assert.EqualError(t, err, "fun")
}

// type Printer struct {
// 	prefix string
// }

// func (p Printer) Errorf(format string, a ...interface{}) error {
// 	err := fmt.Errorf(format, a...)
// 	return Prefix(err, p.prefix)
// }

// func TestErrorf(t *testing.T) {
// 	p := Printer{"test"}
// 	var errorf Errorf = p
// 	s := errorf.Errorf("hello %s", "there")
// 	assert.EqualError(t, s, "test: hello there")
// }

type Stringed struct {
	s string
}

func (s Stringed) String() string {
	return s.s
}

func TestJoin(t *testing.T) {
	joined := New("a", Stringed{"b"}, "c")
	assert.EqualError(t, joined, "a b c")
}
