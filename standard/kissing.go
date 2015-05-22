package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func init() {
	AddScript(func(s *Script) {
		// kissing
		s.The("actors",
			Can("kiss it").And("kissing it").RequiresOne("object"),
			To("kiss it", actorTarget("kiss")),
			//  kissing yourself rule
			WhenCapturing("kissing it", func(g G.Play) {
				source, target := g.The("action.Source"), g.The("action.Target")
				if source == target {
					g.Say(source.Name(), "didn't get much from that.")
					g.StopHere()
				}
			}),
		)
		//  block kissing rule
		s.The("objects",
			Can("kiss").And("kissing").RequiresOne("actor"),
			To("kiss", func(g G.Play) {
				source := g.The("action.Source")
				g.Say(source.Name(), "might not like that")
			}))

		s.Execute("kiss it", Matching("kiss|hug|embrace {{something}}"))
	})
}
