package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard/live"
)

func init() {
	AddScript(func(s *Script) {
		s.The("actors",
			Can("attack it").And("attacking it").RequiresOne("object"),
			To("attack it", func(g G.Play) { ReflectToTarget(g, "report attack") }))

		s.The("objects",
			Can("report attack").And("reporting attack").RequiresOne("actor"),
			To("report attack", func(g G.Play) {
				if actor := g.The("actor"); g.The("player") == actor {
					g.Say("Violence isn't the answer.")
				}
			}))

		s.Execute("attack it", Matching("attack|break|smash|hit|fight|torture {{something}}").
			Or("wreck|crack|destroy|murder|kill|punch|thump {{something}}"))
	})
}
