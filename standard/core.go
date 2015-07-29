package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	"strings"
)

// InitStandardLibrary to ensure all standard library scripts are properly created.
func InitStandardLibrary() *Script {
	// the side effect of importing the standard library
	// adds the scripts to the list of those which get initialized.
	return InitScripts()
}

// FIX: there's no error testing here and its definitely possible to screw things up.
func assignTo(prop G.IObject, rel string, dest G.IObject) {
	// sure hope there's no errors, would relation by value remove the need for transaction?
	if _, parentRel := DirectParent(prop); parentRel != "" {
		// note: an object like the fishFood isnt "in the world", and doest have an owner field to clear.
		prop.Set(parentRel, nil)
	}
	prop.Set(rel, dest)
}

//var articles = regexp.MustCompile(`^(The|A|An|Our|Some)[[:upper:]]`)
//http://www.mudconnect.com/SMF/index.php?topic=74725.0
func sin(s string, set ...string) (ok bool) {
	for _, x := range set {
		if x == s {
			ok = true
			break
		}
	}
	return ok
}

func startsVowel(str string) (vowelSound bool) {
	s := strings.ToUpper(str)
	if sin(s, "A", "E", "I", "O", "U") {
		if !sin(s, "EU", "EW", "ONCE", "ONE", "OUI", "UBI", "UGAND", "UKRAIN", "UKULELE", "ULYSS", "UNA", "UNESCO", "UNI", "UNUM", "URA", "URE", "URI", "URO", "URU", "USA", "USE", "USI", "USU", "UTA", "UTE", "UTI", "UTO") {
			vowelSound = true
		}
	} else if sin(s, "HEIR", "HERB", "HOMAGE", "HONEST", "HONOR", "HONOUR", "HORS", "HOUR") {
		vowelSound = true
	}
	return vowelSound
}

//
func init() {
	AddScript(func(s *Script) {
		// FIX: hierarchy is a work in progress.
		// kinds> stories, rooms, objects > actors (> animals),  props(> openers(> doors,containers), supporters, devices)

		// FIX: IObject and GameObject, ObjectList are confusing versus kind/object
		// Instance might be better...? or some better name for the class of objects
		// i tried "subjects" but hated it more.
		s.The("kinds",
			// printed plural name (text),
			Have("printed name", "text"),
			Have("indefinite article", "text"), //
			AreEither("singular-named").Or("plural-named"),
			AreEither("common-named").Or("proper-named"), //a name used for an individual person, place, or organization, spelled with initial capital letters, e.g., Larry, Mexico, and Boston Red Sox.
			//not normally preceded by an article or other limiting modifier, as any or some. Nor are they usually pluralized
		)

		// vs. descriptions in "kind"
		// it seems to make sense for now to have two separate description fields.
		// rooms like to say their description, while objects like to say their brief initial appearance ( or name, if there's none. )
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
			HaveMany("contents", "objects").
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

		s.The("openers",
			Called("doors"),
			Exist())

		// nothing special: just a handy name to mirror inform's.
		s.The("actors",
			Called("animals"),
			Exist())

		// hrmm.... implies actors are takeable.
		s.The("objects",
			Called("actors"),
			HaveMany("clothing", "objects").
				Implying("objects", HaveOne("wearer", "actor")),
			HaveMany("inventory", "objects").
				Implying("objects", HaveOne("owner", "actor")))

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
		s.The("openers",
			Called("containers"),
			HaveMany("contents", "objects").
				Implying("objects", HaveOne("enclosure", "container")),
			AreEither("opaque").Or("transparent"),
			//AreEither("enterable").Or("not enterable"),
			AreEither("lockable").Or("not lockable").Usually("not lockable"),
			AreEither("locked").Or("unlocked").Usually("unlocked"),
		)

		s.The("props",
			Called("supporters"),
			HaveMany("contents", "objects").
				Implying("objects", HaveOne("support", "supporter")))
	})
}

// reflect to the passed action passing the actors's current whereabouts.
// will have to become more sophisticated for being inside a box.
func ReflectToLocation(action string) G.Callback {
	return func(g G.Play) {
		actor := g.The("actor")
		target := actor.Object("whereabouts")
		target.Go(action, actor)
	}
}

// ReflectToTarget runs the passed action, flipping the source and target.
func ReflectToTarget(action string) G.Callback {
	return func(g G.Play) {
		source := g.The("action.Source")
		target := g.The("action.Target")
		target.Go(action, source)
	}
}

// ReflectWithContext runs the passed action, shifting to target, context, source.
// FIX: i think it'd be better to first use ReflectToTarget, keeping the context as the third parameter
// and then ReflectToContext, possibly re-swapping source and target.
func ReflectWithContext(action string) G.Callback {
	return func(g G.Play) {
		source := g.The("action.Source")
		target := g.The("action.Target")
		context := g.The("action.Context")
		target.Go(action, context, source)
	}
}

