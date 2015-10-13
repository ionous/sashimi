package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func Put(prop string) PutPhrase {
	return PutPhrase{prop: prop}
}

func (p PutPhrase) Onto(supporter string) PutingPhrase {
	p.supporter = supporter
	return PutingPhrase(p)
}

func (p PutingPhrase) Execute(g G.Play) {
	prop, supporter := g.The(p.prop), g.The(p.supporter)
	assignTo(prop, "enclosure", supporter)
}

type putData struct {
	prop, supporter string
}

type PutPhrase putData
type PutingPhrase putData

func init() {
	AddScript(func(s *Script) {
		// 1. source
		s.The("actors",
			// FIX? word-wise this is wrong ( see tickle-it-with, though it is "correct" )
			Can("put it onto").And("putting it onto").RequiresOne("supporter").AndOne("prop"),
			To("put it onto", ReflectWithContext("report put")),
			//  can't put clothes being worn
			WhenCapturing("putting it onto", func(g G.Play) {
				prop := g.The("action.Context")
				if worn := prop.Object("wearer"); worn.Exists() {
					g.Say("You can't put worn clothing.")
					// FIX: try taking off the noun
					g.StopHere()
				}
			}),
			//  can't put what isn't held
			WhenCapturing("putting it onto", func(g G.Play) {
				actor, prop := g.The("action.Source"), g.The("action.Context")
				if carrier := Carrier(prop); carrier != actor {
					g.Say("You aren't holding", ArticleName(g, "action.Context", NameFullStop))
					g.StopHere()
				}
			}),
			//  can't put something onto itself
			WhenCapturing("putting it onto", func(g G.Play) {
				supporter, prop := g.The("action.Target"), g.The("action.Context")
				if supporter == prop {
					g.Say("You can't put something onto itself.")
					g.StopHere()
				}
			}),
			//  can't put onto closed supporters
			WhenCapturing("putting it onto", func(g G.Play) {
				supporter := g.The("action.Target")
				if supporter.Is("closed") {
					g.Say(ArticleName(g, "action.Target", nil), "is closed.")
					g.StopHere()
				}
			}),
		)
		// 2. supporters
		s.The("supporters",
			Can("report put").And("reporting put").RequiresOne("prop").AndOne("actor"),
			To("report put", ReflectWithContext("report placed")))
		// 3. context
		s.The("props",
			Can("report placed").And("reporting placed").RequiresOne("actor").AndOne("supporter"),
			To("report placed", func(g G.Play) {
				g.Go(Put("action.Source").Onto("action.Context"))
				g.Say("You put", ArticleName(g, "action.Source", nil), "onto", ArticleName(g, "action.Context", NameFullStop))
			}))
		// x.
		s.Execute("put it onto",
			Matching("put {{something else}} on|onto {{something}}").
				Or("drop|throw|discard {{something else}} on|onto {{something}}"))
	})
}
