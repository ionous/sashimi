package ident

import (
	"bytes"
	"github.com/ionous/sashimi/util/lang"
	"strings"
	"unicode"
)

// Id uniquely identifies some resource.
// NOTE: go is not capable of comparing slices(!?), and there is no interface for making a type hashable(!?)
// ALSO: the default json encoding --- even when the marshal is implemented!? -- doesnt seem to support a struct as a map key.
// Therefore, Id cannot store its parts as an []string. they must be joined first.
type Id string

func (id Id) Raw() string {
	return string(id)
}

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

// Reserved using an indicator to tag unique strings ids.
// FIX: i think with a judicious use of name vs. id -- this could be removed.
// ie. allow ids to be non-dashed, raw text. the dashing was in part to normalize multiple spaces and capitializations.
func (id Id) Reserved() bool {
	return len(id) > 0 && id[0] == '~'
}

func (id Id) Equals(other Id) bool {
	return other == id
}

// for some reason strings.Compare doesnt exist in go/appengine:
// theres this comment in the string source:
// NOTE(rsc): ... Basically no one should use strings.Compare.
func Compare(a, b Id) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return +1
}

func Join(a, b Id) Id {
	return Id(strings.Join(append(a.Split(), b.Split()...), "-"))
}

func Empty() (ret Id) {
	return
}

// MakeId creates a new string id from the passed raw string.
// Dashes and spaces are treated as word separators; sequences of numbers and sequences of letters are treated as separate words.
// NOTE: Articles ( the, etc. ) are stripped for easier matching at the script/table/level.
func MakeId(name string) (ret Id) {
	if len(name) > 0 {
		if name[0] == '~' {
			ret = Id(name)
		} else {
			type word int
			const (
				noword word = iota
				letter
				number
			)
			var parts parts
			inword, wasUpper := noword, false

			for _, r := range name {
				if r == '-' || r == '_' || r == '=' || unicode.IsSpace(r) {
					inword = noword
					continue
				}

				if unicode.IsDigit(r) {
					if sameWord := inword == number; !sameWord {
						parts.flush()
					}
					parts.WriteRune(r)
					wasUpper = false
					inword = number
				} else if unicode.IsLetter(r) {
					currUpper := unicode.IsUpper(r)
					// classify some common word changes
					sameWord := (inword == letter) && ((wasUpper == currUpper) || (wasUpper && !currUpper))
					if currUpper {
						r = unicode.ToLower(r)
					}
					if !sameWord {
						parts.flush()
					}
					parts.WriteRune(r) // docs say err is always nil
					wasUpper = currUpper
					inword = letter
				}
			}

			dashed := strings.Join(parts.flush(), "-")
			ret = Id(dashed)
		}
	}
	return
}

// Split the id into separated, lower cased components.
func (id Id) Split() (ret []string) {
	if !id.Empty() {
		dashed := string(id)
		if dashed[0] == '~' {
			ret = []string{dashed}
		} else {
			ret = strings.Split(dashed, "-")
		}
	}
	return
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
