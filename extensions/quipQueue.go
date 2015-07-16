package extensions

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
)

type QuipQueue struct {
	quips []G.IObject
}

// QueueQuip schedules the passed quip to be spoken in the future.
func (q *QuipQueue) QueueQuip(quip G.IObject) {
	q.quips = append(q.quips, quip)
}

// ResetQuipQueue removes all pending conversation.
// For testing's sake, returns the number of quips which were pending.
func (q *QuipQueue) ResetQuipQueue() {
	q.quips = nil
}

func (q *QuipQueue) Len() int {
	return len(q.quips)
}

// SetNextQuips updates the next quip for all npcs who have a queued quip.
func (q *QuipQueue) UpdateNextQuips(g G.Play, qm QuipMemory) {
	if Debugging {
		fmt.Println(fmt.Sprintf("! updating %d quips", len(q.quips)))
	}
	// from "slice tricks". this reuses the memory of the quip queue.
	requeue := q.quips[:0]
	// determine what to say next
	// note: queued conversation will never override what an npc already has to say.
	for _, quip := range q.quips {
		npc := quip.Object("subject")
		if npc.Object("next quip").Exists() {
			requeue = append(requeue, quip)
		} else {
			// check to make sure this quip wasn't said in the time since it was queued.
			if quip.Is("repeatable") || !qm.Recollects(quip) {
				nextQuip := g.Add("next quip")
				npc.Set("next quip", nextQuip)
				nextQuip.Set("quip", quip)
				nextQuip.SetIs("casual")
			}
		}
	}
	q.quips = requeue
}
