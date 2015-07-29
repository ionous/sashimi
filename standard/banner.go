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
				g.Say(`{{ $src := action.Source }}
{{ $src.Name }}
{{ if $src.Headline }}{{ $src.Headline }}{{else}}An Interactive fiction{{end}} by {{ $src.Author }}`)
				g.Say(VersionString)
			}))
	})
}
