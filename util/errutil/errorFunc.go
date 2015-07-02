package errutil

//
// ErrorFunc implements the go "error" interface for parameterless functions, typically closures.
//
type Func func() string

//
// Call the function and return a string as per the go "error" interface.
//
func (this Func) Error() string {
	return this()
}
