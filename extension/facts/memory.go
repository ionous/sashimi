package facts

import (
	G "github.com/ionous/sashimi/game"
)

type Memory struct {
	g G.Play
}

func PlayerRecollects(g G.Play, fact string) bool {
	return PlayerMemory(g).Recollects(g.The(fact))
}

func PlayerMemory(g G.Play) Memory {
	return Memory{g}
}

func (m Memory) TriesToLearn(fact G.IObject) (newlyLearned bool) {
	// return newlyLearned
	if recollects := fact.Is("recollected"); !recollects {
		fact.IsNow("recollected")
		newlyLearned = true
	}
	return
}

// LearnQuip causes actors to recollect the passed fact.
// we can use this for facts for now too
// mostly the player will need this -- so just a table with precese is enough
// but it could also be actor, id
func (m Memory) Learns(fact G.IObject) {
	//m.AppendObject(fact)
	fact.IsNow("recollected")
}

// RecollectsQuip determines if the passed fact has been spoken.
func (m Memory) Recollects(fact G.IObject) bool {
	//return m.Contains(fact)
	return fact.Is("recollected")
}
