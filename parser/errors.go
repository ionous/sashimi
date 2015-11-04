package parser

import (
	"fmt"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

func UnknownInput(input string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("That's not something I recognize. You said, '%s'.", input)
	})
}

func DuplicateComprehension(id ident.Id) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("The comprehension '%s' already exists.", id)
	})
}

func InvalidComprehension(id ident.Id) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("The comprehension '%s' isn't valid.", id)
	})
}

func MismatchedNouns(id ident.Id, expected, received int) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("%s expected %d nouns, but matched %d nouns.", id, expected, received)
	})
}
