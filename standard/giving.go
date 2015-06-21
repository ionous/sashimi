package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

// all infom giving rules:
// 	"applies to one carried thing and one visible thing."
//  "can't give what you haven't got"
//  "can't give to yourself"
//  "can't give to a non-person"
//  "can't give clothes being worn"
//  "block giving rule"
//  "the actor doesnt seem interested"
//  "can't exceed carrying capacity when giving"
//  "carry out giving something to"
//  "report an an actor giving something to"
func init() {
	AddScript(func(s *Script) {
		// 1. source
		s.The("actors",
			Can("give it to").And("giving it to").RequiresOne("actor").AndOne("prop"),
			To("give it to", func(g G.Play) {
				presenter, receiver, prop := g.The("action.Source"), g.The("action.Target"), g.The("action.Context")
				receiver.Go("give to", prop, presenter)
			}),
			// "convert give to yourself to examine"
			WhenCapturing("giving it to", func(g G.Play) {
				presenter, receiver := g.The("action.Source"), g.The("action.Target")
				if presenter == receiver {
					g.Say("You can't give to yourself")
					g.StopHere()
				}
			}),
			// "can't give clothes being worn"
			WhenCapturing("giving it to", func(g G.Play) {
				prop := g.The("action.Context")
				if worn := prop.Object("wearer"); worn.Exists() {
					g.Say("You can't give worn clothing.")
					// FIX: try taking off the noun
					g.StopHere()
				}
			}),
			// "you can't give what you haven't got"
			WhenCapturing("giving it to", func(g G.Play) {
				presenter, prop := g.The("action.Source"), g.The("action.Context")
				if carrier, ok := Carrier(prop); !ok || carrier != presenter {
					g.Say("You aren't holding", prop.Name())
					g.StopHere()
				}
			}),
		)
		// 2. receiver
		s.The("actors",
			Can("give to").And("giving to").RequiresOne("prop").AndOne("actor"),
			To("give to", func(g G.Play) {
				receiver, prop, presenter := g.The("action.Source"), g.The("action.Target"), g.The("action.Context")
				prop.Go("give", presenter, receiver)
			}))
		// 3. context
		s.The("props",
			Can("give").And("giving").RequiresTwo("actor"),
			To("give", func(g G.Play) {
				// FIX: should generate a report/response
				receiver := g.The("action.Context")
				g.Say(receiver.Name(), "is unimpressed.")
			}))
		// input
		s.Execute("give it to",
			Matching("give|pay|offer|feed {{something}} {{something else}}").
				Or("give|pay|offer|feed {{something else}} to {{something}}"))
	})
}
