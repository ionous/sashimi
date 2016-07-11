package quips

import (
	G "github.com/ionous/sashimi/game"
)

type Conversation struct {
	G.IObject
}

// Converse provides an object-oriented wrapper around the "conversation" sashimi type.
func Converse(g G.Play) Conversation {
	return Conversation{g.The("conversation")}
}
func (c Conversation) Actor() G.IValue {
	return c.Get("actor")
}
func (c Conversation) Topic() G.IValue {
	return c.Get("topic")
}
func (c Conversation) Quip() G.IValue {
	return c.Get("current")
}
func (c Conversation) History() QuipHistory {
	return QuipHistory{c.Get("current"), c.Get("parent"), c.Get("grandparent")}
}
func (c Conversation) Conversing() bool {
	// note: we dont always have a quip ( ex. greeting can be null )
	return c.Actor().Object().Exists()
}
func (c Conversation) Reset() {
	c.Actor().SetObject(nil)
	c.History().Reset()
}
