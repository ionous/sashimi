package standard

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

var directions []string = []string{"north", "south", "east", "west", "up", "down"}

var Debugging = false

// FIX: like Learn() convert to a game action: actor.Go("move to", dest)
func GoMove(g G.Play) MovePhrase {
	return MovePhrase{g: g}
}

func (move MovePhrase) Object(obj G.IObject) MoveToPhrase {
	move.actor = obj
	return MoveToPhrase(move)
}
func (move MovePhrase) The(obj string) MoveToPhrase {
	return move.Object(move.g.The(obj))
}

func (move MoveToPhrase) To(dest G.IObject) {
	assignTo(move.actor, "whereabouts", dest)
}
func (move MoveToPhrase) ToThe(dest string) {
	move.To(move.g.The(dest))
}

func (move MoveToPhrase) OutOfWorld() {
	assignTo(move.actor, "whereabouts", nil)
}

type MovePhrase struct {
	g     G.Play
	actor G.IObject
}
type MoveToPhrase MovePhrase

func TryMove(actor G.IObject, dir G.IObject, exit G.IObject) {
	if Debugging {
		fmt.Printf("moving %s through %s", dir, exit)
	}
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
			s.The(dir, IsKnownAs(dir[:1]))
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
				exit := from.Object(dir.Text("Name") + "-via")
				if exit.Exists() {
					TryMove(actor, dir, exit)
				} else {
					// try a connected link:
					rev := dir.Object("opposite")
					exit := from.Object(rev.Text("Name") + "-rev-via")
					if exit.Exists() {
						if sources := exit.ObjectList("sources"); len(sources) == 1 {
							TryMove(actor, dir, sources[0])
						}
					} else {
						if Debugging {
							fmt.Printf("couldnt find %s exit", dir)
						}
					}
				}
			}))
		s.The("actors",
			Can("go through it").And("going through it").RequiresOne("door"),
			To("go through it", ReflectToTarget("report pass through")),
		)
		s.The("doors",
			Can("report pass through").And("reporting pass through").RequiresOne("actor"),
			To("report pass through", func(g G.Play) {
				door, actor := g.The("action.Source"), g.The("action.Target")
				if dest := door.Object("destination"); !dest.Exists() {
					if Debugging {
						fmt.Print("couldnt find destination")
					}
				} else if room := dest.Object("whereabouts"); !room.Exists() {
					if Debugging {
						fmt.Print("couldnt find whereabouts")
					}
				} else {
					if Debugging {
						fmt.Print("moving ", actor, " to ", room)
					}
					// FIX: player property change?
					// at the very least a move action.
					GoMove(g).Object(actor).To(room)
					room.Go("report the view")
				}
			}))
		// understandings:
		// FIX: noun matching: so that >go north; >go door. both work.
		// FIX: noun alias: Understand "n" as north.
		s.Execute("go to",
			Matching("go {{something}}"))
		s.Execute("go through it",
			Matching("enter {{something}}"))
	})
}
