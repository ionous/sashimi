package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func init() {
	AddScript(func(s *Script) {
		// "Represents a machine or contrivance of some kind which can be switched on or off."
		s.The("props",
			Called("devices"),
			AreEither("operable").Or("inoperable"),
			AreEither("switched off").Or("switched on"))

		s.The("devices",
			When("printing name text").
				Always(func(g G.Play) {
				text := ArticleName(g, "device", func(device G.IObject) (status string) {
					if device.Is("operable") {
						if device.Is("switched on") {
							status = "switched on"
						} else {
							status = "switched off"
						}
					}
					return status
				})
				g.Say(text)
				g.StopHere()
			}))

		//
		// Turn on:
		//
		s.The("actors",
			Can("switch it on").And("switching it on").RequiresOne("prop"),
			To("switch it on", ReflectToTarget("report switch on")))

		s.The("devices",
			Can("report switch on").And("reporting switch on").RequiresNothing(),
			To("report switch on", func(g G.Play) {
				prop, actor := g.The("prop"), g.The("actor")
				if prop.Is("inoperable") {
					if prop.Is("switched on") {
						prop.Go("report already on", actor)
					} else {
						prop.IsNow("switched on")
						prop.Go("report now on", actor)
					}
				}
			}),
			Can("report already on").And("reporting already on").RequiresOne("actor"),
			To("report already on", func(g G.Play) {
				g.Say("It's already switched on.")
			}),
			Can("report now on").And("reporting now on").RequiresOne("actor"),
			To("report now on", func(g G.Play) {
				prop := g.The("owner")
				g.Say("Now the", prop.Text("Name"), "is on.")
			}))

		//
		// Turn off
		//
		s.The("actors",
			Can("switch it off").And("switching it off").RequiresOne("prop"),
			To("switch it off", ReflectToTarget("report switch off")))

		s.The("devices",
			Can("report switch off").And("reporting switch off").RequiresNothing(),
			To("report switch off", func(g G.Play) {
				prop, actor := g.The("prop"), g.The("actor")
				if prop.Is("switched off") {
					prop.Go("report already off", actor)
				} else {
					prop.IsNow("switched off")
					prop.Go("report now off", actor)
				}
			}),
			Can("report already off").And("reporting already off").RequiresOne("actor"),
			To("report already off", func(g G.Play) {
				g.Say("It's already off.") //[regarding the noun]?
			}),
			Can("report now off").And("reporting now off").RequiresOne("actor"),
			To("report now off", func(g G.Play) {
				prop, _ := g.The("owner"), g.The("actor")
				g.Say("Now the", prop.Text("Name"), "is off.")
			}))

		// understandings:
		// note: inform has "template matching" here --
		// "switch [something switched on]" as switching off.
		// FIX:  inform's  "understand" has many meanings, but i think itd be better here
		// maybe: s.Understand.Or.As; Understand().WhenMatching("").Or()
		s.Execute("switch it on",
			Matching("switch|turn on {{something}}").
				Or("switch {{something}} on"))

		s.Execute("switch it off",
			Matching("turn|switch off {{something}}").
				Or("turn|switch {{something}} off"))

	})
}
