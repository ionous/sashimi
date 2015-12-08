package native

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

func (m Memory) TriesToLearn(quip G.IObject) (newlyLearned bool) {
	// return newlyLearned
	if recollects := quip.Is("recollected"); !recollects {
		quip.IsNow("recollected")
		newlyLearned = true
	}
	return
}

// LearnQuip causes actors to recollect the passed quip.
// we can use this for facts for now too
// mostly the player will need this -- so just a table with precese is enough
// but it could also be actor, id
func (m Memory) Learns(quip G.IObject) {
	//m.AppendObject(quip)
	quip.IsNow("recollected")
}

// RecollectsQuip determines if the passed quip has been spoken.
func (m Memory) Recollects(quip G.IObject) bool {
	//return m.Contains(quip)
	return quip.Is("recollected")
}
