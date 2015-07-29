package extensions

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/standard"
)

// GoGreet(g).Introducing(greeter).To(greeted).With(greeting)
func GoGreet(g G.Play) GreeterPhrase {
	return GreeterPhrase{g: g}
}

func (greet GreeterPhrase) Introducing(greeter G.IObject) GreetedPhrase {
	greet.greeter = greeter
	return GreetedPhrase(greet)
}

func (greet GreeterPhrase) IntroducingThe(greeter string) GreetedPhrase {
	return greet.Introducing(greet.g.The(greeter))
}

func (greet GreetedPhrase) To(greeted G.IObject) GreetingPhrase {
	greet.greeted = greeted
	return GreetingPhrase(greet)
}

func (greet GreetedPhrase) ToThe(greeted string) GreetingPhrase {
	return greet.To(greet.g.The(greeted))
}

func (greet GreetingPhrase) WithDefault() {
	greeting := greet.greeted.Object("greeting")
	greetActor(greet.g, greet.greeter, greet.greeted, greeting)
}

func (greet GreetingPhrase) With(greeting G.IObject) {
	greetActor(greet.g, greet.greeter, greet.greeted, greeting)
}

func (greet GreetingPhrase) WithQuip(greeting string) {
	greet.With(greet.g.The(greeting))
}

type greetingData struct {
	g                G.Play
	greeter, greeted G.IObject
}
type GreeterPhrase greetingData
type GreetedPhrase greetingData
type GreetingPhrase greetingData

func greetActor(g G.Play, greeter, greeted, greeting G.IObject) {
	if greeter == g.The("player") && greeted.Exists() {
		con := g.Global("conversation").(*Conversation)

		if npc, alreadySpeaking := con.Interlocutor.Get(); !alreadySpeaking {
			if standard.Debugging {
				fmt.Println("!", "Now talking to", greeted, "with", greeting)
			}
			con.Interlocutor.Set(greeted)
			if greeting.Exists() {
				// hrmmm....
				//greeted.Go("discuss", greeting)
				if cmt := greeting.Text("comment"); cmt != "" {
					greeter.Says(cmt)
				}
				con.Queue.SetNextQuip(g, greeting)
				con.History.PushQuip(greeting)
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
