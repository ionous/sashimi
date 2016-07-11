package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard/live"
)

func init() {
	AddScript(func(s *Script) {
		s.The("objects",
			Can("print description").And("describing").RequiresNothing(),
			To("print description", func(g G.Play) {
				g.Go(Describe("object"))
			}))

		// FIX: When() puts the contents after the object
		// look at some default actions of the DOM
		// maybe it'd be better to put the print, not in the action,
		// but in a target handler: then this could be after by being in the capture.

		// FIX: After() isnt working well, it goes into the default action
		// but not all objects are containers, so it errors
	})
}
