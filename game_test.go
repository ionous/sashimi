package sashimi

import (
	"github.com/ionous/sashimi/console"
	G "github.com/ionous/sashimi/game"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard"
	x "github.com/smartystreets/goconvey/convey"
	"testing"
)

//
func TestStandardRules(t *testing.T) {
	x.Convey("Given the script", t, func() {
		//s := standard.StandardRules()
		s := InitScripts()
		_, _, err := CompileGameWithConsole(s, console.NewConsole(), true)
		x.So(err, x.ShouldBeNil)
	})
}

//
func TestObjectSet(t *testing.T) {
	x.Convey("Given the script", t, func() {
		s := InitScripts()
		s.The("kinds",
			Have("amBlank", "text"),
			Have("amSet", "text"))

		s.The("kind",
			Called("test"),
			Has("amSet", "original"))

		game, model, err := CompileGameWithConsole(s, console.NewConsole(), true)
		x.So(err, x.ShouldBeNil)
		x.So(model, x.ShouldNotBeNil)

		inst, err := model.Instances.FindInstance("test")
		x.So(inst, x.ShouldNotBeNil)

		x.Convey("the test object", func() {
			gobj, exists := game.Objects[inst.Id()]
			x.So(exists, x.ShouldBeTrue)
			x.So(gobj, x.ShouldNotBeNil)

			obj := R.NewObjectAdapter(game, gobj)

			x.Convey("should have a value, and should be settable", func() {
				obj.Text("amSet")
				x.So(obj.Text("amSet"), x.ShouldEqual, "original")

				obj.SetText("amSet", "new")
				x.So(obj.Text("amSet"), x.ShouldEqual, "new")
			})

			x.Convey("should *not* have a value, and should be settable", func() {
				x.So(obj.Text("amBlank"), x.ShouldEqual, "")

				obj.SetText("amBlank", "not blank any more")
				x.So(obj.Text("amBlank"), x.ShouldEqual, "not blank any more")
			})
		})
	})
}

//
func TestStartupText(t *testing.T) {
	x.Convey("TestStartupText", t, func() {
		s := InitScripts()
		c := console.NewBufCon(nil)

		s.The("story",
			Called("Testing"),
			Has("author", "me"),
			Has("headline", "extra extra"))

		s.The("room",
			Called("Somewhere"),
			Has("description", "an empty room"),
			When("describing").Always(func(g G.Play) {
				g.StopHere()
			}),
		)

		game, model, err := CompileGameWithConsole(s, c, true)
		x.So(err, x.ShouldBeNil)

		story := game.FindFirstOf(model.Classes.FindClass("stories"))
		x.So(story, x.ShouldNotBeNil)

		room := game.FindFirstOf(model.Classes.FindClass("rooms"))
		x.So(room, x.ShouldNotBeNil)

		err = game.SendEvent("starting to play", story.String())
		x.So(err, x.ShouldBeNil)

		game.RunForever()
		x.So(c.Results(), x.ShouldResemble, []string{
			"", // FIX: this line shouldnt exist
			"Testing",
			"extra extra by me",
			standard.VersionString,
			"",
			"Somewhere",
			"an empty room",
		})
	})
}
