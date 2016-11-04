package standard

import (
	"github.com/ionous/mars/script/g"
	"github.com/ionous/mars/std"
	. "github.com/ionous/sashimi/script"
)

func init() {
	AddScript(func(s *Script) {
		s.The("stories",
			Can("print the banner").
				And("printing the banner").RequiresNothing(),

			To("print the banner",
				g.Context{g.Our("story"),
					g.Statements{
						g.Say(g.GetText{"name"}, "."),
						g.Say(g.ChooseText{
							If:    g.Empty{g.GetText{"headline"}},
							False: g.GetText{"headline"},
							// FIX: default for headline in class.
							True: g.MakeText("An interactive fiction"),
						}, "by", g.GetText{"author"}, "."),
						g.Say(std.VersionString),
					},
				}))
	})
}
