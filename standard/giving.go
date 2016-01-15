package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func Give(prop string) GivePropPhrase {
	return GivePropPhrase{prop: prop}
}

func GiveThe(prop G.IObject) GivePropPhrase {
	return GivePropPhrase{prop: prop.Id().String()}
}

func (give GivePropPhrase) To(actor string) GivingPhrase {
	give.actor = actor
	return GivingPhrase(give)
}

func (give GivingPhrase) Execute(g G.Play) {
	prop, actor := g.The(give.prop), g.The(give.actor)
	//added indirection so we can transform props after the rules of taking/giving have run
	actor.Go("acquire it", prop)
}

type givePhraseData struct {
	prop, actor string
}
type GivePropPhrase givePhraseData
type GivingPhrase givePhraseData

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
		s.The("actors",
			Can("acquire it").And("acquiring it").RequiresOne("prop"),
			To("acquire it", func(g G.Play) { ReflectToTarget(g, "be acquired") }))
		s.The("props",
			Can("be acquired").And("being acquired").RequiresOne("actor"),
			To("be acquired", func(g G.Play) {
				actor, prop, rel := g.The("actor"), g.The("prop"), "owner"
				if Debugging {
					par, prev := prop.ParentRelation()
					g.Log(prop, "AssignTo", actor, rel, "from", par, prev)
				}
				AssignTo(prop, rel, actor)
			}))
		// 1. source
		s.The("actors",
			Can("give it to").And("giving it to").RequiresOne("actor").AndOne("prop"),
			To("give it to", func(g G.Play) { ReflectWithContext(g, "report give") }),
			// "convert give to yourself to examine"
			Before("giving it to").Always(func(g G.Play) {
				presenter, receiver := g.The("action.Source"), g.The("action.Target")
				if presenter == receiver {
					g.Say("You can't give to yourself")
					g.StopHere()
				}
			}),
			// "can't give clothes being worn"
			Before("giving it to").Always(func(g G.Play) {
				prop := g.The("action.Context")
				if worn := prop.Object("wearer"); worn.Exists() {
					g.Say("You can't give worn clothing.")
					// FIX: try taking off the noun
					g.StopHere()
				}
			}),
			// "you can't give what you haven't got"
			Before("giving it to").Always(func(g G.Play) {
				presenter, prop := g.The("action.Source"), g.The("action.Context")
				if carrier := Carrier(prop); carrier != presenter {
					g.Say("You aren't holding", ArticleName(g, "action.Context", NameFullStop))
					g.StopHere()
				}
			}),
		)
		// 2. receiver
		s.The("actors",
			Can("report give").And("reporting give").RequiresOne("prop").AndOne("actor"),
			To("report give", func(g G.Play) { ReflectWithContext(g, "report gave") }))
		// 3. context
		s.The("props",
			Can("report gave").And("reporting gave").RequiresTwo("actor"),
			To("report gave", func(g G.Play) {
				receiver := g.The("action.Context")
				receiver.Go("impress")
			}))
		// input
		s.Execute("give it to",
			Matching("give|pay|offer|feed {{something}} {{something else}}").
				Or("give|pay|offer|feed {{something else}} to {{something}}"))
	})
}
