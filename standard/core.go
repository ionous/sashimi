package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func InitStandardLibrary() *Script {
	return InitScripts()
}

// FIX: there's no error testing here and its definitely possible to screw things up.
func Assign(dest G.IObject, rel string, prop G.IObject) {
	// sure hope there's no errors, would relation by value remove the need for transaction?
	if _, parentRel := DirectParent(prop); parentRel != "" {
		// note: an object like the fishFood isnt "in the world", and doesnt have an owner field to clear.
		prop.SetObject(parentRel, nil)
	}
	prop.SetObject(rel, dest)
}

func Give(actor G.IObject, prop G.IObject) {
	Assign(actor, "owner", prop)
}

//
func init() {
	AddScript(func(s *Script) {
		// FIX: hierarchy is a work in progress.
		// kinds> stories, rooms, objects > actors (> animals), doors, props(> containers, supporters, devices)

		// FIX: IObject and GameObject, ObjectList are confusing versus kind/object
		// Instance might be better...? or some better name for the class of objects
		// i tried "subjects" but hated it more.

		// vs. descriptions in "kind"
		// it seems to make sense for now to have two separate description fields.
		// rooms like to say their description, while objects like to say their brief initial appearence ( or name, if there's none. )
		s.The("rooms",
			Have("description", "text"))

		s.The("objects",
			Have("description", "text"),
			Have("brief", "text"))

		// inform's rooms: lighted, dark; unvisited, visited; description, region
		s.The("kinds",
			Called("rooms"),
			AreEither("visited").Or("unvisited").Usually("unvisited"),
		)

		// the class hierarchy means rooms cant contain other rooms.
		s.The("rooms",
			Have("contents", "objects").
				Implying("objects", HaveOne("whereabouts", "room")))

		// things		unlit, lit
		// 	inedible, edible
		//
		// 	unmarked for listing, marked for listing
		// 	described, undescribed : i think, whether to appear in any room descriptions
		// 	mentioned, unmentioned : i think, whether it has appeared in a room description
		// bool	scenery
		// 	wearable
		// 	pushable between rooms
		// 	.handled
		// 	.description (in objects and rooms)
		// 	.initial appearance (brief)
		// 	matching key
		s.The("kinds",
			Called("objects"),
			Exist())

		// hrmmm.... are actors really scenery? handled?
		s.The("objects",
			AreEither("unhandled").Or("handled"),
			AreEither("scenery").Or("not scenery").Usually("not scenery"))

		s.The("objects",
			Called("doors"),
			Exist())

		// nothing special: just a handy name to mirror inform's.
		s.The("actors",
			Called("animals"),
			Exist())

		// hrmm.... implies actors are takeable.
		s.The("objects",
			Called("actors"),
			Have("clothing", "objects").
				Implying("objects", Have("wearer", "actor")),
			Have("inventory", "objects").
				Implying("objects", Have("owner", "actor")))

		// changed: used to have "equipment" for held objects
		// will either implement some sort of "relation with value"
		// or will add a "held","holdable", state to objects.

		s.The("objects",
			Called("props"),
			AreEither("portable").Or("fixed in place"),
		)

		// Usually opaque not transparent, open not closed, unopenable not openable, unlocked not locked.
		// Usually not enterable, lockable.
		// Can have carrying capacity (number)
		s.The("props",
			Called("containers"),
			Have("contents", "objects").
				Implying("objects", Have("enclosure", "container")),
			AreEither("opaque").Or("transparent"),
		//AreEither("enterable").Or("not enterable"),
		//AreEither("lockable").Or("unlockable"),
		//AreEither("locked").Or("unlocked"),
		)

		s.The("props",
			Called("supporters"),
			Have("contents", "objects").
				Implying("objects", Have("support", "supporter")))

		s.The("props",
			Called("devices"),
			AreEither("switched off").Or("switched on"))
	})
}

// will have to become more sophisticated for being inside a box.
func actorLocation(report string) G.Callback {
	return func(g G.Play) {
		actor := g.The("actor")
		target := actor.Object("whereabouts")
		target.Go(report, actor)
	}
}

func actorTarget(report string) G.Callback {
	return func(g G.Play) {
		actor := g.The("actor")
		target := g.The("action.Target")
		target.Go(report, actor)
	}
}

