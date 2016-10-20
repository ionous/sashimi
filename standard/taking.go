package standard

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard/live"
)

// from inform:
// touchable/reach inside checks --, plus:
// can't take yourself: not applicable, can only take props
// can't take other people: not applicable, can only take props
// can't take component parts: not applicable... yet?
// can't take other actor's possessions: not applicable... yet?
// can't take items out of play
// can't take what you're inside: not applicable... yet?
// can't take what's already taken ( carrying or worn )
// can't take scenery
// can only take things: not applicable... yet?
// can't take what's fixed in place
// player's carry all?
// carrying capacity
func init() {
	AddScript(func(s *Script) {
		s.The("actors",
			Can("take it").And("taking it").RequiresOne("prop"),
			To("take it", ReflectToTarget("report take")),
		)

		s.The("props",
			Can("report take").And("reporting take").RequiresOne("actor"),
			To("report take", func(g G.Play) {
				prop, actor := g.The("prop"), g.The("actor")
				// first, only same room:
				actorCeiling, targetCeiling := Enclosure(actor), Enclosure(prop)
				//
				if !actorCeiling.Equals(targetCeiling) {
					g.Say("That isn't available.")
					g.Log(fmt.Sprintf("take ceiling mismatch (%v not %v)", actorCeiling, targetCeiling))
				} else {
					if prop.Is("scenery") {
						g.Say("That isn't available.")
						g.Log("(You can't take scenery.)")
						//g.StopHere() // FIX: should be cancel
						return
					}
					if prop.Is("fixed in place") {
						g.Say("It is fixed in place.")
						//g.StopHere() // FIX: should be cancel
						return
					}
					if parent, _ := prop.ParentRelation(); parent.Exists() {
						if parent.FromClass("actors") {
							if !parent.Equals(actor) {
								g.Say("That'd be stealing!")
							} else {
								g.Say(ArticleName(g, "action.Target", nil), "already has that!")
							}
							return
						}
					}
					// separate report action?
					if actor.Equals(g.The("player")) {
						// the kat food.
						g.Say("You take", DefiniteName(g, "action.Source", NameFullStop))
					} else {
						g.Say(ArticleName(g, "action.Target", nil), "takes", DefiniteName(g, "action.Source", NameFullStop))
					}
					g.Go(Give("prop").To("actor"))
				}
			}))
		// understandings:
		s.Execute("take it",
			Matching("take|get {{something}}").
				Or("pick up {{something}}").
				Or("pick {{something}} up"))
	})
}
