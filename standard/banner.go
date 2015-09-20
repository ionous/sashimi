package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

//
const VersionString = "Sashimi Experimental IF Engine - 0.1"

func init() {
	AddScript(func(s *Script) {
		s.The("stories",
			Can("print the banner").
				And("printing the banner").RequiresNothing(),

			To("print the banner", func(g G.Play) {
				story := g.The("story")
				name := story.Text("name")
				headline := story.Text("headline")
				if headline == "" {
					headline = "An Interactive fiction" // FIX: default for headline in class.
				}
				author := story.Text("author")
				g.Println(name)
				g.Println(headline, "by", author)
				g.Println(VersionString)
			}))
	})
}
