package standard

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

//
// System actions
func init() {
	AddScript(func(s *Script) {
		s.The("kinds",
			Called("stories").WithSingularName("story"),
			Have("author", "text"),
			Have("headline", "text"),
			AreEither("scored").Or("unscored").Usually("unscored"),
			// Inform uses global variables, which would be much nicer.
			// ex. The maximum score is 1.
			Have("score", "num"),
			Have("maximum score", "num"),
			Have("turn count", "num"),
			AreOneOf("playing", "completed", "starting").Usually("starting"),
		)

		s.The("stories",
			Can("commence").And("commencing").RequiresNothing(),
			Can("end the story").And("ending the story").RequiresNothing(),
			Can("end turn").And("ending the turn").RequiresNothing(),
			Before("ending the turn").Always(func(g G.Play) {
				story := g.The("story")
				if story.Is("completed") {
					g.StopHere()
				}
			}),
			To("end turn", func(g G.Play) {
				// almost feel like we should have a "starting the turn" instead
				story := g.The("story")
				turnCount := story.Num("turn count") + 1
				story.SetNum("turn count", turnCount)
				//
				if story.Is("scored") {
					score := story.Num("score")
					status := fmt.Sprintf("%d/%d", int(score), int(turnCount))
					g.The("status bar").SetText("right", status)
				}
			}))

		s.The("stories",
			To("commence", func(g G.Play) {
				// FIX: duplication with end turn
				story := g.The("story")
				if story.Is("scored") {
					score := story.Num("score")
					status := fmt.Sprintf("%d/%d", int(score), int(0))
					g.The("status bar").SetText("right", status)
				}
				room := g.The("player").Object("whereabouts")
				if !room.Exists() {
					rooms := g.List("rooms")
					if rooms.Len() == 0 {
						panic("story has no rooms")
					}
					room = rooms.Get(0).Object()
				}
				story.Go("set initial position", g.The("player"), room).Then(func(g G.Play) {
					story.Go("print the banner").Then(func(g G.Play) {
						room = g.The("player").Object("whereabouts")
						// FIX: Go() should handle both Name() and ref
						story.Go("describe the first room", room).Then(func(g G.Play) {
							story.IsNow("playing")
						})
					})
				})
			}))

		s.The("stories",
			Have("player input", "text"),
			Can("parse player input").And("parsing player input").RequiresNothing())

		s.The("stories",
			To("end the story", func(g G.Play) {
				story := g.The("story")
				g.Say("*** The End ***")
				story.IsNow("completed")

				if story.Is("scored") {
					score, maxScore, turnCount := story.Num("score"), story.Num("maximum score"), story.Num("turn count")
					g.Say(fmt.Sprintf("In that game you scored %d out of a possible %d, in %d turns.",
						int(score), int(maxScore), int(turnCount)+1,
					))
				}
			}))

		s.The("stories",
			Can("set initial position").
				And("setting initial position").
				RequiresOne("actor").
				AndOne("room"),
			To("set initial position", func(g G.Play) {
				player := g.The("action.Target")
				room := g.The("action.Context")
				player.Set("whereabouts", room) // Now("player's whereabouts is $room")
			}))

		s.The("stories",
			Can("describe the first room").
				And("describing the first room").RequiresOne("room"),
			To("describe the first room", func(g G.Play) {
				room := g.The("action.Target")
				room.Go("report the view")
			}),
		)
	})
}
