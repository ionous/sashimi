package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard/live"
)

func init() {
	AddScript(func(s *Script) {
		s.The("actors",
			Can("smell").And("smelling").RequiresNothing(),
			To("smell", func(g G.Play) { ReflectToLocation(g, "report smell") }),

			Can("smell it").And("smelling it").RequiresOne("kind"),
			To("smell it", func(g G.Play) { ReflectToTarget(g, "report smell") }),
		)

		// kinds, to allow rooms and objects
		s.The("kinds",
			Can("report smell").And("reporting smell").RequiresOne("actor"),
			To("report smell", func(g G.Play) {
				actor := g.The("action.Target")
				if g.The("player") == actor {
					g.Say("You smell nothing unexpected.")
				} else {
					g.Say(actor.Text("Name"), "sniffs.")
				}
			}),
		)

		s.Execute("smell it", Matching("smell|sniff {{something}}"))
		s.Execute("smell", Matching("smell|sniff"))
	})
}
