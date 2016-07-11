package live

import G "github.com/ionous/sashimi/game"

func Insert(prop string) InsertPhrase {
	return InsertPhrase{prop: prop}
}

func (p InsertPhrase) Into(container string) InsertingPhrase {
	p.container = container
	return InsertingPhrase(p)
}

func (p InsertingPhrase) Execute(g G.Play) {
	prop, container := g.The(p.prop), g.The(p.container)
	AssignTo(prop, "enclosure", container)
}

type insertData struct {
	prop, container string
}

type InsertPhrase insertData
type InsertingPhrase insertData
