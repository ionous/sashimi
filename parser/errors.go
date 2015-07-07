package parser

import (
	"fmt"
	"github.com/ionous/sashimi/util/errutil"
)

func UnknownInput(input string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("That's not something I recognize. You said, '%s'.", input)
	})
}

func DuplicateComprehension(name string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("The comprehension '%s' already exists.", name)
	})
}

func InvalidComprehension(name string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("The comprehension name '%s' isn't a valid.", name)
	})
}

func MismatchedNouns(name string, expected, received int) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("%s expected %d nouns, but matched %d nouns.", name, expected, received)
	})
}
