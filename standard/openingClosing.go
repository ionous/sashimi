package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func init() {
	// (note: the action uses props, so that stories can make any prop behave similar to an opener. )
	AddScript(func(s *Script) {
		s.The("props",
			Called("openers"),
			AreEither("open").Or("closed"),
			// note: unopenable sounds like it cant become open, rather than it cannot be opened (by someone).
			AreEither("openable").Or("not openable"),
			AreEither("locakable").Or("not lockable").Usually("not lockable"),
			AreEither("unlocked").Or("locked"),
		)

		//
		// Open:
		//
		s.The("actors",
			Can("open it").And("opening it").RequiresOne("prop"),
			To("open it", ReflectToTarget("report open")),
		)

		// "[regarding the noun][They] [aren't] something [we] [can] open."
		s.The("props",
			Can("report open").And("reporting open").RequiresOne("actor"),
			To("report open", func(g G.Play) {
				prop, actor := g.The("prop"), g.The("actor")
				if !prop.Is("openable") {
					prop.Go("report not openable", actor)
				} else {
					if prop.Is("locked") {
						prop.Go("report locked", actor)
					} else {
						if prop.Is("open") {
							prop.Go("report already open", actor)
						} else {
							prop.IsNow("open")
							prop.Go("report now open", actor)
						}
					}
				}
			}),
			Can("report locked").And("reporting locked").RequiresOne("actor"),
			To("report locked", func(g G.Play) {
				g.The("actor").Says("It's locked!")
			}),
			Can("report not openable").And("reporting not openable").RequiresOne("actor"),
			To("report not openable", func(g G.Play) {
				g.Say("That's not something you can open.")
			}),
			Can("report already open").And("reporting already opened").RequiresOne("actor"),
			To("report already open", func(g G.Play) {
				g.Say("It's already opened.")
			}),
			Can("report now open").And("reporting now open").RequiresOne("actor"),
			To("report now open", func(g G.Play) {
				opener := g.The("opener")
				g.Say("The", opener.Text("Name"), "is now open.")
				// if the noun doesnt not enclose the actor
				// list the contents of the noun, as a sentence, tersely, not listing concealed items;
				// FIX? not all openers are opaque/transparent, and not all openers have contents.
				if opener.Is("opaque") {
					listContents(g, "In the", opener)
				}
			}),
		)

		//
		// Close:
		//
		// one visible thing, and requiring light
		s.The("actors",
			Can("close it").And("closing it").RequiresOne("prop"),
			To("close it", ReflectToTarget("report close")),
		)
		s.The("props",
			Can("report close").And("report closing").RequiresOne("actor"),
			To("report close", func(g G.Play) {
				prop, actor := g.The("prop"), g.The("actor")
				if !prop.Is("openable") {
					prop.Go("report not closeable", actor)
				} else {
					// FIX: locked?
					if prop.Is("closed") {
						prop.Go("report already closed", actor)
					} else {
						prop.IsNow("closed")
						prop.Go("report now closed", actor)
					}
				}
			}),
			Can("report not closeable").And("reporting not closeable").RequiresOne("actor"),
			To("report not closeable", func(g G.Play) {
				g.Say("That's not something you can close.")
			}),
			Can("report already closed").And("reporting already closed").RequiresOne("actor"),
			To("report already closed", func(g G.Play) {
				g.Say("It's already closed.") //[regarding the noun]?
			}),
			Can("report now closed").And("reporting now closed").RequiresOne("actor"),
			To("report now closed", func(g G.Play) {
				prop := g.The("prop")
				g.Say("Now the", prop.Text("Name"), "is closed.")
			}),
		)

		// understandings:
		s.Execute("open it", Matching("open {{something}}"))
		s.Execute("close it", Matching("close {{something}}"))
	})
}
