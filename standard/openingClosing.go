package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func init() {
	AddScript(func(s *Script) {
		s.The("props",
			Called("openers"),
			AreEither("open").Or("closed"),
			AreEither("openable").Or("unopenable"),
			AreEither("locakable").Or("not lockable").Usually("not lockable"),
			AreEither("locked").Or("unlocked").Usually("unlocked"),
		)

		//
		// Open:
		//
		s.The("actors",
			Can("open it").And("opening it").RequiresOne("opener"),
			To("open it", actorTarget("open")),
		)

		// "[regarding the noun][They] [aren't] something [we] [can] open."
		s.The("openers",
			Can("open").And("opening").RequiresOne("actor"),
			To("open", func(g G.Play) {
				this, actor := g.The("opener"), g.The("action.Target")
				if this.Is("openable") {
					if this.Is("open") {
						this.Go("report already open", actor)
					} else {
						this.SetIs("open")
						this.Go("report now open", actor)
					}
				}
			}),
			Can("report already open").And("reporting already opened").RequiresOne("actor"),
			To("report already open", func(g G.Play) {
				g.Say("It's already opened.")
			}),
			Can("report now open").And("reporting now open").RequiresOne("actor"),
			To("report now open", func(g G.Play) {
				this, _ := g.The("opener"), g.The("action.Target")
				g.Say("Now the", this.Text("Name"), "is open.")
				// if the noun doesnt not enclose the actor
				// list the contents of the noun, as a sentence, tersely, not listing concealed items;
				if this.Is("opaque") {
					listContents(g, "In the", this)
				}
			}),
		)

		//
		// Close:
		//
		// one visible thing, and requiring light
		s.The("actors",
			Can("close it").And("closing it").RequiresOne("opener"),
			To("close it", actorTarget("close")),
		)
		s.The("openers",
			Can("close").And("closing").RequiresOne("actor"),
			To("close", func(g G.Play) {
				this, actor := g.The("opener"), g.The("action.Target")
				if this.Is("openable") {
					// FIX: locked?
					if this.Is("closed") {
						this.Go("report already closed", actor)
					} else {
						this.SetIs("closed")
						this.Go("report now closed", actor)
					}
				}
			}),

			Can("report already closed").And("reporting already closed").RequiresOne("actor"),
			To("report already closed", func(g G.Play) {
				g.Say("It's already closed.") //[regarding the noun]?
			}),

			Can("report now closed").And("reporting now closed").RequiresOne("actor"),
			To("report now closed", func(g G.Play) {
				this, _ := g.The("opener"), g.The("action.Target")
				g.Say("Now the", this.Text("Name"), "is closed.")
			}),
		)

		// understandings:
		s.Execute("open it", Matching("open {{something}}"))
		s.Execute("close it", Matching("close {{something}}"))
	})
}
