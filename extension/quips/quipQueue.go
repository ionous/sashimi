package quip

import (
	G "github.com/ionous/sashimi/game"
)

type QuipQueue struct {
	G.IList
}

func SetNextQuip(quip string) NextQuipPhrase {
	return NextQuipPhrase{quip}
}

func SetTheNextQuip(quip G.IObject) NextQuipPhrase {
	return SetNextQuip(quip.Id().String())
}

type NextQuipPhrase struct {
	quip string
}

// SetNextQuip for the associated NPC's next round of conversation.
// FIX: I wonder if this should be merged with UpdateNextQuips() and GetPlayerQuips()
func (p NextQuipPhrase) Execute(g G.Play) {
	quip := g.The(p.quip)
	npc := quip.Object("subject")
	npc.Set("next quip", quip)
	quip.IsNow("planned")
	//TheConversation(g).Queue().QueueQuip(quip)
	g.Log(npc, "set next", quip)
}

// QueueQuip schedules the passed quip to be spoken sometime in the future.
func (q QuipQueue) QueueQuip(quip G.IObject) {
	q.AppendObject(quip)
}

// UpdateNextQuips for all npcs who have a queued quip.
// NOTE: queued conversation will never override what an npc already has to say.
func (q QuipQueue) UpdateNextQuips(qm QuipMemory) {
	// from "slice tricks". this reuses the memory of the quip queue.
	requeue := make([]G.IObject, 0, q.Len())
	// determine what to say next
	for i := 0; i < q.Len(); i++ {
		quip := q.Get(i).Object()
		nextQuip := quip.Get("subject").Object().Get("next quip")
		if nextQuip.Object().Exists() {
			requeue = append(requeue, quip)
		} else {
			// check to make sure this quip wasn't said in the time since it was queued.
			if quip.Is("repeatable") || !qm.Recollects(quip) {
				nextQuip.SetObject(quip)
				quip.IsNow("casual")
			}
		}
	}

	q.Reset()
	for _, quip := range requeue {
		q.AppendObject(quip)
	}
}
