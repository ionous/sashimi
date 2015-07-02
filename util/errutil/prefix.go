package errutil

import "fmt"

// Prefix the passed error with the passed string, returning a brand new error.
func Prefix(err error, s string) error {
	return Func(func() string {
		return fmt.Sprintf("%s: %s", s, err)
	})
}
