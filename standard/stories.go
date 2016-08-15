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
		// in lieu of singletons, globals.
		// globals are transmitted to the client in the default view.
		s.The("kinds", Called("globals"), Exist())
		s.The("globals",
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

		endTurn := func(g G.Play) {
			// FIX: we should have a "starting the turn" instead
			// currently, tp get the "frame" counter to see a different number than the last frame
			// due to this frame's player input -- we have to increment both on end turn, and ending the story.
			story := g.The("story")
			turnCount := story.Num("turn count") + 1
			story.SetNum("turn count", turnCount)
			//
			if story.Is("scored") {
				score := story.Num("score")
				status := fmt.Sprintf("%d/%d", int(score), int(turnCount))
				g.The("status bar").SetText("right", status)
			}
		}

		s.The("stories",
			Can("commence").And("commencing").RequiresNothing(),
			Can("end the story").And("ending the story").RequiresNothing(),
			Can("end turn").And("ending the turn").RequiresNothing(),
			Before("commencing").Always(func(g G.Play) {
				inst := g.The("status bar")
				title := g.The("story").Get("name").Text()
				author := g.The("story").Get("author").Text()

				tag := fmt.Sprintf(`"%s" by %s`, title, author)
				inst.Get("left").SetText(title)
				inst.Get("right").SetText(tag)
			}),
			Before("ending the turn").Always(func(g G.Play) {
				story := g.The("story")
				if story.Is("completed") {
					g.StopHere()
				}
			}),
			To("end turn", endTurn),
			After("ending the story").Always(endTurn),
		)

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
					rooms := g.Query("rooms")
					if !rooms.HasNext() {
						panic("story has no rooms")
					}
					room = rooms.Next()
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
