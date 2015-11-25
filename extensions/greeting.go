package extensions

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/standard"
)

// GoGreet(g).Introducing(greeter).To(greeted).With(greeting)
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
	if standard.Debugging {
		g.Log("!", "introducing", greeter, "to", greeted, "with", greeting)
	}
	if greeter == g.The("player") && greeted.Exists() {
		con := TheConversation(g)
		interlocutor := con.Interlocutor()
		if npc := interlocutor.Object(); !npc.Exists() {
			interlocutor.SetObject(greeted)
			// it's not necessary to have a greeting if the npc has some latent conversation options.
			if greeting.Exists() {
				// FIX: doesnt raise an error of any sor when we say go("mispelling"
				greeted.Go("comment", greeting)
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
