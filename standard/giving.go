package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func GoGive(g G.Play) GivePropPhrase {
	return GivePropPhrase{g: g}
}

func (give GivePropPhrase) Prop(prop G.IObject) GivePropToPhrase {
	give.prop = prop
	return GivePropToPhrase(give)
}
func (give GivePropPhrase) The(prop string) GivePropToPhrase {
	return give.Prop(give.g.The(prop))
}

func (give GivePropToPhrase) To(actor G.IObject) {
	assignTo(give.prop, "owner", actor)
}

func (give GivePropToPhrase) ToThe(actor string) {
	give.To(give.g.The(actor))
}

type givePhraseData struct {
	g    G.Play
	prop G.IObject
}
type GivePropPhrase givePhraseData
type GivePropToPhrase givePhraseData

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
			To("give it to", ReflectWithContext("report give")),
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
					g.Say("You aren't holding", prop.Text("Name"))
					g.StopHere()
				}
			}),
		)
		// 2. receiver
		s.The("actors",
			Can("report give").And("reporting give").RequiresOne("prop").AndOne("actor"),
			To("report give", ReflectWithContext("report gave")))
		// 3. context
		s.The("props",
			Can("report gave").And("reporting gave").RequiresTwo("actor"),
			To("report gave", func(g G.Play) {
				// FIX: should generate a report/response
				receiver := g.The("action.Context")
				g.Say(receiver.Text("Name"), "is unimpressed.")
			}))
		// input
		s.Execute("give it to",
			Matching("give|pay|offer|feed {{something}} {{something else}}").
				Or("give|pay|offer|feed {{something else}} to {{something}}"))
	})
}
