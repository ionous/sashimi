package standard

import (
	. "github.com/ionous/sashimi/script"
)

//
func init() {
	AddScript(func(s *Script) {
		// FIX: the player should really be a global variable; not an actor instance.
		// ( or, possibly a game object type of which there is one, with a relation of an actor. )
		s.The("actor",
			Called("player"),
			Is("scenery"),
		)
	})
}
