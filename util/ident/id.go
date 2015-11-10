package ident

import (
	"bytes"
	"github.com/satori/go.uuid"
	"unicode"
)

//
// TitleCaseWord uniquely identifying some resource.
//
type Id string

//
// Return the Id as a regular string.
//
func (id Id) String() (ret string) {
	if id.Empty() {
		ret = "<empty>"
	} else {
		ret = string(id)
	}
	return ret
}

func (id Id) Empty() bool {
	return id == ""
}

func Empty() Id {
	return ""
}

func MakeUniqueId() Id {
	return Id(uuid.NewV4().String())
}

//
// Create a new string id from the passed raw string.
// Dashes and spaces are treated as word separators.
// Other Illegal characters ( leading digits and non-word characters ) are stripped.
// Articles ( the, etc. ) are stripped for easier matching at the script/table/level.
//
func MakeId(name string) Id {
	var buffer bytes.Buffer
	started, inword, wasUpper := false, false, false

	for _, r := range name {
		if r == '-' || r == '_' || r == '=' || unicode.IsSpace(r) {
			inword = false
			continue
		}

		// this test similar to go scanner's is identifier alg:
		// it has _ which we treat as a space, and
		// it tested i>0, but since we are stripping spaces, 'started' is better.
		if unicode.IsLetter(r) || (started && unicode.IsDigit(r)) {
			started = true
			nowUpper := unicode.IsUpper(r)
			// classify some common word changes
			sameWord := inword && ((wasUpper == nowUpper) || (wasUpper && !nowUpper))
			if !sameWord {
				r = unicode.ToUpper(r)
			} else {
				r = unicode.ToLower(r)
			}
			buffer.WriteRune(r) // docs say err is always nil
			wasUpper = nowUpper
			inword = true
		}
	}

	return Id(buffer.String())
}

//
// Break the title cased string id into separated, lower cased components.
//
func (id Id) Split() (ret []string) {
	p := parts{}
	for _, r := range id {
		if !unicode.IsUpper(r) {
			p.WriteRune(r)
		} else {
			p.flush()
			p.WriteRune(unicode.ToLower(r))
		}
	}
	return p.flush()
}

type parts struct {
	bytes.Buffer
	arr []string
}

func (p *parts) flush() []string {
	if p.Len() > 0 {
		p.arr = append(p.arr, p.String())
		p.Reset()
	}
	return p.arr
}
