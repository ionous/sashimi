package standard

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

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
					g.Say(fmt.Sprintf("You aren't holding %s.", prop.Text("Name")))
					g.StopHere()
				}
			}),
			//  can't insert something into itself
			WhenCapturing("inserting it into", func(g G.Play) {
				container, prop := g.The("action.Target"), g.The("action.Context")
				if container == prop {
					g.Say("can't insert something into itself.")
					g.StopHere()
				}
			}),
			//  can't insert into closed containers
			WhenCapturing("inserting it into", func(g G.Play) {
				container := g.The("action.Target")
				if container.Is("closed") {
					g.Say(fmt.Sprintf("%s is closed.", container.Text("Name")))
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
				prop, container := g.The("action.Source"), g.The("action.Context")
				Assign(prop, "enclosure", container)
				g.Say("You insert {{action.Source.Name}} into {{action.Context.Name}}.")
			}))
		// input: actor, container, prop
		s.Execute("insert it into",
			Matching("put|insert {{something else}} in|inside|into {{something}}").
				Or("drop {{something else}} in|into|down {{something}}"))
	})
}
