package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func init() {
	AddScript(func(s *Script) {
		s.The("actors",
			Can("listen").And("listening").RequiresNothing(),
			To("listen", ReflectToLocation("report listen")),

			Can("listen to").And("listening to").RequiresOne("kind"),
			To("listen to", ReflectToTarget("report listen")),
		)
		// kinds, to allow rooms and objects
		s.The("kinds",
			Can("report listen").And("reporting listen").RequiresOne("actor"),
			To("report listen", func(g G.Play) {
				actor := g.The("action.Target")
				if g.The("player") == actor {
					g.Say("You hear nothing unexpected.")
				} else {
					g.Say(actor.Text("Name"), "listens.")
				}
			}))
		s.Execute("listen to", Matching("listen to {{something}}").Or("listen {{something}}"))
		s.Execute("listen", Matching("listen"))
	})
}
