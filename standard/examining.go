package standard

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

var _ = fmt.Sprint

func init() {
	AddScript(func(s *Script) {
		// one visible thing
		// examine studio: You can't see any such thing; sad face.
		s.The("actors",
			Can("examine it").And("examining it").RequiresOne("object"),
			To("examine it", func(g G.Play) { ReflectToTarget(g, "be examined") }),
		)
		// the default action prints the place holder text
		// the events system prints the specifics and prevents the defaults as needed
		// users can override for a particular object the entire thing
		s.The("objects",
			Can("be examined").And("being examined").RequiresOne("actor"),
			To("be examined", func(g G.Play) {
				object := g.The("object")
				object.Go("print details")
				object.Go("print contents")
			}))

		s.The("objects",
			Can("print contents").And("printing contents").RequiresNothing())

		s.The("objects",
			Can("print details").And("printing details").RequiresNothing(),
			To("print details", func(g G.Play) {
				object := g.The("object")
				desc := object.Text("description")
				if desc != "" {
					g.Say(desc)
				} else {
					//g.Say("You see nothing special about:")
					object.Go("print name")
				}
			}))

		s.The("containers",
			When("printing contents").Always(func(g G.Play) {
				if c := g.The("container"); c.Is("open") || c.Is("transparent") {
					ListContents(g, "In the", c)
				}
			}))

		// report an actor examining:
		// where best to do that switch?
		// carry out in inform seems to be limited to the player;....
		///	if the actor is not the player:
		//	say "[The actor] [look] closely at [the noun]." (A).

		s.The("supporters",
			When("printing contents").Always(func(g G.Play) {
				this := g.The("supporter")
				ListContents(g, "On the", this)
			}))

		s.Execute("examine it",
			Matching("examine|x|watch|describe|check {{something}}").
				Or("look|l {{something}}").
				Or("look|l at {{something}}"))
	})
}
