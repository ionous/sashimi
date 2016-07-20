package errutil

// ErrorFunc implements the go Error() interface for parameterless functions, typically closures.
type Func func() string

// Call the function and return a string as per the go Error() interface.
func (e Func) Error() string {
	return e()
}
