package quips

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
	greet.greeting = string(greeting.Id())
	return GreetingPhrase(greet)
}

func (greet GreetingPhrase) Execute(g G.Play) {
	greeter, greeted := g.The(greet.greeter), g.The(greet.greeted)
	var greeting G.IObject
	if greet.greeting != "" {
		greeting = g.The(greet.greeting)
	}
	if greeting == nil || !greeting.Exists() {
		greeting = greeted.Object("greeting")
		if !greeting.Exists() {
			greeting = g.The("default greeting")
			g.Log("falling back on global default greeting", greet)
		}
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
		greeted.Go("be greeted by", greeter, greeting)
	}
}
