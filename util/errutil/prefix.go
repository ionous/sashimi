package errutil

import "strings"

// Prefix the passed error with the passed string, returning a brand new error.
func Prefix(err error, s string) error {
	return Func(func() string {
		return strings.Join([]string{s, err.Error()}, ": ")
	})
}
