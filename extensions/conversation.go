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
					if lastQuip := qh.MostRecent(g); lastQuip.Exists() {
						if greeting := greeter.Object("greeting"); greeting.Exists() {
							qh.PushQuip(greeting)
							QueueQuip(g, greeting)
						}
					} else {
						if speaker == lastQuip.Object("speaker") {
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
				qh.ClearQuips()
			}))

		s.The("stories",
			When("ending the turn").Always(func(g G.Play) {
				currentQuip := qh.MostRecent(g)
				currentRestricts := currentQuip.Is("restrictive")
				// handle queued conversation, unless the current quip is restrictive.
				if !currentRestricts {
					// from slice tricks, this reuses the memory of the quip queue
					requeue := quipQueue[:0]
					// determine what to say next
					// note: queued conversation will never override what an npc already has to say.
					for _, quip := range quipQueue {
						npc := quip.Object("Speaker")
						if npc.Object("next quip").Exists() {
							requeue = append(requeue, quip)
						} else {
							// check to make sure this quip wasn't said in the time since it was queued.
							if quip.Is("repeatable") || !Recollects(g, quip) {
								nextQuip := g.Add("next quip")
								npc.Set("next quip", nextQuip)
								nextQuip.Set("quip", quip)
								nextQuip.SetIs("casual")
							}
						}
					}
					quipQueue = requeue
				}

				perform := func(actor G.IObject) {
					if nextQuip := actor.Object("next quip"); nextQuip.Exists() {
						if !currentRestricts || nextQuip.Is("planned") {
							quip := nextQuip.Object("quip")
							actor.Go("discuss", quip)
						}
						// this removes the planned conversation which was just said,
						// and any casual conversation that couldnt be said due to restriction.
						nextQuip.Remove()
					}
				}

				// process the current speaker first:
				currentInterlocutor := currentQuip.Object("speaker")
				perform(currentInterlocutor)

				// process anyone else who might have something to say:
				g.Visit("actors", func(actor G.IObject) (okay bool) {
					// threaded conversation tests:
					// repeat with target running through **visible** people who are not the player:
					if actor != currentInterlocutor {
						perform(actor)
					}
					return okay
				})
			}))

		s.The("actors",
			Can("discuss").And("discussing").RequiresOne("quip"),
			To("discuss", func(g G.Play) {
				actor, quip := g.The("actor"), g.The("quip")
				if quip.Object("speaker") == actor {
					reply := quip.Text("reply")
					actor.Says(reply)
					LearnQuip(g, quip)
				}
			}))

		s.The("kinds",
			Called("next quips"),
			Have("quip", "quip"),
			// these have the same meaning as "immediate obligatory" and "immediate optional".
			// casual quips lose their relevance whenever a new casual or planned quip is set,
			// and the moment after a player has spoken ( stop any planned casual follow-ups ).
			//
			// note: there is a gap in the original logic --
			// if the current quip is restrictive and if the person isnt the current interlocutor,
			// then the immediate optional conversation doesn't cleared; it sticks in there until the player chooses some unrestrictive quip.
			// but: it's difficult to get immediate conversation assigned to a person who isnt the current interlocutor
			// because the shortcuts always refer to the current interlocutor.
			// the gap is likely an oversight.
			AreEither("planned").Or("casual"),
		)
		s.The("actors",
			// FIX? with pointers, it wouldnt be too difficult to have parts now
			// an auto-created association.
			Have("next quip", "next quip"))

		// x. discussing action to choose and execute a line of dialog
		// saying hello to queue an npc greeting
		// present player choices
		// inject player choice as discussing
		// x. move spoken quips to the recollection table.
		// check that every dialog line has an npc
		// check that every npc has a greeting
		// perhaps a Requires for Has
	})
}
