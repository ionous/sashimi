package live

import G "github.com/ionous/sashimi/game"

func Give(prop string) GivePropPhrase {
	return GivePropPhrase{prop: prop}
}

func GiveThe(prop G.IObject) GivePropPhrase {
	return GivePropPhrase{prop: prop.Id().String()}
}

func (give GivePropPhrase) To(actor string) GivingPhrase {
	give.actor = actor
	return GivingPhrase(give)
}

func (give GivingPhrase) Execute(g G.Play) {
	prop, actor := g.The(give.prop), g.The(give.actor)
	//added indirection so we can transform props after the rules of taking/giving have run
	actor.Go("acquire it", prop)
}

type givePhraseData struct {
	prop, actor string
}
type GivePropPhrase givePhraseData
type GivingPhrase givePhraseData
