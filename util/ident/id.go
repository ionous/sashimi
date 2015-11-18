package ident

import (
	"bytes"
	"github.com/ionous/sashimi/util/lang"
	"github.com/satori/go.uuid"
	"strings"
	"unicode"
)

// Id uniquely identifies some resource.
// NOTE: go is not capable of comparing slices(!?), and there is no interface for making a type hashable(!?)
// ALSO: the default json encoding --- even when the marshal is implemented!? -- doesnt seem to support a struct as a map key.
// Therefore, Id cannot store its parts as an []string. they must be joined first.
type Id string

// String representation of the id (ex. for fmt), currently TitleCase.
func (id Id) String() (ret string) {
	if id.Empty() {
		ret = "<empty>"
	} else {
		parts := id.Split()
		for i, s := range parts {
			parts[i] = lang.Capitalize(s)
		}
		ret = strings.Join(parts, "")
	}
	return ret
}

func (id Id) Empty() bool {
	return len(id) == 0
}

func Compare(a, b Id) int {
	return strings.Compare(string(a), string(b))
}

func Join(a, b Id) Id {
	return Id(strings.Join(append(a.Split(), b.Split()...), "-"))
}

func Empty() (ret Id) {
	return
}

func MakeUniqueId() Id {
	str := uuid.NewV4().String()
	return Id(str)
}

// MakeId creates a new string id from the passed raw string.
// Dashes and spaces are treated as word separators.
// Other Illegal characters ( leading digits and non-word characters ) are stripped.
// NOTE: Articles ( the, etc. ) are stripped for easier matching at the script/table/level.
func MakeId(name string) Id {
	var parts parts
	started, inword, wasUpper := false, false, false

	for _, r := range name {
		if r == '-' || r == '_' || r == '=' || unicode.IsSpace(r) {
			inword = false
			continue
		}

		// this test similar to go scanner's is identifier alg:
		// it has _, which we treat as a space, and
		// its i>0, but since we are stripping spaces, 'started' is better.
		if unicode.IsLetter(r) || (started && unicode.IsDigit(r)) {
			started = true
			nowUpper := unicode.IsUpper(r)
			// classify some common word changes
			sameWord := inword && ((wasUpper == nowUpper) || (wasUpper && !nowUpper))
			if nowUpper {
				r = unicode.ToLower(r)
			}
			if !sameWord {
				parts.flush()
			}
			parts.WriteRune(r) // docs say err is always nil
			wasUpper = nowUpper
			inword = true
		}
	}

	dashed := strings.Join(parts.flush(), "-")
	return Id(dashed)
}

// Split the id into separated, lower cased components.
func (id Id) Split() []string {
	dashed := string(id)
	return strings.Split(dashed, "-")
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
