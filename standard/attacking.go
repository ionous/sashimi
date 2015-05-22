package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func init() {
	AddScript(func(s *Script) {
		s.The("actors",
			Can("attack it").And("attacking it").RequiresOne("object"),
			To("attack it", actorTarget("attack")))

		s.The("objects",
			Can("attack").And("attacking").RequiresOne("actor"),
			To("attack", func(g G.Play) {
				g.Say("violence isn't the answer")
			}))

		s.Execute("attack it", Matching("attack|break|smash|hit|fight|torture {{something}}").
			Or("wreck|crack|destroy|murder|kill|punch|thump {{something}}"))
	})
}
