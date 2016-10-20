package parser

import (
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

func UnknownInput(input string) error {
	return errutil.New("That's not something I recognize. You said:", input)
}

func DuplicateComprehension(id ident.Id) error {
	return errutil.New("duplicate comprehension", id)
}

func InvalidComprehension(id ident.Id) error {
	return errutil.New("invalid comprehension", id)
}

func MismatchedNouns(id ident.Id, expected, received int) error {
	return errutil.New("mismatched nouns", id, "expected", expected, "received", received)
}
