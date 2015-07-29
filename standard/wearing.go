package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func Clothe(actor string) ClothePhrase {
	return ClothePhrase{actor: actor}
}

func (p ClothePhrase) With(prop string) WearingPhrase {
	p.clothing = prop
	return WearingPhrase(p)
}

func (p WearingPhrase) Execute(g G.Play) {
	actor, clothing := g.The(p.actor), g.The(p.actor)
	assignTo(clothing, "wearer", actor)
}

type wearData struct {
	actor, clothing string
}
type ClothePhrase wearData
type WearingPhrase wearData

func init() {
	AddScript(func(s *Script) {
		// FIX? the action definition effectively establishes a new way of declaring properties
		// it requires their own special storage, even.
		// why? just like tables -- this could be done with a class.
		// and then we could store anything, text, etc. just as was desired.
		// you could start by fixing this internally, and then come back to change the requires interface.
		s.The("actors",
			Can("wear it").And("wearing it").RequiresOne("prop"),
			To("wear it", ReflectToTarget("report wear")),
		)

		s.The("props",
			AreEither("wearable").Or("not wearable").Usually("not wearable"))

		s.The("props",
			Can("report wear").And("reporting wear").RequiresOne("actor"),
			To("report wear", func(g G.Play) {
				actor, prop := g.The("actor"), g.The("prop")
				if prop.Is("not wearable") {
					g.Say("That's not something you can wear.")
				} else {
					g.Go(Clothe("actor").With("prop"))
					g.Say("Now", actor.Text("name"), "is wearing the", prop.Text("name"))
				}
			}))

		s.Execute("wear it",
			Matching("wear|don {{something}}").
				Or("put on {{something}}").
				Or("put {{something}} on"))
	})
}
