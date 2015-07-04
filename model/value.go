package model

//
// a script variable variant
//
type IValue interface {
	// return the underlying value of this variant
	Any() interface{}
	// return a (formatted) string representation of the value.
	String() string
}
