package extensions

import (
	G "github.com/ionous/sashimi/game"
)

func Converse(g G.Play, qh QuipHistory) {

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
				//
				talker := quip.Object("speaker")
				reply := quip.Text("reply")
				talker.Says(reply)
				LearnQuip(g, quip)
			}
			// this removes the planned conversation which was just said,
			// and any casual conversation that couldn't be said due to restriction.
			nextQuip.Remove()
		}
	}

	// process the current speaker first:
	interlocutor := currentQuip.Object("speaker")
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
