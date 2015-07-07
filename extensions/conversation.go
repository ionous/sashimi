package extensions

import (
	G "github.com/ionous/sashimi/game"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
)

func init() {
	AddScript(func(s *Script) {
		s.The("kinds",
			Called("quips"),
			//Comment comes first in dialog, and is said by the player.
			//If there is no comment then it is considered “npc-directed”.
			//For instance, a greeting when the player selects an NPC.
			Have("comment", "text"),
			Have("responder", "actor"),
			Have("reply", "text"),
			//Have("hook", "text"), // displayed on the menu
			//performative, informative, questioning: used for ask about, tell about, or simply state the quip name
			AreEither("repeatable").Or("one time").Usually("one time"),
			AreEither("restrictive").Or("unrestricted").Usually("unrestricted"),
		//really important, unimportant, ...: from my extension to add priority sorting
		)
		s.The("kinds",
			Called("following quips"),
			Have("leading", "quip"),
			Have("following", "quip"),
			AreOneOf("indirectly following", "directly following", "tangential"),
		)
		s.The("kinds",
			Called("pending quips"),
			Have("speaker", "actor"),
			AreEither("immediate").Or("postponed"),
			AreEither("obligatory").Or("optional"),
		)
		// FIX: may relate to actors
		s.The("kinds",
			Called("recollections"),
			Have("quip", "quip"),
		)
		s.The("kinds", Called("facts"),
			// FIX: interestingly, kinds should have names
			// having the same property as a parent class probably shouldnt be an error
			Have("summary", "text"))

		var interlocutor G.IObject

		s.The("actors",
			Can("greet").And("greeting").RequiresOne("actor"),
			To("greet", func(g G.Play) {
				speaker, receiver := g.The("action.Source"), g.The("action.Target")
				if speaker == g.The("player") {
					interlocutor = receiver
					//				receiver.
				}
			}))

		s.The("actors",
			Can("depart").And("departing").RequiresNothing(),
			To("depart", func(g G.Play) {
				// if speaker == g.The("player") {

				// }
			}))

		s.The("stories",
			When("ending the turn").Always(func(g G.Play) {
				// Filter to quips which have player comments.
				// Filter to quips which quip supply the interlocutor.
				// Exclude one-time quips, checking the recollection table.
				// Filter restrictive quips to limit available responses to those which directly follow. And, the opposite of restrictive quips: include those which indirectly follow.
				// Include/exclude quips based on the recollection of others.
				// Include/exclude quips based on the knowledge of facts.
				adapt := g.(*R.GameEventAdapter)
				for _, o := range adapt.Game.Model.Instances {
					if isQuip := o.Class().CompatibleWith("Quip"); isQuip {
						o := R.NewObjectAdapter(adapt.Game, o)
						if o.Text("comment") != "" {
							//	if o.Text
						}
					}
				}
			}))
	})
}
