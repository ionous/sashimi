package parser

import "github.com/ionous/sashimi/util/ident"

// NewMatcher gets called back every new parser matching cycle.
type IMakeMatchers interface {
	NewMatcher(ident.Id) (IMatch, error)
}

// IMatch provides an algorithm which the Parser uses to match and execute commands.
// MatchNoun gets called successively for each word in a user's input.
type IMatch interface {
	// MatchNoun transforms a word into one specific noun.
	// ex. "glasses", "the" => "horn-rimmed kryptonian disguise device"
	MatchNoun(word string, article string) error
}
