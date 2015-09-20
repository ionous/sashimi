package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func init() {
	AddScript(func(s *Script) {
		// one visible thing
		// examine studio: You can't see any such thing; sad face.
		s.The("actors",
			Can("examine it").And("examining it").RequiresOne("object"),
			To("examine it", ReflectToTarget("report examine")),
		)
		// the default action prints the place holder text
		// the events system prints the specifics and prevents the defaults as needed
		// users can override for a particular object the entire thing
		s.The("objects",
			Can("report examine").And("reporting examine").RequiresOne("actor"),
			To("report examine", func(g G.Play) {
				object := g.The("object")
				desc := object.Text("description")
				if desc != "" {
					g.Say(desc)
				} else {
					g.Say("You see nothing special about:")
					object.Go("print name")
				}
			}))

		s.The("containers",
			After("reporting examine").Always(func(g G.Play) {
				this := g.The("action.Source")
				if (this.Is("open") || this.Is("transparent")) && !this.Is("scenery") {
					listContents(g, "In the", this)
				}
			}))

		// report an actor examining:
		// where best to do that switch?
		// carry out in inform seems to be limited to the player;....
		///	if the actor is not the player:
		//	say "[The actor] [look] closely at [the noun]." (A).

		s.The("supporters",
			After("reporting examine").Always(func(g G.Play) {
				this := g.The("action.Source")
				listContents(g, "On the", this)
			}))

		s.Execute("examine it",
			Matching("examine|x|watch|describe|check {{something}}").
				Or("look|l {{something}}").
				Or("look|l at {{something}}"))
	})
}
