package game

import (
	"github.com/ionous/sashimi/util/ident"
)

// There might be a few implementations:
// * valid objects
// * invalid objects while debugging which panics or errors
// * invalid objects during play which eat all errors
// * testing which logs all calls
type IObject interface {
	Id() ident.Id
	Exists() bool      // FIX: added for obj.Object() tests, alternatives?
	Class(string) bool // FIX: seems to programmery, alternatives?

	// property access
	Is(string) bool
	SetIs(string)

	Num(string) float32
	SetNum(string, float32)

	Object(string) IObject
	ObjectList(string) []IObject
	Set(string, IObject)

	Text(string) string
	SetText(string, string)

	// other built ins
	Go(action string, withTargetAndContext ...IObject)
	Says(string)
}
