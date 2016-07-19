package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard/live"
	"github.com/ionous/sashimi/util/lang"
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

/// insert it into, receive insertion, being inserted.
func init() {
	AddScript(func(s *Script) {
		// 1. source
		s.The("actors",
			// FIX? word-wise this is wrong ( see tickle-it-with, though it is "correct" )
			Can("insert it into").And("inserting it into").RequiresOne("container").AndOne("prop"),
			To("insert it into", func(g G.Play) { ReflectWithContext(g, "receive insertion") }),
			//  can't insert clothes being worn
			Before("inserting it into").Always(func(g G.Play) {
				prop := g.The("action.Context")
				if worn := prop.Object("wearer"); worn.Exists() {
					g.Say("You can't insert worn clothing.")
					// FIX: try taking off the noun
					g.StopHere()
				}
			}),
			//  can't insert what isn't held
			Before("inserting it into").Always(func(g G.Play) {
				actor, prop := g.The("action.Source"), g.The("action.Context")
				if carrier := Carrier(prop); carrier != actor {
					g.Say("You aren't holding", ArticleName(g, "action.Context", NameFullStop))
					g.StopHere()
				}
			}),
			//  can't insert something into itself
			Before("inserting it into").Always(func(g G.Play) {
				container, prop := g.The("action.Target"), g.The("action.Context")
				if container == prop {
					g.Say("You can't insert something into itself.")
					g.StopHere()
				}
			}),
		)

		// 2. containers
		// FIX FIX FIX could this be "receive"? the more shared events the better,
		// and that would certainly work for any acquisition: actor, supporter,..?
		// the only problem of course is that "be inserted" using specific reporting
		// maybe, rather than chaining events 1->2->3 we could do: 1->{2,3}
		// and check the return result?[the ran default action status of the evt.]
		s.The("containers",
			Can("receive insertion").And("receiving insertion").RequiresOne("prop").AndOne("actor"),
			//  can't insert into closed containers
			Before("receiving insertion").Always(func(g G.Play) {
				container := g.The("container")
				if container.Is("closed") {
					g.Say(lang.Capitalize(DefiniteName(g, "container", nil)), "is closed.")
					g.StopHere()
				}
			}),
			To("receive insertion", func(g G.Play) { ReflectWithContext(g, "be inserted") }))

		// 3. context
		s.The("props",
			Can("be inserted").And("being inserted").RequiresOne("actor").AndOne("container"),
			To("be inserted", func(g G.Play) {
				g.Go(Say("You insert", ArticleName(g, "action.Source", nil), "into", ArticleName(g, "action.Context", NameFullStop)).OnOneLine()).Then(func(g G.Play) {
					g.Go(Insert("action.Source").Into("action.Context"))
				})
			}))
		// input: actor, container, prop
		s.Execute("insert it into",
			Matching("put|insert {{something else}} in|inside|into {{something}}").
				Or("drop {{something else}} in|into|down {{something}}"))
	})
}
