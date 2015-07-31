package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func Insert(prop string) InsertPhrase {
	return InsertPhrase{prop: prop}
}

func (p InsertPhrase) Into(container string) InsertingPhrase {
	p.container = container
	return InsertingPhrase(p)
}

func (p InsertingPhrase) Execute(g G.Play) {
	prop, container := g.The(p.prop), g.The(p.container)
	assignTo(prop, "enclosure", container)
}

type insertData struct {
	prop, container string
}

type InsertPhrase insertData
type InsertingPhrase insertData

// from inform:
// 	"insert applies into two things", doesnt use preferably held. [ possibly a mistake? ]
//  "convert insert into drop where possible" ( if the second noun is down, or if the actor is in the second noun )
// 		?FIX? i dont understand.
//  "can't insert what isn't held"
// 		? again, isnt that point of "something preferably held"?
//  "can't insert something into itself"
//  "can't insert into closed containers"
// 	X "can't insert into what's not a container": implicit in this definition
//  "can't insert clothes being worn"
//  FIX: can't insert if this exceeds carrying capacity
// * "carry out inserting": now in second noun
// * "concise report inserting rule": "done"
// * "standard reporting rule": "actor put thing into thing"
func init() {
	AddScript(func(s *Script) {
		// 1. source
		s.The("actors",
			// FIX? word-wise this is wrong ( see tickle-it-with, though it is "correct" )
			Can("insert it into").And("inserting it into").RequiresOne("container").AndOne("prop"),
			To("insert it into", ReflectWithContext("report insert")),
			//  can't insert clothes being worn
			WhenCapturing("inserting it into", func(g G.Play) {
				prop := g.The("action.Context")
				if worn := prop.Object("wearer"); worn.Exists() {
					g.Say("You can't insert worn clothing.")
					// FIX: try taking off the noun
					g.StopHere()
				}
			}),
			//  can't insert what isn't held
			WhenCapturing("inserting it into", func(g G.Play) {
				actor, prop := g.The("action.Source"), g.The("action.Context")
				if carrier, ok := Carrier(prop); !ok || carrier != actor {
					g.Say("You aren't holding", ArticleName(g, "action.Context", NameFullStop))
					g.StopHere()
				}
			}),
			//  can't insert something into itself
			WhenCapturing("inserting it into", func(g G.Play) {
				container, prop := g.The("action.Target"), g.The("action.Context")
				if container == prop {
					g.Say("You can't insert something into itself.")
					g.StopHere()
				}
			}),
			//  can't insert into closed containers
			WhenCapturing("inserting it into", func(g G.Play) {
				container := g.The("action.Target")
				if container.Is("closed") {
					g.Say(ArticleName(g, "action.Target", nil), "is closed.")
					g.StopHere()
				}
			}),
		)
		// 2. containers
		s.The("containers",
			Can("report insert").And("reporting insert").RequiresOne("prop").AndOne("actor"),
			To("report insert", ReflectWithContext("report insertion")))
		// 3. context
		s.The("props",
			Can("report insertion").And("reporting insertion").RequiresOne("actor").AndOne("container"),
			To("report insertion", func(g G.Play) {
				g.Go(Insert("action.Source").Into("action.Context"))
				g.Say("You insert", ArticleName(g, "action.Source", nil), "into", ArticleName(g, "action.Context", NameFullStop))
			}))
		// input: actor, container, prop
		s.Execute("insert it into",
			Matching("put|insert {{something else}} in|inside|into {{something}}").
				Or("drop {{something else}} in|into|down {{something}}"))
	})
}
