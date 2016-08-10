package errutil

// Func delays error string generation until the error string actually get displayed. It does this by implementing the go Error() interface for parameterless functions.
// For example: funError:= errutil.Func(func() string { return "fun" })
type Func func() string

// Error implements go's Error() interface for a function.
func (e Func) Error() string {
	return e()
}
