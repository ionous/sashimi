package extensions

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	_ "github.com/ionous/sashimi/standard"
)

func init() {
	AddScript(func(s *Script) {
		s.The("kinds",
			Called("quips"),
			//Comment comes first in dialog, and is said by the player.
			//If there is no comment then it is considered “npc-directed”.
			//For instance, a greeting when the player selects an NPC.
			Have("comment", "text"),
			Have("speaker", "actor"),
			Have("reply", "text"),
			//Have("hook", "text"), // displayed on the menu
			//performative, informative, questioning: used for ask about, tell about, or simply state the quip name
			AreEither("repeatable").Or("one time").Usually("one time"),
			AreEither("restrictive").Or("unrestricted").Usually("unrestricted"),
		//really important, unimportant, ...: from my extension to add priority sorting
		)

		// FIX: data not kinds.
		s.The("kinds",
			Called("following quips"),
			Have("leading", "quip"),
			Have("following", "quip"),
			AreEither("indirectly following").Or("directly following"),
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

		// this is overspecification --
		// we've got recollection after all.
		// FIX: we just need fast sorting.
		qh := QuipHistory{}

		s.The("actors",
			Have("greeting", "quip"))

		s.The("actors",
			Can("greet").And("greeting").RequiresOne("actor"),
			To("greet", func(g G.Play) {
				greeter, speaker := g.The("action.Source"), g.The("action.Target")
				if greeter == g.The("player") {
					quips := GetQuipPool(g)
					if lastQuip, ok := quips.MostRecent(qh); !ok {
						if greeting := greeter.Object("greeting"); !greeting.Exists() {
						} else {
							qh.Push(greeting.Id())
						}
						// QUEUE!
					} else {
						interlocutor := quips.Interlocutor(lastQuip)
						if interlocutor == speaker.Id() {
							g.Say("You're already speaking to them!")
						} else {
							g.Say("You're already speaking to someone!")
						}
					}
				}
			}))

		s.The("actors",
			Can("depart").And("departing").RequiresNothing(),
			To("depart", func(g G.Play) {
				qh.Clear()
			}))

		s.The("stories",
			When("ending the turn").Always(func(g G.Play) {
				// have a turn of thinking

				// present player choices
				//qp := GetQuipPool(g)
				//list := qp.GetPlayerQuips(qh)
			}))

		s.The("quips",
			Can("discuss").And("discussing").RequiresNothing(),
			To("discuss", func(g G.Play) {
				quip := g.The("quip")
				queue := g.Add("queued quip")
				// default is immediate/obligatory.]
				// NOTE: inform doesnt have event parameters.
				// threaded convesation simply exposes a handful of different actions to queue things in different ways.
				queue.Set("quip", quip)
			}))

		s.The("kinds",
			Called("queued quips"),
			Have("quip", "quip"),
			AreEither("immediate").Or("postponed"),
			AreEither("obligatory").Or("optional"),
		)

		// discussing action to choose and exectue a line of dialog
		// saying hello to queue an npc greeting
		// present player choices
		// inject player choice as discussing
		// move spoken quips to the recollection table.
		// check that every dialog line has an npc
		// check that every npc has a greeting
		// perhaps a Requires for Has
	})
}
