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
			To("kiss it", func(g G.Play) { ReflectToTarget(g, "report kiss") }),
			//  kissing yourself rule
			WhenCapturing("kissing it", func(g G.Play) {
				source, target := g.The("action.Source"), g.The("action.Target")
				if source == target {
					g.Say(source.Text("Name"), "didn't get much from that.")
					g.StopHere()
				}
			}),
		)
		//  block kissing rule
		s.The("objects",
			Can("report kiss").And("reporting kiss").RequiresOne("actor"),
			To("report kiss", func(g G.Play) {
				source := g.The("action.Source")
				g.Say(source.Text("Name"), "might not like that.")
			}))

		s.Execute("kiss it", Matching("kiss|hug|embrace {{something}}"))
	})
}
