package model

//
// a script variable variant
//
type IValue interface {
	// return the underlying value of this variant, true if ever set.
	Any() (interface{}, bool)
	// return a (formatted) string representation of the value.
	String() string
}
