package native

import (
	G "github.com/ionous/sashimi/game"
)

func ChangeTopic(topic string) TopicPhrase {
	return TopicPhrase{topic}
}

type TopicPhrase struct {
	topic string
}

func (t TopicPhrase) Execute(g G.Play) {
	con := Converse(g)
	con.Topic().SetObject(g.The(t.topic))
}
