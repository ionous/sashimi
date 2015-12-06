package quip

import (
	G "github.com/ionous/sashimi/game"
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
		g.Log("conversing with", npc)

		// sometimes conversations want to loop, without the player saying anything
		for again := true; again; {
			again = false
			// CHANGE; in case the npc replies to itself, re-update the queu.
			// (ex. the vending machine says "you seem perfectly safe to me" and then re-greets)
			currentQuip := c.History().MostRecent()
			currentRestricts := currentQuip.Exists() && currentQuip.Is("restrictive")

			// handle queued conversation, unless the current quip is restrictive.
			// restrictive is supposed to limit responses to a specific sub-set
			// ( to those who "directly follow" )
			// GetPlayerQuips() figures that out.
			if currentRestricts {
				g.Log("quip", currentQuip, "restricts player response; not updating next quips")
			} else {
				g.Log("updating next quips")
				c.Queue().UpdateNextQuips(PlayerMemory(g))
			}

			// process the current speaker first:
			if c.perform(npc, currentRestricts) {
				again = true
			}

			// process anyone else who might have something to say:
			for i, actors := 0, g.List("actors"); i < actors.Len(); i++ {
				actor := actors.Get(i).Object()
				// FIX? repeat with **visible** people who are not the player....
				if actor != npc {
					if c.perform(actor, currentRestricts) {
						again = true
					}
				}
			}
		}

		// we might have changed conversations...
		// note: if the player has nothing to say, conversation ends.
		if npc := interlocutor.Object(); npc.Exists() {
			player := g.The("player")
			player.Go("print conversation choices", npc)
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
