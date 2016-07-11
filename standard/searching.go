package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard/live"
)

func init() {
	AddScript(func(s *Script) {
		// searching: requiring light; FIX: what does searching a room do?
		s.The("actors",
			Can("search it").And("searching it").RequiresOne("prop"),
			To("search it", func(g G.Play) { ReflectToTarget(g, "report search") }))
		s.The("props",
			Can("report search").And("reporting search").RequiresOne("actor"),
			To("report search", func(g G.Play) {
				g.Say("You find nothing unexpected.")
			}))

		// WARNING/FIX: multi-word statements must appear before their single word variants
		// ( or the parser will attempt to match the setcond word as a noun )
		s.Execute("search it", Matching("search {{something}}").
			Or("look inside|in|into|through {{something}}"))
	})
}
