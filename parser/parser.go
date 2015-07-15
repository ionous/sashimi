package parser

//
// A collection of input matchers
//
type Parser struct {
	comps map[string]*Comprehension
}

//
// Create a new parser.
//
func NewParser() *Parser {
	return &Parser{
		make(map[string]*Comprehension),
	}
}

//
// NewComprehension adds the named pattern set.
// Returns an object to which new
// Name must be unique ( used to help with error-handling and auto-documentation. )
//
func (p *Parser) NewComprehension(name string, matcher NewMatcher) (
	ret *Comprehension,
	err error,
) {
	if name == "" {
		err = InvalidComprehension(name)
	} else if _, exists := p.comps[name]; exists {
		err = DuplicateComprehension(name)
	} else {
		comp := &Comprehension{p, name, matcher, nil}
		p.comps[name] = comp
		ret = comp
	}
	return ret, err
}

//
// Parse the input, and generate a matching command.
// Returns the command found regardless of error.
//
func (p *Parser) ParseInput(input string) (ret Matched, err error) {
	matched := false
	for _, c := range p.comps {
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
