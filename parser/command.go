package parser

//
// NewMatcher gets called back every new parser matching cycle.
//
type NewMatcher func() (IMatch, error)

//
// IMatch provides an algorithm which the Parser uses to match and execute commands.
// MatchNoun gets called successively for each word in a user's input.
// Execute gets called when input has been exhausted.
//
type IMatch interface {
	// MatchNoun transforms a word into one specific noun.
	// ex. "glasses", "the" => "horn-rimmed kryptonian disguise device"
	MatchNoun(word string, article string) error
	// Matched gets called after all nouns in an input have been parsed succesfully.
	OnMatch() error
}

type Matched struct {
	Pattern *Pattern
	OnMatch func() error
}
