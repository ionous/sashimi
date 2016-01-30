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

	Get(string) IValue
	List(string) IList

	// other built ins
	Go(action string, withTargetAndContext ...IObject) IPromise

	// FIX: this should probably just be an action.
	// Go("say", ...)
	Says(string)

	// old: property access
	Num(string) float32
	SetNum(string, float32)

	Object(string) IObject
	ObjectList(string) []IObject
	Set(string, IObject)

	Text(string) string
	SetText(string, string)
}

type IValue interface {
	Num() float32
	SetNum(float32)

	Object() IObject
	SetObject(IObject)

	Text() string
	SetText(string)

	State() ident.Id
	SetState(ident.Id)
}

type IList interface {
	Len() int
	Get(int) IValue
	Pop() IValue

	// IValue, IObject, numbers, text
	Contains(interface{}) bool

	AppendNum(float32)
	AppendText(string)
	AppendObject(IObject)

	Reset()
}
