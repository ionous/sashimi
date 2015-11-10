package parser

import "github.com/ionous/sashimi/util/ident"

// Comprehension contains a set of patterns against which user input gets matched.
// Comprehensions are created via NewComprehension() and expanded via LearnPattern().
type Comprehension struct {
	id       ident.Id
	patterns []Pattern
}

// Name of this pattern set.
func (c *Comprehension) Id() ident.Id {
	return c.id
}

// LearnPattern adds a new pattern to this pattern set.
// Patterns are tried successively to match a user's input.
func (c *Comprehension) LearnPattern(pattern string) (p Pattern, err error) {
	// split the pattern into groups separated by tags
	groups, tags := tokenize(pattern)
	if nouns, e := newNounCheck(groups, tags); e != nil {
		err = e
	} else {
		p = Pattern{c, nouns.exp, pattern, nouns.matchIndices()}
		c.patterns = append(c.patterns, p)
	}
	return p, err
}

// Patterns used by this pattern set.
// Note, it's not a copy because pattern objects are not modifiable by the caller.
func (c *Comprehension) Patterns() []Pattern {
	return c.patterns
}

// Parse the paseed input, trying each of the known patterns in turn.
// If one of the patterns matched, found will contain that pattern --
// even if the nouns from that pattern failed to match.
func (c *Comprehension) TryParse(input string, matchFactory IMakeMatchers) (found *Pattern, matched IMatch, err error) {
	// for all patterns in those sets:
	for _, p := range c.patterns {
		// try the pattern:
		if match, ok := p.TryPattern(input); ok {
			// let the caller know that something matched:
			found = &p
			// try the nouns:
			if imatch, e := matchFactory.NewMatcher(c.id); e != nil {
				err = e // provisional return
			} else if e := match.MatchNouns(imatch); e != nil {
				err = e // provisional return
			} else {
				matched, err = imatch, nil
				break
			}
		}
	}
	if found == nil {
		err = UnknownInput(input)
	}
	return found, matched, err
}
