package extensions

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/standard"
)

type Conversation struct {
	G.IObject
}

func TheConversation(g G.Play) Conversation {
	return Conversation{g.The("conversation")}
}

func (c Conversation) Interlocutor() G.IValue {
	return c.Get("interlocutor")
}

func (c Conversation) History() QuipHistory {
	return QuipHistory{c.Get("parent"), c.Get("grandparent"), c.Get("greatgrand")}
}

func (c Conversation) Queue() QuipQueue {
	return QuipQueue{c.List("queue")}
}

func (c Conversation) Depart() (wasTalking G.IObject) {
	interlocutor := c.Interlocutor()
	if npc := interlocutor.Object(); npc.Exists() {
		c.History().Reset()
		c.Queue().Reset()
		npc.Set("next quip", nil)
		wasTalking = npc
		interlocutor.SetObject(nil)
	}
	return wasTalking
}

func Converse(g G.Play) {
	c := TheConversation(g)
	interlocutor := c.Interlocutor()
	if npc := interlocutor.Object(); npc.Exists() {
		if standard.Debugging {
			g.Log("conversing...")
		}
		currentQuip := c.History().MostRecent()
		currentRestricts := currentQuip.Exists() && currentQuip.Is("restrictive")
		// handle queued conversation, unless the current quip is restrictive.
		if !currentRestricts {
			c.Queue().UpdateNextQuips(PlayerMemory(g))
		}

		// sometimes conversations want to loop, without the player saying anything; for now this doesnt advance the queue -- just looks at next quips; may need to be re-visited.
		for again := true; again; {
			again = false
			// process the current speaker first:
			if c.perform(npc, currentRestricts) {
				again = true
			}

			// process anyone else who might have something to say:
			for i, actors := 0, g.List("actors"); i < actors.Len(); i++ {
				actor := actors.Get(i).Object()
				// threaded conversation tests:
				// repeat with target running through **visible** people who are not the player:
				if actor != npc {
					if c.perform(actor, currentRestricts) {
						again = true
					}
				}
			}
		}

		// we might have changed conversations...
		if npc := interlocutor.Object(); npc.Exists() {
			g.The("player").Go("print conversation choices", npc)
		}
	}
}

func (c Conversation) perform(
	actor G.IObject,
	currentRestricts bool,
) (
	spoke bool,
) {
	next := actor.Get("next quip")
	if nextQuip := next.Object(); nextQuip.Exists() {
		if !currentRestricts || nextQuip.Is("planned") {
			actor.Go("discuss", nextQuip)
			spoke = true
		}
		// this removes the planned conversation which was just said,
		// and any casual conversation that couldn't be said due to restriction.
		next.SetObject(nil)
	}
	return spoke
}
