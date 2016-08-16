package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard/live"
)

func init() {
	AddScript(func(s *Script) {
		s.The("stories",
			Can("print the banner").
				And("printing the banner").RequiresNothing(),

			To("print the banner", func(g G.Play) {
				story := g.The("story")
				name := story.Text("name")
				headline := story.Text("headline")
				if !(len(headline) > 0) {
					headline = "An Interactive fiction" // FIX: default for headline in class.
				}
				author := story.Text("author")
				g.Say(name + ".")
				g.Say(headline, "by", author+".")
				g.Say(VersionString)
			}))
	})
}
