package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func init() {
	AddScript(func(s *Script) {
		s.The("actors",
			Can("eat it").And("eating it").RequiresOne("prop"),
			To("eat it", ReflectToTarget("report eat")),
		)

		s.The("props",
			Can("report eat").And("reporting eat").RequiresOne("actor"),
			To("report eat", func(g G.Play) {
				if actor := g.The("actor"); g.The("player") == actor {
					g.Say("That's not something you can eat.")
				} else {
					actor.Go("impress")
				}
			}),
		)
		s.Execute("eat it", Matching("eat {{something}}"))
	})
}
