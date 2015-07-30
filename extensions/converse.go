package extensions

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/standard"
)

func (con *Conversation) Converse(g G.Play) {
	if npc, ok := con.Interlocutor.Get(); ok {
		if standard.Debugging {
			fmt.Println("conversing...")
		}
		currentQuip := con.History.MostRecent(g)
		currentRestricts := currentQuip.Exists() && currentQuip.Is("restrictive")
		// handle queued conversation, unless the current quip is restrictive.
		if !currentRestricts {
			con.Queue.UpdateNextQuips(g, con.Memory)
		}

		// process the current speaker first:
		con.perform(npc, currentRestricts)

		// process anyone else who might have something to say:
		g.Visit("actors", func(actor G.IObject) (okay bool) {
			// threaded conversation tests:
			// repeat with target running through **visible** people who are not the player:
			if actor != npc {
				con.perform(actor, currentRestricts)
			}
			return okay
		})

		// we might have changed conversations...
		if npc, ok := con.Interlocutor.Get(); ok {
			g.The("player").Go("print conversation choices", npc)
		}
	}
}

func (con *Conversation) perform(actor G.IObject, currentRestricts bool) {
	if nextQuip := actor.Object("next quip"); nextQuip.Exists() {
		if !currentRestricts || nextQuip.Is("planned") {
			quip := nextQuip.Object("quip")
			talker := quip.Object("subject")
			talker.Go("discuss", quip)
		}
		// this removes the planned conversation which was just said,
		// and any casual conversation that couldn't be said due to restriction.
		nextQuip.Remove()
	}
}
