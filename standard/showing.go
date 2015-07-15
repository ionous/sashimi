package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
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
			To("show it to", ReflectWithContext("report show")),
			// "you can't show what you haven't got"
			WhenCapturing("showing it to", func(g G.Play) {
				presenter, _, prop := g.The("action.Source"), g.The("action.Target"), g.The("action.Context")
				if carrier, ok := Carrier(prop); !ok || carrier != presenter {
					g.Say("You aren't holding", prop.Text("Name"))
					g.StopHere()
				}
			}),
			// "convert show to yourself to examine"
			WhenCapturing("showing it to", func(g G.Play) {
				presenter, receiver, prop := g.The("action.Source"), g.The("action.Target"), g.The("action.Context")
				if presenter == receiver {
					presenter.Go("examine it", prop)
					g.StopHere()
				}
			}),
		)
		// 2. receiver
		s.The("actors",
			Can("report show").And("reporting show").RequiresOne("prop").AndOne("actor"),
			To("report show", ReflectWithContext("report shown")))
		// 3. context
		s.The("props",
			Can("report shown").And("reporting shown").RequiresTwo("actor"),
			To("report shown", func(g G.Play) {
				_, _, receiver := g.The("action.Source"), g.The("action.Target"), g.The("action.Context")
				g.Say(receiver.Text("Name"), "is unimpressed")
			}))
		// input
		s.Execute("show it to",
			Matching("show|present|display {{something}} {{something else}}").
				Or("show|present|display {{something else}} to {{something}}"))
	})
}
