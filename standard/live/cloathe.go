package live

import G "github.com/ionous/sashimi/game"

func Clothe(actor string) ClothePhrase {
	return ClothePhrase{actor: actor}
}

func (p ClothePhrase) With(prop string) WearingPhrase {
	p.clothing = prop
	return WearingPhrase(p)
}

func (p WearingPhrase) Execute(g G.Play) {
	actor, clothing := g.The(p.actor), g.The(p.actor)
	AssignTo(clothing, "wearer", actor)
}

type wearData struct {
	actor, clothing string
}
type ClothePhrase wearData
type WearingPhrase wearData
