package sbuf

import (
	"strings"
)

type StringBuffer struct {
	parts []Stringer
}

type Stringer interface {
	String() string
}

// S adds a string
func (sbuf *StringBuffer) B(b bool) *StringBuffer {
	sbuf.parts = append(sbuf.parts, Bool{b})
	return sbuf
}

// S adds a string
func (sbuf *StringBuffer) S(s string) *StringBuffer {
	sbuf.parts = append(sbuf.parts, String{s})
	return sbuf
}

// V adds the default string representation via Stringer
func (sbuf *StringBuffer) V(v interface{}) *StringBuffer {
	sbuf.parts = append(sbuf.parts, Str{v})
	return sbuf
}

// R adds a single rune
func (sbuf *StringBuffer) R(r rune) *StringBuffer {
	sbuf.parts = append(sbuf.parts, Rune{r})
	return sbuf
}

// D adds a decimal
func (sbuf *StringBuffer) D(i int64) *StringBuffer {
	sbuf.parts = append(sbuf.parts, Int64{i})
	return sbuf
}

// E adds an error
func (sbuf *StringBuffer) E(e error) *StringBuffer {
	sbuf.parts = append(sbuf.parts, Error{e})
	return sbuf
}

func (sbuf *StringBuffer) I(i int) *StringBuffer {
	sbuf.parts = append(sbuf.parts, Int{i})
	return sbuf
}

func (sbuf *StringBuffer) String() string {
	return sbuf.Join("")
}

func (sbuf *StringBuffer) Join(sep string) string {
	// FIX? speed isnt really important, but should probably use buffer directly
	strs := make([]string, len(sbuf.parts))
	for i, p := range sbuf.parts {
		strs[i] = p.String()
	}
	return strings.Join(strs, sep)
}
