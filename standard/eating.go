package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard/live"
)

func init() {
	AddScript(func(s *Script) {
		s.The("actors",
			Can("eat it").And("eating it").RequiresOne("prop"),
			To("eat it", func(g G.Play) { ReflectToTarget(g, "report eat") }),
		)

		s.The("props", AreEither("inedible").Or("edible"))

		s.The("props",
			Can("report eat").And("reporting eat").RequiresOne("actor"),
			To("report eat", func(g G.Play) {
				if g.The("prop").Is("inedible") {
					g.Say("That's not something you can eat.")
				} else {
					g.The("actor").Go("impress")
				}
			}),
		)
		s.Execute("eat it", Matching("eat {{something}}"))
	})
}
