package standard

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

// from inform:
//  touchable/reach inside checks --, plus:
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
			To("take it", actorTarget("take")),
		)

		s.The("props",
			Can("take").And("taking").RequiresOne("actor"),
			To("take", func(g G.Play) {
				prop, actor := g.The("action.Source"), g.The("action.Target")
				// first, only same room:
				actorCeiling, _ := Enclosure(actor)
				targetCeiling, _ := Enclosure(prop)
				//
				if actorCeiling != targetCeiling {
					g.Say("That isn't available.")
					g.Log(fmt.Sprintf("take ceiling mismatch (%v!=%v)", actorCeiling, targetCeiling))
				} else {
					if prop.Is("scenery") {
						g.Say("That isn't available.")
						g.Log("can't take scenery")
						//g.StopHere() // FIX: should be cancel
						return
					}
					if prop.Is("fixed in place") {
						g.Say("It is fixed in place.")
						//g.StopHere() // FIX: should be cancel
						return
					}
					parent, rel := DirectParent(prop)
					if rel != "" {
						if parent.Class("actor") {
							if parent != actor {
								g.Say("That'd be stealing!")
							} else {
								g.Say("{{action.Target.Name}} already has that!")
							}
							return
						}
					}
					Give(actor, prop)
					// separate report action?
					if actor == g.The("player") {
						g.Say("You take the {{action.Source.Name}}.")
					} else {
						g.Say("{{action.Target.Name}} takes the {{action.Source.Name}}.")
					}
				}
			}))
		// understandings:
		s.Execute("take it",
			Matching("take|get {{something}}").
				Or("pick up {{something}}").
				Or("pick {{something}} up"))
	})
}
