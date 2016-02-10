package standard

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/util/lang"
	"strings"
)

// FIX: we have the concept floating in other fixes of "function" globals
// and that might be needed for this, where we really dont want *shared* globals
// you would want this tied to session, if at all possible.
var Debugging bool

// InitStandardLibrary to ensure all standard library scripts are properly created.
func InitStandardLibrary() *Script {
	// the side effect of importing the standard library
	// adds the scripts to the list of those which get initialized.
	return InitScripts()
}

// FIX: there's no error testing here and its definitely possible to screw things up.
func AssignTo(prop G.IObject, rel string, dest G.IObject) {
	// sure hope there's no errors, would relation by value remove the need for transaction?
	if _, parentRel := prop.ParentRelation(); parentRel != "" {
		// note: an object like the fishFood isnt "in the world", and doest have an owner field to clear.
		prop.Set(parentRel, nil)
	}
	prop.Set(rel, dest)
}

//
func init() {
	AddScript(func(s *Script) {
		// FIX: hierarchy is a work in progress.
		// kinds> stories, rooms, objects > actors (> animals),  props(> openers(> doors,containers), supporters, devices)

		// FIX: IObject and Instance, ObjectList are confusing versus kind/object
		// Instance might be better...? or some better name for the class of objects
		// i tried "subjects" but hated it more.
		s.The("kinds",
			// printed plural name (text),
			Have("printed name", "text"),
			// FIX: missing a vs. the?
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

		// CAN WE DEFAULT (USUALLY(X)) DOORS TO fixed-in-place???
		s.The("doors", Before("reporting take").Always(func(g G.Play) {
			g.Say("It is fixed in place.")
			g.StopHere()
		}))

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
func ReflectToLocation(g G.Play, action string) {
	actor := g.The("actor")
	target := actor.Object("whereabouts")
	//g.Log("reflecting", action, actor, target)
	target.Go(action, actor)
}

// ReflectToTarget runs the passed action, flipping the source and target.
func ReflectToTarget(g G.Play, action string) {
	source := g.The("action.Source")
	target := g.The("action.Target")
	//g.Log("reflecting", action, source, target)
	target.Go(action, source)
}

// ReflectWithContext runs the passed action, shifting to target, context, source.
// FIX: i think it'd be better to first use ReflectToTarget, keeping the context as the third parameter
// and then ReflectToContext, possibly re-swapping source and target.
func ReflectWithContext(g G.Play, action string) {
	source := g.The("action.Source")
	target := g.The("action.Target")
	context := g.The("action.Context")
	//g.Log("reflecting", action, source, target, context)
	target.Go(action, context, source)
}

func init() {
	AddScript(func(s *Script) {

		// one visible thing, and requring light
		s.The("actors",
			Can("look").And("looking").RequiresNothing(),
			To("look",
				// func( g G.Play) { ReflectToLocation(g,"report the view") }
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
			To("look under it", func(g G.Play) { ReflectToTarget(g, "report look under") }),
		)

		// FIX: should generate a report/response instead?
		s.The("actors",
			Can("impress").And("impressing").RequiresNothing(),
			To("impress", func(g G.Play) {
				g.Say(lang.Capitalize(ArticleName(g, "actor", nil)), "is unimpressed.")
			}))

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
func ListContents(g G.Play, header string, obj G.IObject) (printed bool) {
	// if something described which is not scenery is on the noun and something which is not the player is on the noun:
	// obviously a filterd callback, visitor, would be nice FilterList("contents", func() ... )
	// FIX: if something has scenery objets, they appear as contents,
	// but then the list is empty. ( ex. lab coat, but it might happen elsewhere )
	// we'd maybe need to know if something printed?
	if contents := obj.ObjectList("contents"); len(contents) > 0 {
		g.Say(header, obj.Text("Name"), "is:")
		for _, content := range contents {
			content.Go("print description")
		}
		printed = true
	}
	return printed
}

func NameFullStop(G.IObject) string {
	return ""
}

type NameStatus func(obj G.IObject) string

func ArticleName(g G.Play, which string, status NameStatus) string {
	return articleName(g, which, false, status)
}
func DefiniteName(g G.Play, which string, status NameStatus) string {
	return articleName(g, which, true, status)
}
func articleName(g G.Play, which string, definite bool, status NameStatus) string {
	obj := g.The(which)
	text := obj.Text("Name")
	if obj.Is("proper-named") {
		text = lang.Titleize(text)
	} else {
		article := ""
		if definite {
			article = "the"
		} else {
			article = obj.Text("indefinite article")
			if article == "" {
				if obj.Is("plural-named") {
					article = "some"
				} else if lang.StartsWithVowel(text) {
					article = "an"
				} else {
					article = "a"
				}
			}
		}
		text = strings.Join([]string{article, strings.ToLower(text)}, " ")
	}
	if status != nil {
		if s := status(obj); s != "" {
			text = fmt.Sprintf("%s (%s)", text, s)
		} else {
			text = text + "."
		}
	}
	return text
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
				if text := ArticleName(g, "object", NameFullStop); len(text) > 0 {
					text = lang.Capitalize(text)
					g.Say(text)
				}
			}))

		s.The("containers",
			When("printing name text").
				Always(func(g G.Play) {
				// FIX: conditional return instead of Always?
				// or some way ( dependency injection ) to get at the event object
				// of course, rules producing values and stacks might work too.
				// FIX: a container is an opener... where do we print the opener status name
				// put this on doors for now.
				text := ArticleName(g, "container", func(obj G.IObject) (status string) {
					if obj.Is("closed") {
						if obj.Is("hinged") {
							status = "closed"
						}
					} else if list := obj.ObjectList("contents"); len(list) == 0 {
						if obj.Is("transparent") || obj.Is("open") {
							status = "empty"
						}
					}
					return status
				})
				text = lang.Capitalize(text)
				g.Say(text)
				g.StopHere()
			}))

		s.The("doors",
			When("printing name text").
				Always(func(g G.Play) {
				text := DefiniteName(g, "door", func(obj G.IObject) (status string) {
					if obj.Is("hinged") {
						if obj.Is("open") {
							status = "open"
						} else {
							status = "closed"
						}
					}
					return status
				})
				text = lang.Capitalize(text)
				g.Say(text)
				g.StopHere()
			}))

		s.The("rooms",
			Can("report the view").And("reporting the view").RequiresNothing(),
			When("reporting the view").Always(func(g G.Play) {
				room := g.The("room")
				g.The("status bar").SetText("left", lang.Titleize(room.Text("Name")))
			}),
			After("reporting the view").Always(func(g G.Play) {
				g.Go(Change("room").To("visited"))
			}),
			To("report the view", func(g G.Play) {
				g.Go(View("room"))
			}))
	})
}

// FUTURE? interestingly, while we wouldnt be able to encode them without special work, the contents of the phrases are fixed: we could have After("reporting").Execute(Phrase). maybe "standard" phrases could put themselves in some sort of wrapper? around the model? tho not quite sure howd that work.
type ViewRoomPhrase struct {
	object string
}

func View(object string) ViewRoomPhrase {
	return ViewRoomPhrase{object}
}

func (p ViewRoomPhrase) Execute(g G.Play) {
	room := g.The("room")
	// sometines a blank like is printed without this
	// (maybe certain directions? or going through doors?)
	// not sure why, so leaving this for consistency
	g.Say(Lines("", room.Get("Name").Text()))
	g.Say(Lines(room.Get("description").Text(), ""))
	// FIX? uses 1 to exclude the player....
	// again, this happens because we dont know if print description actually did anything (re:scenery, etc.)
	if contents := room.ObjectList("contents"); len(contents) > 1 {
		for _, obj := range contents {
			obj.Go("print description")
			g.Say("")
		}
	}
}
