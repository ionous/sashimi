package standard

import (
	"bitbucket.org/pkg/inflect"
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
				text = inflect.Capitalize(text)
				g.Say(text)
				g.StopHere()
			}))

		s.The("props", Can("report inoperable").And("reporting inoperable").RequiresNothing(),
			To("report inoperable", func(g G.Play) {
				g.Say("It's inoperable.")
			}))

		//
		// Turn on:
		//
		s.The("actors",
			Can("switch it on").And("switching it on").RequiresOne("prop"),
			To("switch it on", ReflectToTarget("report switched on")))

		s.The("devices",
			Can("report switched on").And("reporting switched on").RequiresOne("actor"),
			To("report switched on", func(g G.Play) {
				device, actor := g.The("device"), g.The("actor")
				if device.Is("inoperable") {
					device.Go("report inoperable")
				} else {
					if device.Is("switched on") {
						device.Go("report already on", actor)
					} else {
						device.IsNow("switched on")
						device.Go("report now on", actor)
					}
				}
			}),
			Can("report already on").And("reporting already on").RequiresOne("actor"),
			To("report already on", func(g G.Play) {
				g.Say("It's already switched on.")
			}),
			Can("report now on").And("reporting now on").RequiresOne("actor"),
			To("report now on", func(g G.Play) {
				device := g.The("device")
				g.Say("Now the", device.Text("Name"), "is on.")
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
				device, actor := g.The("device"), g.The("actor")
				if device.Is("switched off") {
					device.Go("report already off", actor)
				} else {
					device.IsNow("switched off")
					device.Go("report now off", actor)
				}
			}),
			Can("report already off").And("reporting already off").RequiresOne("actor"),
			To("report already off", func(g G.Play) {
				g.Say("It's already off.") //[regarding the noun]?
			}),
			Can("report now off").And("reporting now off").RequiresOne("actor"),
			To("report now off", func(g G.Play) {
				device, _ := g.The("device"), g.The("actor")
				g.Say("Now the", device.Text("Name"), "is off.")
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
