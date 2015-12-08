package native

import (
	G "github.com/ionous/sashimi/game"
)

// Introduce with out invoking events.
// ex. g.Go(Introdue(greeter).To(greeted).With(greeting))
// for all the silly things go complains about
// silently discarding return values is somehow perfectly fine. sigh.
func Introduce(greeter string) GreeterPhrase {
	return GreeterPhrase{greeter: greeter}
}

func (greet GreeterPhrase) To(greeted string) GreetedPhrase {
	greet.greeted = greeted
	return GreetedPhrase(greet)
}

func (greet GreetedPhrase) WithDefault() GreetingPhrase {
	return greet.With("")
}

func (greet GreetedPhrase) With(greeting string) GreetingPhrase {
	greet.greeting = greeting
	return GreetingPhrase(greet)
}

func (greet GreetedPhrase) WithQuip(greeting G.IObject) GreetingPhrase {
	greet.greeting = greeting.Id().String()
	return GreetingPhrase(greet)
}

func (greet GreetingPhrase) Execute(g G.Play) {
	greeter, greeted := g.The(greet.greeter), g.The(greet.greeted)
	var greeting G.IObject
	if greet.greeting == "" {
		greeting = greeted.Object("greeting")
	} else {
		greeting = g.The(greet.greeting)
	}
	greetActor(g, greeter, greeted, greeting)
}

type greetingData struct {
	greeter, greeted, greeting string
}
type GreeterPhrase greetingData
type GreetedPhrase greetingData
type GreetingPhrase greetingData

func greetActor(g G.Play, greeter, greeted, greeting G.IObject) {
	g.Log(greeter, "greeting", greeted, "with", greeting)
	if greeter == g.The("player") && greeted.Exists() {
		c := Converse(g)
		if npc := c.Actor().Object(); !npc.Exists() {
			c.Actor().SetObject(greeted)
			// it's not necessary to have a greeting if the npc has some latent conversation options.
			if greeting.Exists() {
				// FIX: doesnt raise an error of any sor when we say go("mispelling"
				greeter.Go("comment", greeting)
			}
		} else {
			if npc == greeted {
				g.Say("You're already speaking to them!")
			} else {
				g.Say("You're already speaking to someone!")
			}
		}
	}
}
