package standard

import (
	G "github.com/ionous/sashimi/game"
)

// Change("door").To("open")
func Change(target string) ChangePhrase {
	return ChangePhrase{target: target}
}

func (p ChangePhrase) To(state string) G.RuntimePhrase {
	p.state = state
	return p
}

func (p ChangePhrase) Execute(g G.Play) {
	g.The(p.target).IsNow(p.state)
}

type ChangePhrase struct {
	target, state string
}
