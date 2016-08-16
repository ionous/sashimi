package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard/live"
)

// all infom showing rules:
// 	"applies to one carried thing and one visible thing"
//  "you can't show what you haven't got"
//  	why does it need to do this, since it already applies to something carried?
//  "convert show to yourself to examine, tests if the actors are the same, and calls "convert examining""
//  	why it needs a special "convert" function?
//  "block showing - says: "the actor is unimpressed""
// 		why not an explicit report?
func init() {
	AddScript(func(s *Script) {
		// 1. source
		s.The("actors",
			Can("show it to").And("showing it to").RequiresOne("actor").AndOne("prop"),
			To("show it to", func(g G.Play) { ReflectWithContext(g, "report show") }),
			// "you can't show what you haven't got"
			Before("showing it to").Always(func(g G.Play) {
				presenter, _, prop := g.The("action.Source"), g.The("action.Target"), g.The("action.Context")
				if carrier := Carrier(prop); !carrier.Equals(presenter) {
					g.Say("You aren't holding", ArticleName(g, "action.Context", NameFullStop))
					g.StopHere()
				}
			}),
			// "convert show to yourself to examine"
			Before("showing it to").Always(func(g G.Play) {
				presenter, receiver, prop := g.The("action.Source"), g.The("action.Target"), g.The("action.Context")
				if presenter.Equals(receiver) {
					presenter.Go("examine it", prop)
					g.StopHere()
				}
			}),
		)
		// 2. receiver
		s.The("actors",
			Can("report show").And("reporting show").RequiresOne("prop").AndOne("actor"),
			To("report show", func(g G.Play) { ReflectWithContext(g, "report shown") }))

		// 3. context
		s.The("props",
			Can("report shown").And("reporting shown").RequiresTwo("actor"),
			To("report shown", func(g G.Play) {
				_, _, receiver := g.The("action.Source"), g.The("action.Target"), g.The("action.Context")
				receiver.Go("impress")
			}))
		// input
		s.Execute("show it to",
			Matching("show|present|display {{something}} {{something else}}").
				Or("show|present|display {{something else}} to {{something}}"))
	})
}
