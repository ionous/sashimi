package live

import G "github.com/ionous/sashimi/game"

func Put(prop string) PutPhrase {
	return PutPhrase{prop: prop}
}

func (p PutPhrase) Onto(supporter string) PutingPhrase {
	p.supporter = supporter
	return PutingPhrase(p)
}

func (p PutingPhrase) Execute(g G.Play) {
	prop, supporter := g.The(p.prop), g.The(p.supporter)
	AssignTo(prop, "enclosure", supporter)
}

type putData struct {
	prop, supporter string
}

type PutPhrase putData
type PutingPhrase putData
