package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	"log"
	"strings"
)

var directions []string = []string{"north", "south", "east", "west", "up", "down"}

func Move(actor G.IObject, dest G.IObject) (changed bool) {
	if ok := actor.Object("whereabouts") != dest; ok {
		Assign(actor, "whereabouts", dest)
		changed = true
	}
	return changed
}

func TryMove(actor G.IObject, dir G.IObject, exit G.IObject) {
	log.Printf("moving %s through %s", dir, exit)
	actor.Go("go through it", exit)
}

func init() {
	AddScript(func(s *Script) {
		// 1. A Room (contains) Doors
		s.The("openers",
			Called("doors"),
			Exist())

		// 2. An Exit (has a matching) Entrance
		s.The("doors",
			// exiting using a door leads to one location
			HaveOne("destination", "door").
				// one door can be the target of many other doors
				Implying("doors", HaveMany("sources", "doors")),
		)

		// 3. A Room+Travel Direction (has a matching) Exit
		// FIX: without relation by value we have to give each room a set of explict directions
		// each direction relation points to the matching door
		for _, dir := range directions {
			// moving in a direction, takes us to a room's entrance.
			s.The("rooms", HaveOne(dir+"-via", "door").
				// FIX: opposite relation shouldnt be required
				Implying("doors", HaveMany("x-via-"+dir, "rooms")))
			// the reverse direction is needed because we dont all of the directions at compile time
			// ( we have the default set, but users could add more )
			s.The("rooms", HaveOne(dir+"-rev-via", "door").
				Implying("doors", HaveMany("x-rev-via-"+dir, "rooms")))
		}

		// Directions:
		s.The("kinds", Called("directions"),
			HaveOne("opposite", "direction").
				//FIX: the reverse shouldnt be required.
				Implying("directions", HaveOne("x-opposite", "direction")),
		)

		for i := 0; i < len(directions)/2; i++ {
			a, b := directions[2*i], directions[2*i+1]
			s.The("direction", Called(a), Has("opposite", b))
			s.The("direction", Called(b), Has("opposite", a))
		}

		// FIX: change logs to reports
		s.The("actors",
			Can("go to").And("going to").RequiresOne("direction"),
			To("go to", func(g G.Play) {
				actor, dir := g.The("actor"), g.The("action.Target")
				from := actor.Object("whereabouts")
				// try the forward direction:
				exit := from.Object(dir.Name() + "-via")
				if exit.Exists() {
					TryMove(actor, dir, exit)
				} else {
					// try a connected link:
					rev := dir.Object("opposite")
					exit := from.Object(rev.Name() + "-rev-via")
					if exit.Exists() {
						if sources := exit.ObjectList("sources"); len(sources) == 1 {
							TryMove(actor, dir, sources[0])
						}
					} else {
						log.Printf("couldnt find %s exit", dir)
					}
				}
			}))
		s.The("actors",
			Can("go through it").And("going through it").RequiresOne("door"),
			To("go through it", actorTarget("pass through")),
		)
		s.The("doors",
			Can("pass through").And("passing through").RequiresOne("actor"),
			To("pass through", func(g G.Play) {
				door, actor := g.The("action.Source"), g.The("action.Target")
				if dest := door.Object("destination"); !dest.Exists() {
					log.Print("couldnt find destination")
				} else if room := dest.Object("whereabouts"); !room.Exists() {
					log.Print("couldnt find whereabouts")
				} else {
					log.Print("moving ", actor, " to ", room)
					if Move(actor, room) {
						// FIX: duplicated in stories describe the first room
						room.Go("report the view")
						room.SetIs("visited")
						g.The("status bar").SetText("left", strings.Title(room.Name()))
					}
				}
			}))
		// understandings:
		// FIX: noun matching: so that >go north; >go door. both work.
		// FIX: noun aliass: Understand "n" as north.
		s.Execute("go to",
			Matching("go {{something}}"))
		s.Execute("go through it",
			Matching("enter {{something}}"))
	})
}