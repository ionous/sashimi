package facts

import (
	G "github.com/ionous/sashimi/game"
)

// FIX: replace  with player, go learn
// ALSO: if this were in the "fact" package, it could be: fact.Learn
// and maybe prop.Give?
func Learn(fact string) FactPhrase {
	return FactPhrase{fact}
}
func LearnThe(fact G.IObject) FactPhrase {
	return FactPhrase{string(fact.Id())}
}

type FactPhrase struct {
	fact string
}

func (p FactPhrase) Execute(g G.Play) {
	PlayerMemory(g).Learns(g.The(p.fact))
}

func PlayerLearns(g G.Play, fact string) (newlyLearned bool) {
	return PlayerMemory(g).TriesToLearn(g.The(fact))
}