func init() {
	AddScript(func(s *Script) {

		// one visible thing, and requring light
		s.The("actors",
			Can("look").And("looking").RequiresNothing(),
			To("look",
				// ReflectToLocation("report the view")
				// reflect to location will send the actor as a parameter,
				// but report the view doesnt expect parameters.
				func(g G.Play) {
					actor := g.The("actor")
					target := actor.Object("whereabouts")
					target.Go("report the view")
				}),
		)

		// one visible thing and requiring light.
		s.The("actors",
			Can("look under it").And("looking under it").RequiresOne("object"),
			To("look under it", ReflectToTarget("report look under")),
		)

		// "taking inventory" in inform
		// again, as with some other actions: for players this happens in carry out, for npcs in report.
		// i'm sure that's useful... somehow....
		s.The("actors",
			Can("report inventory").And("reporting inventory").RequiresNothing(),
			To("report inventory", func(g G.Play) {
				this := g.The("actor")
				source := []string{"Clothing", "Inventory"}
				for _, s := range source {
					contents := this.ObjectList(s)
					if len(contents) > 0 {
						g.Say(s + ":")
						for _, obj := range contents {
							obj.Go("print name")
						}
					} else {
						g.Say(s, "(none).")
					}
				}
			}),
		)

		// FIX: for some reason, the order must be biggest match to smallest, the other way doesnt work.
		s.Execute("report inventory", Matching("inventory|inv|i"))
		s.Execute("look", Matching("look|l"))
		s.Execute("look under it", Matching("look under {{something}}"))
	})
}

//
// when is the right time for functions versus callbacks?
func listContents(g G.Play, header string, obj G.IObject) (printed bool) {
	// if something described which is not scenery is on the noun and something which is not the player is on the noun:
	// obviously a filterd callback, visitor, would be nice FilterList("contents", func() ... )
	contents := obj.ObjectList("contents")
	// FIX: if something has scenery objets, they appear as contents,
	// but then the list is empty. ( ex. lab coat, but it might happen elsewhere )
	// we'd maybe need to know if something printed?
	if len(contents) > 0 {
		g.Say(header, obj.Text("Name"), "is:")
		for _, content := range contents {
			content.Go("print description")
		}
		printed = true
	}
	return printed
}

//
// System actions
func init() {
	AddScript(func(s *Script) {

		// inform has two entries for some actions (looking under as an example, jumping as a counter example):
		// 1. carry out an actor looking under: if the player
		// 2. report an actor looking under: if not the player
		// its not exactly clear to me why, the docs give guidelines for this, but not rationale.
		// it might be interesting to queue says, if they need to be cancelled or held back.
		// keep in mind, most of these really want to be animations, and only sometimes voice.
		s.The("objects",
			Can("report look under").And("reporting look under").RequiresOne("actor"),
			To("report look under", func(g G.Play) {
				source, actor := g.The("action.Source"), g.The("action.Target")
				if g.The("player") == actor {
					g.Say("You find nothing of interest.")
				} else {
					g.Say(actor.Text("Name"), "looks under the", source.Text("Name"), ".")
				}
			}))

		s.The("objects",
			Can("print name").And("printing name text").RequiresNothing(),
			To("print name", func(g G.Play) {
				obj := g.The("object")
				text := obj.Text("Name")
				if obj.Is("proper-named") {
					text = strings.Title(text)
				} else {
					article := obj.Text("indefinite article")
					if article == "" {
						if startsVowel(text) {
							article = "An"
						} else {
							article = "A"
						}
					}
					text = strings.Join([]string{article, strings.ToLower(text)}, " ")
				}
				g.Say(text + ".")
			}))

		s.The("containers",
			When("printing name text").
				Always(func(g G.Play) {
				// FIX: conditional return instead of Always?
				// or some way ( dependency injection ) to get at the event object
				// of course, rules producing values and stacks might work too.
				this := g.The("container")
				list := this.ObjectList("contents")
				// FIX: a container is an opener... where do we print the opener status name
				// put this on doors for now.
				if this.Is("transparent") && len(list) == 0 {
					g.Say(this.Text("Name"), "(empty)")
				} else {
					g.Say(this.Text("Name"))
				}
				g.StopHere()
			}))

		s.The("doors",
			When("printing name text").
				Always(func(g G.Play) {
				this := g.The("door")
				if this.Is("openable") {
					x := map[bool]string{true: "(open)", false: "(closed)"}
					status := x[this.Is("open")]
					g.Say(this.Text("Name"), status)
					g.StopHere()
				}
			}))

		s.The("rooms",
			Can("report the view").And("reporting the view").RequiresNothing(),
			To("report the view", func(g G.Play) {
				room := g.The("room")
				g.Say(Lines("", room.Text("Name")))
				g.Say(Lines(room.Text("description"), ""))
				if contents := room.ObjectList("contents"); len(contents) > 0 {
					g.Say("You can see:")
					for _, obj := range contents {
						obj.Go("print description")
					}
				}
				// FIX: duplicated in stories describe the first room
				room.IsNow("visited")
				g.The("status bar").SetText("left", strings.Title(room.Text("Name")))
			}))
	})
}
