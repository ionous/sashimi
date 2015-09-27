package game

import (
	"github.com/ionous/sashimi/util/ident"
)

// IObject might have a few implementations:
// * valid objects
// * invalid objects while debugging which panics or errors
// * invalid objects during play which eat all errors
// * testing which logs all calls
// FIX: the interface would be easier to use, and easier to provide implementations of,
// if get and set returned variants -- then the variant values could handle the setting.
// (except: Is and IsNow are nice )
type IObject interface {
	Id() ident.Id
	Exists() bool      // FIX: added for obj.Object() tests, alternatives?
	Class(string) bool // FIX: seems to programmery, alternatives?
	// Remove a previously new'd data object.
	Remove()

	// property access
	Is(string) bool
	IsNow(string)

	Num(string) float32
	SetNum(string, float32)

	Object(string) IObject
	ObjectList(string) []IObject
	Set(string, IObject)

	Text(string) string
	SetText(string, string)

	// other built ins
	Go(action string, withTargetAndContext ...IObject)

	// FIX: this should probably just be an action.
	// Go("say", ...)
	Says(string)
}
