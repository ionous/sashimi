package parser

//import "fmt"

//
// Transform a series of words into a series of nouns,
// implemented by Parser clients.
//
type IMatch interface {
	// Transform a word into one specific noun.
	// ex. "glasses", "the" => "horn-rimmed kryptonian disguise device"
	MatchNoun(word string, article string) (string, error)
}

//
// An action triggered by user input,
// implemented by Parser clients.
//
type ICommand interface {
	// Helper to parse input into nouns.
	// The number, and kind, of nouns depends on the command in question.
	NewMatcher() IMatch

	// Run some function for a set of matched nouns.
	RunCommand(...string) error
}

//
// Result of Parser.AddCommand(),
// implemented by the parser itself.
//
type ILearn interface {
	// ex. "examine {{something}}"
	LearnPattern(string) error
}

//
// Result of a successful Parser.Parse()
//
type CommandMatch struct {
	Command ICommand
	Nouns   []string
	Pattern string
}

func (this CommandMatch) Run() error {
	return this.Command.RunCommand(this.Nouns...)
}
