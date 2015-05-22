package parser

import (
	"fmt"
)

//
// A collection of input matchers
//
type Parser struct {
	comprehension map[string]*Comprehension
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
// Add a new command returns an interface to which new patterns can be added.
// Name must be unique ( used to help with error-handling and auto-documentation. )
//
func (this *Parser) AddCommand(name string, cmd ICommand, nounCount int) (
	ret ILearn,
	err error,
) {
	if name == "" {
		err = fmt.Errorf("add commmand expects a valid name")
	} else if _, exists := this.comprehension[name]; exists {
		err = fmt.Errorf("command %s already exists", name)
	} else if nounCount < 0 || nounCount > 2 {
		err = fmt.Errorf("noun count %d out of range [0-2]", nounCount)
	} else {
		comp := &Comprehension{this, name, cmd, nounCount, nil}
		this.comprehension[name] = comp
		ret = comp
	}
	return ret, err
}

//
// Read the input, and generate a matching command.
// found is set if the command was found but not matched
// err is only nil if the command was matched
//
func (this *Parser) Parse(input string) (found string, ret CommandMatch, err error) {
	// default return, cleared on success:
	err = UnknownCommand{}
	// for all registered commands:
	for _, c := range this.comprehension {
		// for all patterns in those commands:
		for _, p := range c.patterns {
			// try the pattern:
			if match, ok := p.TryParse(input); ok {
				// let the caller know that something matched:
				found = c.name
				// try the nouns:
				if nouns, e := match.MatchNouns(c.command.NewMatcher()); e == nil {
					ret, err = CommandMatch{c.command, nouns, p.pattern}, nil
					break
				} else {
					err = e // provisional error, but keep going.
					continue
				}
			}
		}
	}
	return found, ret, err
}

type UnknownCommand struct{}

func (UnknownCommand) Error() string {
	return "That's not something I recognize."
}
