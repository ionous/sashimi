package parser

import "github.com/ionous/sashimi/util/ident"

// Parser collects input matchers.
type Parser map[ident.Id]*Comprehension

func NewParser() Parser {
	return make(map[ident.Id]*Comprehension)
}

// NewComprehension adds a pattern set.
// Name must be unique ( used to help with error-handling and auto-documentation. )
func (p Parser) NewComprehension(id ident.Id, matcher NewMatcher) (
	ret *Comprehension,
	err error,
) {
	if id.Empty() {
		err = InvalidComprehension(id)
	} else if _, exists := p[id]; exists {
		err = DuplicateComprehension(id)
	} else {
		comp := &Comprehension{id, matcher, nil}
		p[id] = comp
		ret = comp
	}

	return ret, err
}

// Parse the input, and generate a matching command.
// Returns the command found regardless of error.
func (p Parser) ParseInputString(input string) (ret Matched, err error) {
	matched := false
	for _, c := range p {
		if pattern, matcher, e := c.TryParse(input); e == nil || matcher != nil {
			ret, err = Matched{pattern, func() error { return matcher.OnMatch() }}, e
			matched = true
			break
		}
	}
	if !matched {
		err = UnknownInput(input)
	}
	return ret, err
}
