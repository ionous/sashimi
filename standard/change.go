package standard

import (
	G "github.com/ionous/sashimi/game"
)

// Change("door").To("open")
func Change(target string) ChangePhrase {
	return ChangePhrase{target: target}
}

func (p ChangePhrase) To(state string) ChangePhrase {
	return p.And(state)
}

func (p ChangePhrase) And(state string) ChangePhrase {
	p.states = append(p.states, state)
	return p
}

func (p ChangePhrase) Execute(g G.Play) {
	tgt := g.The(p.target)
	for _, state := range p.states {
		tgt.IsNow(state)
	}
}

type ChangePhrase struct {
	target string
	states []string
}