func init() {
	AddScript(func(s *Script) {
		s.The("actors",
			Can("jump").And("jumping").RequiresNothing(),
			To("jump", actorLocation("report jumping")),
		)

		s.The("actors",
			Can("listen").And("listening").RequiresNothing(),
			To("listen", actorLocation("report the sound")),

			Can("listen to").And("listening to").RequiresOne("kind"),
			To("listen to", actorTarget("report the sound")),
		)

		// one visible thing, and requring light
		s.The("actors",
			Can("look").And("looking").RequiresNothing(),
			To("look", //actorLocation("report the view")
				func(g G.Play) {
					actor := g.The("actor")
					target := actor.Object("whereabouts")
					target.Go("report the view")
				}),
		)

		// one visible thing and requiring light.
		s.The("actors",
			Can("look under it").And("looking under it").RequiresOne("object"),
			To("look under it", actorTarget("look under")),
		)

		// "taking inventory" in inform
		// again, as with some other actions: for players this happens in carry out, for npcs in report.
		// i'm sure that's useful... somehow....
		s.The("actors",
			Can("report inventory").And("reporting inventory").RequiresNothing(),
			To("report inventory", func(g G.Play) {
				this := g.The("actor")
				source := []string{"clothing", "inventory"}
				for _, s := range source {
					contents := this.ObjectList(s)
					if len(contents) > 0 {
						g.Say(s, ":")
						for _, obj := range contents {
							obj.Go("print name")
						}
					} else {
						g.Say(s, "(none).")
					}
				}
			}),
		)

		// searching: requiring light; FIX: what does searching a room do?
		s.The("actors",
			Can("search it").And("searching it").RequiresOne("prop"),
			To("search it", actorTarget("search")))
		s.The("props",
			Can("search").And("searching").RequiresOne("actor"))

		// smelling
		s.The("actors",
			Can("smell").And("smelling").RequiresNothing(),
			To("smell", actorLocation("report the smell")),

			Can("smell it").And("smelling it").RequiresOne("kind"),
			To("smell it", actorTarget("report the smell")),
		)

		// WARNING/FIX: multi-word statements must appear before their single word variants
		// ( or the parser will attempt to match the setcond word as a noun )
		s.Execute("jump", Matching("jump|skip|hop"))
		s.Execute("search it", Matching("search {{something}}").
			Or("look inside|in|into|through {{something}}"))
		s.Execute("smell it", Matching("smell|sniff {{something}}"))
		s.Execute("smell", Matching("smell|sniff"))
		// FIX: for some reason, the order must be biggest match to smallest, the other way doesnt work.
		s.Execute("report inventory", Matching("inventory|inv|i"))
		s.Execute("look under it", Matching("look under {{something}}"))
		s.Execute("look", Matching("look|l"))
		s.Execute("listen to", Matching("listen to {{something}}").Or("listen {{something}}"))
		s.Execute("listen", Matching("listen"))
	})
}

//
// when is the right time for functions versus callbacks?
func listContents(g G.Play, header string, this G.IObject) (printed bool) {
	// if something described which is not scenery is on the noun and something which is not the player is on the noun:
	// obviously a filterd callback, visitor, would be nice FilterList("contents", func() ... )
	contents := this.ObjectList("contents")
	if len(contents) > 0 {
		g.Say(header, this.Name(), "is:")
		for _, obj := range contents {
			obj.Go("print description")
		}
		g.Say("")
		printed = true
	}
	return printed
}

//
// System actions
func init() {
	AddScript(func(s *Script) {
		s.The("kinds",
			Can("report jumping").And("reporting jumping").RequiresOne("actor"),
			To("report jumping", func(g G.Play) {
				actor := g.The("action.Target")
				// FIX? inform often, but not always, tests for trying silently,
				// "if the action is not silent" ...
				// seems... strange. why report if if its silent?
				if g.The("player") == actor {
					g.Say("You jump on the spot.")
				} else {
					g.Say(actor.Name(), "jumps on the spot.")
				}
			}))

		// inform has two entries for some actions (looking under as an example, jumping as a counter example):
		// 1. carry out an actor looking under: if the player
		// 2. report an actor looking under: if not the player
		// its not exactly clear to me why, the docs give guidelines for this, but not rationale.
		// it might be interesting to queue says, if they need to be cancelled or held back.
		// keep in mind, most of these really want to be animations, and only sometimes voice.
		s.The("objects",
			Can("look under").And("looking under").RequiresOne("actor"),
			To("look under", func(g G.Play) {
				source, actor := g.The("action.Source"), g.The("action.Target")
				if g.The("player") == actor {
					g.Say("You find nothing of interest.")
				} else {
					g.Say(actor.Name(), "looks under the", source.Name(), ".")
				}
			}))

		// kinds, to allow rooms and objects
		s.The("kinds",
			Can("report the smell").And("reporting the smell").RequiresOne("actor"),
			To("report the smell", func(g G.Play) {
				actor := g.The("action.Target")
				if g.The("player") == actor {
					g.Say("You smell nothing unexpected.")
				} else {
					g.Say(actor.Name(), "sniffs.")
				}
			}),
			Can("report the sound").And("reporting the sound").RequiresOne("actor"),
			To("report the sound", func(g G.Play) {
				actor := g.The("action.Target")
				if g.The("player") == actor {
					g.Say("You hear nothing unexpected.")
				} else {
					g.Say(actor.Name(), "listens.")
				}
			}))

		s.The("objects",
			Can("print name").And("printing name text").RequiresNothing(),
			To("print name", func(g G.Play) {
				obj := g.The("object")
				g.Say(obj.Name())
			}))

		s.The("containers",
			When("printing name text").
				Always(func(g G.Play) {
				// FIX: conditional return instead of Always?
				// or some way ( dependency injection ) to get at the event object
				// of course, rules producing values and stacks might work too.
				this := g.The("container")
				list := this.ObjectList("contents")
				if this.Is("transparent") && len(list) == 0 {
					g.Say(this.Name(), "(empty)")
				} else {
					g.Say(this.Name())
				}
				g.StopHere()
			}))

		s.The("rooms",
			Can("report the view").And("room describing").RequiresNothing(),
			To("report the view", func(g G.Play) {
				room := g.The("room")
				g.Say(room.Name())
				//
				desc := room.Text("description")
				if desc == "" {
					// FIX: this .Name() is confusing: possibly support obj.Text("name") instead?
					desc = room.Name()
				}
				g.Say(desc)
				for _, obj := range room.ObjectList("contents") {
					obj.Go("print description")
				}
			}))
	})
}
