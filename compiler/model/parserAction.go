package model

import "github.com/ionous/sashimi/util/ident"

// ParserAction commands that beome an action
type ParserAction struct {
	Action   ident.Id
	Commands []string
}
