package model

//
// a script variable variant
type IValue interface {
	// return the class/property this value was derived from
	Property() IProperty
	// return the underlying value of this variant, true if every set.
	Any() (interface{}, bool)
	// set the value of the variant
	SetAny(v interface{}) error
	// return the value of this variant as a string
	String() string
}
