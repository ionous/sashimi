package model

import (
	"bytes"
	"regexp"
	"unicode"
)

//
// FIX: pretty much everything that has a map of M.StringId also needs a find by name
//
type StringPair struct {
	id  StringId
	str string
}

//
// TitleCaseWord uniquely identifying some resource.
//
type StringId string

//
// Return the StringId as a regular string.
//
func (this StringId) String() string {
	return string(this)
}

//
// Create a new string id from the passed raw string.
// Dashes and spaces are treated as word separators.
// Other Illegal characters ( leading digits and non-word characters ) are stripped.
// Articles ( the, etc. ) are stripped for easier matching at the script/table/level.
//
func MakeStringId(name string) StringId {
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

	_, str := stripArticle(buffer.String())
	return StringId(str)
}

//
// Break the title cased string id into separated, lower cased components.
//
func (this StringId) Split() (ret []string) {
	p := parts{}
	for _, r := range this {
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

func (this *parts) flush() []string {
	if this.Len() > 0 {
		this.arr = append(this.arr, this.String())
		this.Reset()
	}
	return this.arr
}

//
// unfortunately, the parser has this too.
// FIX: move to some common space?
//
var articles = regexp.MustCompile(`^(The|A|An|Our|Some)[[:upper:]]`)

func stripArticle(name string) (article, bare string) {
	pair := articles.FindStringIndex(name)
	if pair == nil {
		bare = name
	} else {
		split := pair[1] - 1
		article = name[:split]
		bare = name[split:]
	}
	return article, bare
}
