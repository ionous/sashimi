package extensions

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/standard"
)

func (con *Conversation) Depart() (wasTalking G.IObject) {
	if npc, ok := con.Interlocutor.Get(); ok {
		con.History.ClearQuips()
		con.Queue.ResetQuipQueue()
		npc.Object("next quip").Remove()
		wasTalking = npc
		con.Interlocutor.Clear()
	}
	return wasTalking
}

func (con *Conversation) Converse(g G.Play) {
	if npc, ok := con.Interlocutor.Get(); ok {
		if standard.Debugging {
			g.Log("conversing...")
		}
		currentQuip := con.History.MostRecent(g)
		currentRestricts := currentQuip.Exists() && currentQuip.Is("restrictive")
		// handle queued conversation, unless the current quip is restrictive.
		if !currentRestricts {
			con.Queue.UpdateNextQuips(g, con.Memory)
		}

		// sometimes conversations want to loop, without the player saying anything; for now this doesnt advance the queue -- just looks at next quips; may need to be re-visited.
		for again := true; again; {
			again = false
			// process the current speaker first:
			if con.perform(npc, currentRestricts) {
				again = true
			}

			// process anyone else who might have something to say:
			g.Visit("actors", func(actor G.IObject) (okay bool) {
				// threaded conversation tests:
				// repeat with target running through **visible** people who are not the player:
				if actor != npc {
					if con.perform(actor, currentRestricts) {
						again = true
					}
				}
				return okay
			})
		}

		// we might have changed conversations...
		if npc, ok := con.Interlocutor.Get(); ok {
			g.The("player").Go("print conversation choices", npc)
		}
	}
}

func (con *Conversation) perform(actor G.IObject, currentRestricts bool,
) (
	spoke bool,
) {
	if nextQuip := actor.Object("next quip"); nextQuip.Exists() {
		if !currentRestricts || nextQuip.Is("planned") {
			quip := nextQuip.Object("quip")
			actor.Go("discuss", quip)
			spoke = true
		}
		// this removes the planned conversation which was just said,
		// and any casual conversation that couldn't be said due to restriction.
		nextQuip.Remove()
	}
	return spoke
}
