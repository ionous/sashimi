package parser

import "github.com/ionous/sashimi/util/lang"

// Match helps input parsing, returned by TryParse().
type Match struct {
	pattern    *Pattern
	words      []string // from regexp.FindSubStringMatch
	wordToNoun []int    // noun index => regexp match index
}

// Pattern which matched.
func (m Match) Pattern() *Pattern {
	return m.pattern
}

// MatchNouns changes words into nouns.
// Returns error if any of the nouns could not be parsed.
func (m Match) MatchNouns(matcher IMatch) (err error) {
	for _, matchIdx := range m.wordToNoun {
		word := m.words[matchIdx]
		article, name := lang.SliceArticle(word)
		if e := matcher.MatchNoun(name, article); e != nil {
			err = e
			break
		}
	}
	return err
}
