package extensions

import (
	G "github.com/ionous/sashimi/game"
)

func Converse(g G.Play, interlocutor G.IObject, qh QuipHistory) {
	if interlocutor.Exists() {
		currentQuip := qh.MostRecent(g)
		currentRestricts := currentQuip.Exists() && currentQuip.Is("restrictive")
		// handle queued conversation, unless the current quip is restrictive.
		if !currentRestricts {
			QuipQueue.UpdateNextQuips(g)
		}

		perform := func(actor G.IObject) {
			if nextQuip := actor.Object("next quip"); nextQuip.Exists() {
				if !currentRestricts || nextQuip.Is("planned") {
					quip := nextQuip.Object("quip")
					//
					talker := quip.Object("subject")
					reply := quip.Text("reply")
					talker.Says(reply)
					QuipMemory.Learn(quip)
				}
				// this removes the planned conversation which was just said,
				// and any casual conversation that couldn't be said due to restriction.
				nextQuip.Remove()
			}
		}

		// process the current speaker first:
		perform(interlocutor)

		// process anyone else who might have something to say:
		g.Visit("actors", func(actor G.IObject) (okay bool) {
			// threaded conversation tests:
			// repeat with target running through **visible** people who are not the player:
			if actor != interlocutor {
				perform(actor)
			}
			return okay
		})
	}
}
