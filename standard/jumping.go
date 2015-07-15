package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func init() {
	AddScript(func(s *Script) {
		s.The("actors",
			Can("jump").And("jumping").RequiresNothing(),
			To("jump", ReflectToLocation("report jump")),
		)

		// kinds, to allow rooms and objects
		s.The("kinds",
			Can("report jump").And("reporting jump").RequiresOne("actor"),
			To("report jump", func(g G.Play) {
				actor := g.The("action.Target")
				// FIX? inform often, but not always, tests for trying silently,
				// "if the action is not silent" ...
				// seems... strange. why report if if its silent?
				if g.The("player") == actor {
					g.Say("You jump on the spot.")
				} else {
					g.Say(actor.Text("Name"), "jumps on the spot.")
				}
			}))

		s.Execute("jump", Matching("jump|skip|hop"))
	})
}
