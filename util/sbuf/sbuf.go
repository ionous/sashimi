package sbuf

import (
	"bytes"
	"strconv"
)

type StringBuffer struct {
	// alt could be an array of object and format
	// dont print till actually needed
	// as long as things are ideompotent all is well.
	b bytes.Buffer
}
type Stringer interface {
	String() string
}

func New() *StringBuffer {
	return &StringBuffer{}
}

func NewString(s string) *StringBuffer {
	return (&StringBuffer{}).S(s)
}

// S adds a string
func (sbuf *StringBuffer) S(s string) *StringBuffer {
	sbuf.b.WriteString(s)
	return sbuf
}

// V adds the default string representation via Stringer
func (sbuf *StringBuffer) V(v Stringer) *StringBuffer {
	return sbuf.S(v.String())
}

// R adds a single rune
func (sbuf *StringBuffer) R(r rune) *StringBuffer {
	sbuf.b.WriteRune(r)
	return sbuf
}

// D adds a decimal
func (sbuf *StringBuffer) D(i int64) *StringBuffer {
	return sbuf.S(strconv.FormatInt(i, 10))
}

// E adds an error
func (sbuf *StringBuffer) E(e error) *StringBuffer {
	return sbuf.S(e.Error())
}

func (sbuf *StringBuffer) Itoa(i int) *StringBuffer {
	return sbuf.S(strconv.FormatInt(int64(i), 10))
}

func (sbuf *StringBuffer) String() string {
	return sbuf.b.String()
}

func (sbuf *StringBuffer) Error() error {
	return errorBuf(sbuf.b)
}

type errorBuf bytes.Buffer

func (ebuf errorBuf) Error() string {
	b := bytes.Buffer(ebuf)
	return b.String()
}

// sbuf.NewString("bad nouns for").S(action).R(',').S(event).R(':').V(end).V(classes).Error()

//fmt.Errorf("bad nouns for %s,%s: %d, %s?", action, event, end, classes)
