package parser

import (
	"regexp"
)

//
// Pattern holds the results of Comprehension.LearnPattern().
//
type Pattern struct {
	c *Comprehension
	*regexp.Regexp
	pattern    string
	wordToNoun []int // noun index => regexp match index
}

//
// Comprehension which contains this pattern.
//
func (p *Pattern) Comprehension() *Comprehension {
	return p.c
}

//
// Pattern string used to generate this pattern.
//
func (p *Pattern) Pattern() string {
	return p.pattern
}

//
// String returns pattern.
//
func (p *Pattern) String() string {
	return p.Pattern()
}

// TryPattern returns a Match helper if the passed input matches this pattern.
func (p *Pattern) TryPattern(input string) (ret Match, okay bool) {
	if words := p.FindStringSubmatch(input); len(words) > 0 {
		ret, okay = Match{p, words, p.wordToNoun}, true
	}
	return ret, okay
}
