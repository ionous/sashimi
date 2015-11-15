package extensions

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

func (p NextQuipPhrase) Execute(g G.Play) {
	con := TheConversation(g)
	con.Queue().SetNextQuip(g.The(p.quip))
}

// SetNextQuip for the associated NPC's next round of conversation.
// FIX: I wonder if this should be merged with UpdateNextQuips() and GetPlayerQuips()
// rather than a queue -- a pool of next quips -- and it selects the best of the set.
// ( though player is technically from all quips... )
func (q QuipQueue) SetNextQuip(quip G.IObject) {
	npc := quip.Object("subject")
	npc.Set("next quip", quip)
	quip.IsNow("planned")
}

// QueueQuip schedules the passed quip to be spoken sometime in the future.
func (q QuipQueue) QueueQuip(quip G.IObject) {
	q.AppendObject(quip)
}

// UpdateNextQuips for all npcs who have a queued quip.
func (q QuipQueue) UpdateNextQuips(qm QuipMemory) {
	// from "slice tricks". this reuses the memory of the quip queue.
	requeue := make([]G.IObject, 0, q.Len())
	// determine what to say next
	// note: queued conversation will never override what an npc already has to say.
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
