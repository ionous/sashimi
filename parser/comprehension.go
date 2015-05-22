package parser

import (
	"fmt"
	"regexp"
)

//
// created via Parser.New()
//
type Comprehension struct {
	parser    *Parser
	name      string
	command   ICommand
	nounCount int
	patterns  []Pattern
}

type Pattern struct {
	*regexp.Regexp
	pattern     string
	matchByNoun []int // noun index => regexp match index
}

type TryParseResult int

//
// implements ILearn
//
func (this *Comprehension) LearnPattern(pattern string) (err error) {
	// split the pattern into groups separated by tags
	groups, tags := tokenize(pattern)
	nouns := newNounCheck(this.nounCount)
	if exp, e := patternize(nouns, groups, tags); e != nil {
		err = e
	} else if !nouns.foundAllNouns() {
		// FIX: allow missing nouns to be specified via... callbacks?
		err = fmt.Errorf("`%s` pattern `%s` didn't specify enough nouns (%d of %d)", this.name, pattern, nouns.count, this.nounCount)
	} else {
		indices := nouns.matchIndices()
		p := Pattern{exp, pattern, indices[0:this.nounCount]}
		this.patterns = append(this.patterns, p)
	}
	return err
}

//
// If this pattern matches the passed input, return a new pattern matcher
//
func (this *Pattern) TryParse(input string) (ret PatternMatch, okay bool) {
	if match := this.FindStringSubmatch(input); len(match) > 0 {
		ret, okay = PatternMatch{match, this.matchByNoun}, true
	}
	return ret, okay
}

//
// Helper for input parsing, returned by TryMatch()
//
type PatternMatch struct {
	match       []string // from regexp.FindSubStringMatch
	matchByNoun []int    // noun index => regexp match index
}

//
// Change all the words into nouns.
// Returns error if any of the nouns could not be parsed
//
func (this *PatternMatch) MatchNouns(matcher IMatch) (nouns []string, err error) {
	for _, matchIdx := range this.matchByNoun {
		words := this.match[matchIdx]
		name, article := NormalizeNoun(words)
		if noun, e := matcher.MatchNoun(name, article); e == nil {
			nouns = append(nouns, noun)
		} else {
			err = e
			break
		}
	}
	return nouns, err
}
