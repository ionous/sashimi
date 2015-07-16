package extensions

import (
	G "github.com/ionous/sashimi/game"
)

func Converse(g G.Play) {
	con := g.Global("conversation").(*Conversation)
	if npc, ok := con.Interlocutor.Get(); ok {

		currentQuip := con.History.MostRecent(g)
		currentRestricts := currentQuip.Exists() && currentQuip.Is("restrictive")
		// handle queued conversation, unless the current quip is restrictive.
		if !currentRestricts {
			con.Queue.UpdateNextQuips(g, con.Memory)
		}

		perform := func(actor G.IObject) {
			if nextQuip := actor.Object("next quip"); nextQuip.Exists() {
				if !currentRestricts || nextQuip.Is("planned") {
					quip := nextQuip.Object("quip")
					//
					talker := quip.Object("subject")
					reply := quip.Text("reply")
					talker.Says(reply)
					con.Memory.Learn(quip)
				}
				// this removes the planned conversation which was just said,
				// and any casual conversation that couldn't be said due to restriction.
				nextQuip.Remove()
			}
		}

		// process the current speaker first:
		perform(npc)

		// process anyone else who might have something to say:
		g.Visit("actors", func(actor G.IObject) (okay bool) {
			// threaded conversation tests:
			// repeat with target running through **visible** people who are not the player:
			if actor != npc {
				perform(actor)
			}
			return okay
		})
	}
}
