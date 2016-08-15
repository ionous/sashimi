package game

import (
	"github.com/ionous/sashimi/util/ident"
)

// IObject might have a few implementations:
// * valid objects
// * invalid objects while debugging which panics or errors
// * invalid objects during play which eat all errors
// * testing which logs all calls
type IObject interface {
	Id() ident.Id
	Exists() bool // FIX: added for obj.Object() tests, alternatives?
	// FromClass returns true if the object was derived from the passed plural named class.
	// FIX: seems to programmery, alternatives?
	FromClass(string) bool
	// Parent returns the spatial parent, enclosure, of an object.
	// by default, there is no such object, the standard rules, etc. define their own thing.
	ParentRelation() (IObject, string)

	Is(string) bool
	IsNow(string)

	// Get returns the named property.
	Get(string) IValue

	// List returns the named list.
	List(string) IList

	// Go run an action defined on this object.
	Go(action string, withTargetAndContext ...IObject) IPromise

	// FIX: this should probably just be an action.
	// Go("say", ...)
	Says(string)

	// deprecated: prefer Get()
	Num(string) float64
	SetNum(string, float64)

	Object(string) IObject
	ObjectList(string) []IObject
	Set(string, IObject)

	Text(string) string
	SetText(string, string)
}

// IValue provides access to an object property
type IValue interface {
	Num() float64
	SetNum(float64)

	Object() IObject
	SetObject(IObject)

	Text() string
	SetText(string)

	State() ident.Id
	SetState(ident.Id)
}

// IQuery provides access to an object list
type IQuery interface {
	HasNext() bool
	Next() IObject
}

// IValue provides access to an object list
type IList interface {
	// Len returns the number of values in the list.
	Len() int
	// Get returns the value of the nth values in the list.
	Get(int) IValue
	// Contains returns true if the passed parameter exists in the list.
	// The parameter can be an IValue, IObject, number, text.
	Contains(interface{}) bool

	// AppendNum adds the passed number to the end of the list.
	AppendNum(float64)
	// AppendText adds the passed string to the end of the list.
	AppendText(string)
	// AppendObject adds the passed object to the end of the list.
	AppendObject(IObject)

	// Pop removes the last value from the list.
	Pop() IValue
	// Reset the list by completely emptying it.
	Reset()
}
