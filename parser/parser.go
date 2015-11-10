package parser

// P collects input matchers.
type P struct {
	IMakeMatchers
	Comprehensions
}

func NewParser(m IMakeMatchers) P {
	return P{m, make(Comprehensions)}
}

// ParseInput to generate a matching command.
// Returns the command found regardless of error.
func (parser P) ParseInput(input string) (p *Pattern, m IMatch, err error) {
	matched := false
	for _, c := range parser.Comprehensions {
		if pattern, matcher, e := c.TryParse(input, parser.IMakeMatchers); e == nil || matcher != nil {
			p, m, err = pattern, matcher, e
			matched = true
			break
		}
	}
	if !matched {
		err = UnknownInput(input)
	}
	return
}
