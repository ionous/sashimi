package fishgen

import (
	"fmt"
	. "github.com/ionous/sashimi/extensions"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/lang"
	"strings"
)

// from script...
func Lines(a ...string) string {
	return strings.Join(a, "\n")
}

var Callbacks = map[ident.Id]G.Callback{

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/openingClosing.go:25
	"a028e313823833dd7dad493203bd86fc": func(g G.Play) { ReflectToTarget(g, "be closed by") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/kissing.go:13
	"a03c941d25050d3a9974d9165bb8f1": func(g G.Play) {
		source := g.The("action.Source")
		g.Say(source.Text("Name"), "might not like that.")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/puttingItOn.go:35
	"a0ca28fd1e1f6003219b799c59e5925": func(g G.Play) {
		g.Go(Put("action.Source").Onto("action.Context"))
		g.Say("You put", ArticleName(g, "action.Source", nil), "onto", ArticleName(g, "action.Context", NameFullStop))
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/openingClosing.go:25
	"a0cf32d8daa137026c62f8494129": func(g G.Play) {
		g.Say("It's already closed.") //[regarding the noun]?
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/stories.go:37
	"a1157c67ba6d0bacc2a189d6ac0f88": func(g G.Play) {
		story := g.The("story")
		g.Say("*** The End ***")
		story.IsNow("completed")

		if story.Is("scored") {
			score, maxScore, turnCount := story.Num("score"), story.Num("maximum score"), story.Num("turn count")
			g.Say(fmt.Sprintf("In that game you scored %d out of a possible %d, in %d turns.",
				int(score), int(maxScore), int(turnCount)+1,
			))
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/openingClosing.go:25
	"a14a3d5d1cca145f13f50ca1309af311": func(g G.Play) {
		prop := g.The("prop")
		g.Say("Now the", prop.Text("Name"), "is closed.")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/insertingInto.go:49
	"a174d13dbae209ddf968f9f16e9": func(g G.Play) {
		prop := g.The("action.Context")
		if worn := prop.Object("wearer"); worn.Exists() {
			g.Say("You can't insert worn clothing.")
			// FIX: try taking off the noun
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/movement.go:96
	"a279a3ec2a6f0098be46df9a674f96f": func(g G.Play) {
		actor, dir := g.The("actor"), g.The("action.Target")
		from := actor.Object("whereabouts")
		// try the forward direction:
		departingDoor := from.Object(dir.Text("Name") + "-via")
		if departingDoor.Exists() {
			TryMove(actor, dir, departingDoor)
		} else {
			// // try the opposite direction link:
			// rev := dir.Object("opposite")
			// exit := from.Object(rev.Text("Name") + "-rev-via")
			// if exit.Exists() {
			// 	if sources := exit.ObjectList("sources"); len(sources) == 1 {
			// 		TryMove(actor, dir, sources[0])
			// 	}
			//} else {
			if Debugging {
				g.Log("couldnt find %s exit", dir)
			}
			g.Say("You can't move that direction.")
			//}
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"a27a04f5b30dcdb2fb3bca460ffa83a3": func(g G.Play) {
		if PlayerLearns(g, "examinedTelegraph") {
			g.The("evil fish").Says(`"So," blubs the evil fish. "How about it? Little food over here?"`)
		} else {
			fishComments := []string{
				//Table of Insulting Fish Comments
				`"Yeah, yeah," says the fish. "You having some trouble with the message, there? Confused? Something I could clear up for you?"`,
				`"Oookay, genius kid has some troubles in the reading comprehension department." The fish taps his head meaningfully against the side of the tank. "I'm so hungry I could eat my way out, you get my meaning?"`,
				`"I'll translate for you," screams the fish in toothy fury. "It says GIVE FOOD TO FISH!! How much more HELP do you NEED???"`,
			}
			i := g.Random(len(fishComments))
			comment := fishComments[i]
			g.The("evil fish").Says(comment)
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/wearing.go:37
	"a2afcc3585516a895061346c3d5fa3f": func(g G.Play) { ReflectToTarget(g, "report wear") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/giving.go:48
	"a38fd7888ecbfb8e121d362dbfed": func(g G.Play) { ReflectToTarget(g, "be acquired") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/core.go:180
	"a4e06ed8921a1c0e863719f30eed2": func(g G.Play) {
		g.Say(lang.Capitalize(ArticleName(g, "actor", nil)), "is unimpressed.")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"a5003a18e42eccfb92e40558889d715": func(g G.Play) {
		g.Say("You're saving all your lovin for someone a lot cuddlier.")
		g.StopHere()
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"a5c72be52c4f2336436967baf5754621": func(g G.Play) {
		if g.The("window").Is("open") {
			g.Say(`Through the windows you get a lovely view of the street outside. At the moment, the glass is thrown open, and a light breeze is blowing through.`)
		} else {
			g.Say(`Through the windows, you get a lovely view of the street outside -- the little fountain on the corner, the slightly dilapidated but nonetheless magnificent Jugendstil architecture of the facing building. The glass itself is shut, however.`)
		}
		g.StopHere()
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"a60e331213831e69ad2e4bc2294f1": func(g G.Play) {
		if PlayerLearns(g, "openedWindow") {
			g.The("evil fish").Says(`"Thank god some air," says the fish. "Man, it was getting hard to breathe in here." Two beats pass. "Oh wait."`)
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"a795e203a59b3f39adb4b28ca8eef4c3": func(g G.Play) {
		fish := g.The("evil fish")
		receiver := g.The("action.Context")
		if fish == receiver {
			fish.Says("What are you, some kind of sadist? I don't want to see a bunch of cloths! What kind of f'ing good, 'scuse my French, is that supposed to do me? I don't even wear pants for God's sake!")
			g.Say("He really looks upset. You start wondering whether apoplexy is an ailment common to fish.")
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/puttingItOn.go:35
	"a8e0802cdf9350defff5d67b08b6": func(g G.Play) {
		supporter, prop := g.The("action.Target"), g.The("action.Context")
		if supporter == prop {
			g.Say("You can't put something onto itself.")
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"a98fd9211e9759d489fb4522c66d852f": func(g G.Play) {
		g.The("evil fish").Says(`The evil fish notices you sniffing the air. "Vanilla Raspberry Roast," it remarks. "You really miss her, don't you."`)
		g.Say("You glance over, startled, but the fish's mouth is open in a piscine equivalent of a laugh. You stifle the urge to skewer the thing...")
		g.StopHere()
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/giving.go:48
	"aa0f3ad3144d3b009591e22f4f7f9aa2": func(g G.Play) {
		actor, prop, rel := g.The("actor"), g.The("prop"), "owner"
		if Debugging {
			par, prev := prop.ParentRelation()
			g.Log(prop, "AssignTo", actor, rel, "from", par, prev)
		}
		AssignTo(prop, rel, actor)
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"aa170312949d078407f725833b2f5060": func(g G.Play) {
		if PlayerLearns(g, "examinedCloths") {
			g.The("evil fish").Says("Whatcha looking at? I can't see through the doors, you know.")
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/extensions/conversation.go:107
	"ab0570189d883833d088ab9206526dc3": func(g G.Play) {
		con := TheConversation(g)
		if npc := con.Depart(); npc.Exists() {
			if Debugging {
				g.Log("!", g.The("actor"), "departing", npc)
			}
			g.Say("(", lang.Capitalize(DefiniteName(g, "actor", nil)), "says goodbye.", ")")
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/taking.go:27
	"ad645ba09d5207de6d3ab358466b7355": func(g G.Play) { ReflectToTarget(g, "report take") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/puttingItOn.go:35
	"adc65b6c27ecd7b4aabf6bdd973b": func(g G.Play) {
		supporter := g.The("action.Target")
		if supporter.Is("closed") {
			g.Say(ArticleName(g, "action.Target", nil), "is closed.")
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/showing.go:21
	"aea3c76840abada678c70a2122fe83ff": func(g G.Play) {
		presenter, _, prop := g.The("action.Source"), g.The("action.Target"), g.The("action.Context")
		if carrier := Carrier(prop); carrier != presenter {
			g.Say("You aren't holding", ArticleName(g, "action.Context", NameFullStop))
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"afe1549f4f3744702939c3f69275690": func(g G.Play) {
		if PlayerLearns(g, "examinedFishOnce") {
			g.Say("The fish glares at you, as though to underline this point.")
		} else if PlayerLearns(g, "examinedFishTwice") {
			g.The("evil fish").Says(`"If you're looking for signs of malnutrition," says the fish, "LOOK NO FURTHER!!" And it sucks in its gills until you can see its ribcage.`)
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"b0b3278650e296e6c97cd7ab4ad833a": func(g G.Play) {
		if PlayerLearns(g, "examinedPaints") {
			g.The("evil fish").Says(`"Tons of useful stuff in there," hollers in the fish, in a syncopated burble.`)
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/listening.go:12
	"b0b8e208c4cc475e100e14fbcf2f4ae": func(g G.Play) {
		actor := g.The("action.Target")
		if g.The("player") == actor {
			g.Say("You hear nothing unexpected.")
		} else {
			g.Say(actor.Text("Name"), "listens.")
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/core.go:180
	"b0d267fe6bec1f08710b46623fe815b": func(g G.Play) { ReflectToTarget(g, "report look under") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"b35f0712de0ecce8238feaa309229a9": func(g G.Play) {
		if g.The("object") == g.The("fish food") {
			g.The("actor").Go("feed it", g.The("evil fish"))
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/devices.go:36
	"b44ad3f16c77942ddebfbcbbdb1381c": func(g G.Play) {
		g.Say("It's already switched on.")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/openingClosing.go:25
	"b53ef45a3e1238a3aa47a2bdc8d9da4": func(g G.Play) {
		g.The("actor").Says("It's locked!")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/showing.go:21
	"b5bd6ab4d4df9616f49967630082b1e": func(g G.Play) { ReflectWithContext(g, "report show") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/openingClosing.go:25
	"b6e604472836e917e8ad1dd6ba2b93": func(g G.Play) { ReflectToTarget(g, "be opened by") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/extensions/conversation.go:107
	"b703927b07424ab1a58cc61fc9287c": func(g G.Play) { ReflectToTarget(g, "be discussed") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/devices.go:36
	"b703969791fe6a4b4b43764fe9215df0": func(g G.Play) {
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
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/examining.go:17
	"b8fe8f66db86e05e87b099ee578f0c": func(g G.Play) {
		if c := g.The("container"); c.Is("open") || c.Is("transparent") {
			ListContents(g, "In the", c)
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"b9f230f5cf3e2fed1850f455e7": func(g G.Play) {
		receiver := g.The("action.Context")
		fish := g.The("evil fish")
		if fish == receiver {
			fish.Says("That's the ticket, sweetie! Bring it on.")
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"bb18d25b98660857d88e5f8d3747cdd": func(g G.Play) {
		if PlayerLearns(g, "examinedPainting") {
			g.The("evil fish").Says(`A ferocious banging from the aquarium attracts your attention as you go to look at the painting. "Hey!" screams the fish. "She doesn't like strangers looking at her paintings before they're DOONNNE!"`)
			g.The("player").Says(`"Shut up, you," you reply casually. "I'm not a stranger." But the fish puts you off a little bit, and your heart is already in your mouth before you see the painting itself...`)
		} else {
			g.Say("Once is really enough. It's pretty much embedded in your consciousness now.")
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/debugging.go:13
	"bb4ae3b5d6669a4eb5d919a52800816": func(g G.Play) {
		target := g.The("action.Target")
		parent, relation := target.ParentRelation()
		if relation == "" {
			g.Say(target.Text("Name"), "=>", "out of world")
		} else {
			g.Say(target.Text("Name"), "=>", relation, parent.Text("Name"))
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/giving.go:48
	"bb5aa921e7bf8a1a2a243bc496da3d37": func(g G.Play) { ReflectWithContext(g, "report gave") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/giving.go:48
	"bbbc9f864e3bc39e4f95286879cfb14": func(g G.Play) {
		prop := g.The("action.Context")
		if worn := prop.Object("wearer"); worn.Exists() {
			g.Say("You can't give worn clothing.")
			// FIX: try taking off the noun
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/devices.go:36
	"bc33e94a52fffabe6b20bcc54955ce3": func(g G.Play) {
		device, _ := g.The("device"), g.The("actor")
		g.Say("Now the", device.Text("Name"), "is off.")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/smelling.go:12
	"bd1fabe270be6ab6a4724cc8f142da1": func(g G.Play) { ReflectToLocation(g, "report smell") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/insertingInto.go:49
	"bd9ceebcdcc7ebb57be7187da0b3916": func(g G.Play) { ReflectWithContext(g, "be inserted") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/stories.go:37
	"be2fb26f4dfcdf76f7d89ab1a25c254": func(g G.Play) {
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
		story.Go("set initial position", g.The("player"), room)
		story.Go("print the banner") // see: banner.go
		room = g.The("player").Object("whereabouts")
		// FIX: Go() should handle both Name() and ref
		story.Go("describe the first room", room)
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/core.go:180
	"be3c2951941a98acc45dfd2d4705308": func(g G.Play) {
		room := g.The("room")
		// sometines a blank like is printed without this
		// (maybe certain directions? or going through doors?)
		// not sure why, so leaving this for consistency
		g.Say(Lines("", room.Text("Name")))
		g.Say(Lines(room.Text("description"), ""))
		// FIX? uses 1 to exclude the player....
		// again, this happens because we dont know if print description actually did anything (re:scenery, etc.)
		if contents := room.ObjectList("contents"); len(contents) > 1 {
			for _, obj := range contents {
				obj.Go("print description")
				g.Say("")
			}
		}
		// FIX: duplicated in stories describe the first room
		room.IsNow("visited")
		g.The("status bar").SetText("left", lang.Titleize(room.Text("Name")))
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/giving.go:48
	"be81d6d6fb1282452bb02cb39edd0ab7": func(g G.Play) {
		presenter, receiver := g.The("action.Source"), g.The("action.Target")
		if presenter == receiver {
			g.Say("You can't give to yourself")
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/insertingInto.go:49
	"c092ea1a525fdfe12f7b255a354": func(g G.Play) {
		container, prop := g.The("action.Target"), g.The("action.Context")
		if container == prop {
			g.Say("You can't insert something into itself.")
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/giving.go:48
	"c1cd32951100ecf601c99e78a416faa": func(g G.Play) { ReflectWithContext(g, "report give") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"c239f2d147803170152e8566a7": func(g G.Play) {
		if PlayerLearns(g, "openedCabinet") {
			g.The("evil fish").Says(`"There ya go," says the fish. "The girl is getting WARMER."`)
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/showing.go:21
	"c2f1521949bf2b49274bdf1aacd83349": func(g G.Play) {
		presenter, receiver, prop := g.The("action.Source"), g.The("action.Target"), g.The("action.Context")
		if presenter == receiver {
			presenter.Go("examine it", prop)
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/searching.go:13
	"c346a8302cd8e56c361e834a576fde7": func(g G.Play) {
		g.Say("You find nothing unexpected.")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"c3955aa0d0fc455e328e1faa99333f88": func(g G.Play) {
		fish, player := g.The("evil fish"), g.The("player")
		fish.Says(`"Oh, for--!" The evil fish breaks out in exasperation and hives. "Screw the screwing around with the screwtop. SHE never has to do that."`)
		player.Says(`"Well, SHE is not here," you reply. "What do you suggest?"`)
		fish.Says(`">FEED FISH<" says the fish promptly, making fishy faces and pointing at you with his fin. "Simplicity. Try it."`)
		g.StopHere()
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"c5351efbb7da045f6635e48a3e1958d": func(g G.Play) {
		fishFood := g.The("fish food")
		if !fishFood.Is("found") {
			fishFood.IsNow("found")
			g.Go(GiveThe(fishFood).To("the player"))
			g.Say("Poking around the cloths reveals -- ha HA! -- a vehemently orange can of fish food.")
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/extensions/conversation.go:107
	"c5f3dd965e3297200df816e37587df": func(g G.Play) {
		player := g.The("player")
		if player.Is("inputing dialog") {
			if quips := GetPlayerQuips(g); len(quips) == 0 {
				if Debugging {
					g.Log("! no conversation choices !")
				}
				player.IsNow("not inputing dialog")
				player.Go("depart") // safety first
			} else {
				input := g.The("story").Get("player input").Text()
				var choice int
				if _, e := fmt.Sscan(input, &choice); e != nil {
					g.Say("Please choose a number from the menu.", input)
				} else if choice < 1 || choice > len(quips) {
					g.Log(fmt.Sprintf("Please choose a number from the menu; number: %d of %d", choice, len(quips)))
				} else {
					quip := quips[choice-1]
					if Debugging {
						g.Log("!", player, "chose", quip)
					}
					player.IsNow("not inputing dialog")
					player.Go("comment", quip)
				}
				g.StopHere()
			}
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/openingClosing.go:25
	"c76bb26a42aa9719a9a0313b2ec1193d": func(g G.Play) {
		g.Say("That's not something you can open.")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/jumping.go:12
	"c7e571f4cb23c1766247c3657077f3b9": func(g G.Play) { ReflectToLocation(g, "report jump") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/puttingItOn.go:35
	"c7ebd1166c4885183f433a6de51411e1": func(g G.Play) { ReflectWithContext(g, "report put") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/devices.go:36
	"c8f1338cac6671fc5d6920c2c0f7b893": func(g G.Play) { ReflectToTarget(g, "report switched on") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/extensions/conversation.go:107
	"c937462c89e5bdcf0dd8d0922d98": func(g G.Play) {
		Converse(g)
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/taking.go:27
	"c9d7e290147df7748f2222edd327fcb2": func(g G.Play) {
		prop, actor := g.The("prop"), g.The("actor")
		// first, only same room:
		actorCeiling, targetCeiling := Enclosure(actor), Enclosure(prop)
		//
		if actorCeiling != targetCeiling {
			g.Say("That isn't available.")
			g.Log(fmt.Sprintf("take ceiling mismatch (%v!=%v)", actorCeiling, targetCeiling))
		} else {
			if prop.Is("scenery") {
				g.Say("That isn't available.")
				g.Log("(You can't take scenery.)")
				//g.StopHere() // FIX: should be cancel
				return
			}
			if prop.Is("fixed in place") {
				g.Say("It is fixed in place.")
				//g.StopHere() // FIX: should be cancel
				return
			}
			if parent, _ := prop.ParentRelation(); parent.Exists() {
				if parent.FromClass("actors") {
					if parent != actor {
						g.Say("That'd be stealing!")
					} else {
						g.Say(ArticleName(g, "action.Target", nil), "already has that!")
					}
					return
				}
			}
			g.Go(Give("prop").To("actor"))
			// separate report action?
			if actor == g.The("player") {
				// the kat food.
				g.Say("You take", DefiniteName(g, "action.Source", NameFullStop))
			} else {
				g.Say(ArticleName(g, "action.Target", nil), "takes", DefiniteName(g, "action.Source", NameFullStop))
			}
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"ca750763ef503a6b3d7b11b970801455": func(g G.Play) {
		if PlayerLearns(g, "lookedUnderCabinet") {
			g.The("evil fish").Says(`"Dustbunnies," predicts the fish, with telling accuracy. It executes what for all the world looks like a fishy shudder. "Lemme tell you, one time I accidentally flopped outta the tank, and I was TWO HOURS on the floor with those things STARING ME IN THE NOSE. It was frightening."`)
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/describing.go:39
	"cabe59e7088cb62f7dceba9b137394": func(g G.Play) {
		g.Go(Describe("object"))
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/openingClosing.go:25
	"cce18d15e14eddefbd4c7741dce6f5ce": func(g G.Play) {
		prop, actor := g.The("prop"), g.The("actor")
		if !prop.Is("hinged") {
			prop.Go("report unopenable", actor)
		} else {
			if prop.Is("locked") {
				prop.Go("report locked", actor)
			} else {
				if prop.Is("open") {
					prop.Go("report already open", actor)
				} else {
					prop.IsNow("open")
					prop.Go("report now open", actor)
				}
			}
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/openingClosing.go:25
	"cd40d14902a558771bcc65444bc871f": func(g G.Play) {
		opener := g.The("opener")
		g.Say("The", opener.Text("Name"), "is now open.")
		// if the noun doesnt not enclose the actor
		// list the contents of the noun, as a sentence, tersely, not listing concealed items;
		// FIX? not all openers are opaque/transparent, and not all openers have contents.
		if opener.Is("opaque") {
			ListContents(g, "In the", opener)
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/openingClosing.go:25
	"cddd2b508635fcd0fcc15232ddfee0d9": func(g G.Play) {
		prop, actor := g.The("prop"), g.The("actor")
		if !prop.Is("hinged") {
			prop.Go("report not closeable", actor)
		} else {
			// FIX: locked?
			if prop.Is("closed") {
				prop.Go("report already closed", actor)
			} else {
				prop.IsNow("closed")
				prop.Go("report now closed", actor)
			}
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/puttingItOn.go:35
	"ce6c9fbff93abf1dbaa2997c357e86dc": func(g G.Play) {
		actor, prop := g.The("action.Source"), g.The("action.Context")
		if carrier := Carrier(prop); carrier != actor {
			g.Say("You aren't holding", ArticleName(g, "action.Context", NameFullStop))
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"cf141e4b8a35226dfdb1ea917e278a": func(g G.Play) {
		if PlayerLearns(g, "lookedUnderTable") {
			g.The("evil fish").Says(`"You're not going to find anything down there," whines the fish. "I mean, c'mon. It's the fricking floor. Please tell me you can see that. I can see that. I'm a myopic fish in a tank ten feet away and I can tell you there is nothing there but floor."`)
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"cfa04e62556a30555426ce11": func(g G.Play) {
		player := g.The("player")
		if c := Carrier(g.The("fish food")); c != player {
			g.Say("You need the fish food first!")
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"d01cb53d42b7e1ff4ae8bf8f8c468188": func(g G.Play) {
		if PlayerLearns(g, "examinedSeaweed") {
			g.The("evil fish").Says(`"Nice, hunh?" blubs the fish, taking a stabbing bite out of one just by way of demonstration. "Look so good I could eat it."`)
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"d3aad960bedd830691389807b6897": func(g G.Play) {
		g.The("evil fish").Says(`"Er," says the fish. "Does that, like, EVER help??"`)
		g.StopHere()
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"d3d7b2c9b0c8ddbc95756b981d73e5f": func(g G.Play) {
		// [185-192) Table of Fish Banter
		banter := []string{
			`"Hey, nice SKIN TONE," shouts the evil fish. His words reach you in a spitting gurgle of aquarium water. "You gone over to a pure eggplant diet these days?"`,
			"The evil fish is floating belly up! ...oh, curse. He was toying with you. As soon as he sees you looking, he goes back to swimming around.",
			"The evil fish darts to the bottom of the tank and moves the gravel around with his nose.",
			"The evil fish is swimming around the tank in lazy circles.",
			"The evil fish begins to butt his pointy nose against the glass walls of the tank.",
		}
		i := g.Random(len(banter))
		comment := banter[i]
		if i == 0 {
			g.The("evil fish").Says(comment)
		} else {
			g.Say(comment)
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/stories.go:37
	"d4983e203d6eeadd7b989e46d6": func(g G.Play) {
		room := g.The("action.Target")
		/// FIX: visited should happen elsewhere
		room.Go("report the view")
		room.IsNow("visited")
		g.The("status bar").SetText("left", lang.Titleize(room.Text("Name")))
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/devices.go:36
	"d5a222bae6a821c254abbe62ff2d879": func(g G.Play) {
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
		text = lang.Capitalize(text)
		g.Say(text)
		g.StopHere()
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/devices.go:36
	"d61637c94b5914263b9b2e888c6ee3f0": func(g G.Play) { ReflectToTarget(g, "report switch off") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/openingClosing.go:25
	"d621a08c3de87b0ee4cfe7971a9": func(g G.Play) {
		g.Say("That's not something you can close.")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/core.go:180
	"d6fb6189d5ad739d1f229b238d88c": func(g G.Play) {
		g.Say("It is fixed in place.")
		g.StopHere()
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/extensions/conversation.go:107
	"d73f8083c58056e0b106f0330cabd14e": func(g G.Play) {
		g.Go(Introduce("action.Source").To("action.Target").WithDefault())
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/giving.go:48
	"d81e1bcae38ca5faa6da34c0abc814e": func(g G.Play) {
		presenter, prop := g.The("action.Source"), g.The("action.Context")
		if carrier := Carrier(prop); carrier != presenter {
			g.Say("You aren't holding", ArticleName(g, "action.Context", NameFullStop))
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"d87833ece051f1acbf636ab531b13708": func(g G.Play) {
		if PlayerLearns(g, "tookPaints") {
			g.The("evil fish").Says(`"Boy," says the fish, apparently to himself, "I sure hope that's some food she's finding for me in there. You know, the yummy food in the ORANGE CAN."`)
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"d87a6fc604d8a034efb8a192199a1a": func(g G.Play) {
		g.Say("That doesn't make much sense.")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"d95e07eb79eaf3fa2f3302d5ec0bf09": func(g G.Play) {
		if g.The("studio").Is("visited") {
			g.Say("Decorated with Britney's signature flair. It was her innate sense of style that first made you forgive her that ludicrous name. And here it is displayed to the fullest: deep-hued drapes on the walls, the windows flung open with their stunning view of old Vienna, the faint smell of coffee that clings to everything. Her easel stands over by the windows, where the light is brightest.")
		} else {
			g.Say(
				`This is Britney's studio. You haven't been around here for a while, because of how busy you've been with work, and she's made a few changes -- the aquarium in the corner, for instance. But it still brings back a certain emotional sweetness from the days when you had just met for the first time... when you used to spend hours on the sofa...
You shake your head. No time for fantasy. Must feed fish.`)
		}
		g.StopHere()
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/wearing.go:37
	"dba3f28c029c8ebf74fe7ac2d94e0ac0": func(g G.Play) {
		actor, prop := g.The("actor"), g.The("prop")
		if prop.Is("not wearable") {
			g.Say("That's not something you can wear.")
		} else {
			g.Go(Clothe("actor").With("prop"))
			g.Say("Now", actor.Text("name"), "is wearing the", prop.Text("name"))
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/core.go:180
	"dba636750406cd8d7bd0da3": func(g G.Play) {
		actor := g.The("actor")
		target := actor.Object("whereabouts")
		target.Go("report the view")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/kissing.go:13
	"dbc0e2eebefb53eb032ffd9e2de": func(g G.Play) { ReflectToTarget(g, "report kiss") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/eating.go:12
	"dc5f2a70d4f8b2f822f0d4f669c330b": func(g G.Play) { ReflectToTarget(g, "report eat") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/examining.go:17
	"dcddc49d44fe886a8f96d952f": func(g G.Play) {
		this := g.The("supporter")
		ListContents(g, "On the", this)
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"de68b28196dad365db8afd062982": func(g G.Play) {
		g.Say("No, you'd better leave it. It'd freak her out if you moved it.")
		g.StopHere()
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"deb7efe1826c5a9b5819f6073fbc089a": func(g G.Play) {
		if PlayerLearns(g, "examinedBouquet") {
			g.The("evil fish").Says(`"Oh, you shouldn't have," says the fish. "For me??"`)
			g.Say("You just respond with a livid glare.")
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"df1fbacb62e40599f5b43f7e52eb": func(g G.Play) {
		if PlayerLearns(g, "closedCabinet") && !g.The("fish food").Is("found") {
			g.The("evil fish").Says(`"Ooh, what do you think, Bob? I think we're going to have to dock the girl a few points. HAVE ANOTHER LOOK, sweetcakes, there's a doll."`)
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"dfc7b1de160d5dea8bc04e39ff06899c": func(g G.Play) {
		g.Say("Oh, it's tempting. But it would get you in a world of hurt later on.")
		g.StopHere()
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/attacking.go:12
	"e092832c33702a363a16fa5d40421bf": func(g G.Play) {
		if actor := g.The("actor"); g.The("player") == actor {
			g.Say("Violence isn't the answer.")
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/attacking.go:12
	"e1067e858bff37ff95ed58744c3655e0": func(g G.Play) { ReflectToTarget(g, "report attack") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/debugging.go:13
	"e27d852c08e661f57fed7a37fb44b6": func(g G.Play) {
		target := g.The("action.Target")
		contents := target.ObjectList("contents")
		g.Say("debugging contents of", target.Text("name"))
		for _, v := range contents {
			g.Say(v.Id().String())
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/extensions/conversation.go:107
	"e2d83166d1ece71900065119ffe42599": func(g G.Play) {
		player, talker, talkedTo := g.The("player"), g.The("action.Source"), g.The("action.Target")
		if player == talker {
			if quips := GetPlayerQuips(g); len(quips) == 0 {
				if Debugging {
					g.Log("! no conversation choices !")
				}
				player.Go("depart") // safety first
			} else {
				if Debugging {
					g.Log("!", talker, "printing", talkedTo, quips)
				}
				// FIX: the console should grab this to label the list, and add the header numbers./
				text := fmt.Sprintf("%s: ", player.Text("name"))
				g.Say(Lines("", text))
				for i, quip := range quips {
					cmt := quip.Text("comment")             // FIX: is this good? should it be slug, or name
					text := fmt.Sprintf("%d: %s", i+1, cmt) // FIX? template instead of fmt
					g.Say(text)                             // FIX FIX: CAN "SAY" TEXT BE SCOPED TO THE EVENT IN THE CMD OUTPUT.
				}
				player.IsNow("inputing dialog")
			}
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/showing.go:21
	"e460183c17f770fec808d132736f24e": func(g G.Play) {
		_, _, receiver := g.The("action.Source"), g.The("action.Target"), g.The("action.Context")
		receiver.Go("impress")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/stories.go:37
	"e549441d59279c1801dcd222a57f": func(g G.Play) {
		story := g.The("story")
		turnCount := story.Num("turn count") + 1
		story.SetNum("turn count", turnCount)
		//
		if story.Is("scored") {
			score := story.Num("score")
			status := fmt.Sprintf("%d/%d", int(score), int(turnCount))
			g.The("status bar").SetText("right", status)
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/devices.go:36
	"e54bc28e4a6d6ee8f15cde43d041f41": func(g G.Play) {
		g.Say("It's already off.") //[regarding the noun]?
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/examining.go:17
	"e5d3985209c83582a8cec2119c5f808": func(g G.Play) {
		object := g.The("object")
		desc := object.Text("description")
		if desc != "" {
			g.Say(desc)
		} else {
			//g.Say("You see nothing special about:")
			object.Go("print name")
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/insertingInto.go:49
	"e74361077dff040e0f5ac2c": func(g G.Play) {
		actor, prop := g.The("action.Source"), g.The("action.Context")
		if carrier := Carrier(prop); carrier != actor {
			g.Say("You aren't holding", ArticleName(g, "action.Context", NameFullStop))
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/insertingInto.go:49
	"e7456c6895a9da6aa": func(g G.Play) { ReflectWithContext(g, "receive insertion") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"e81a1f6706715fbda9c4a073ca7ef98": func(g G.Play) {
		g.Say("The fish swims adroitly out of range of your bare hand.")
		g.The("evil fish").Says(`"Hey," he says, and the bubbles of his breath brush against your fingers. "Count yourself lucky I don't bite you right now, you stinking mammal."`)
		g.StopHere()
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/insertingInto.go:49
	"e8b94915f1dace4e2b2f59f3a9162c": func(g G.Play) {
		g.Go(Insert("action.Source").Into("action.Context"))
		g.Say("You insert", ArticleName(g, "action.Source", nil), "into", ArticleName(g, "action.Context", NameFullStop))
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/extensions/conversation.go:107
	"e8ea14871dd41a378b7d216c974bfdd0": func(g G.Play) {
		talker, quip := g.The("actor"), g.The("quip")
		if Debugging {
			g.Log("!", talker, "discussing", quip)
		}
		if reply := quip.Text("reply"); reply != "" {
			talker.Says(reply)
		}
		PlayerMemory(g).Learns(quip)
		TheConversation(g).History().PushQuip(quip)
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/stories.go:37
	"e994180f6a5d75a96cbdd199f4385b7": func(g G.Play) {
		story := g.The("story")
		if story.Is("completed") {
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"e9cbdbcdcc4ceac8f98fc99adca7268": func(g G.Play) {
		g.Say(g.The("aquarium").Text("description"))
		g.StopHere()
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/showing.go:21
	"ea0713e6af256dfe532c004ff5a23f5": func(g G.Play) { ReflectWithContext(g, "report shown") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"ea5f6c77ffaf10b74b109e760": func(g G.Play) {
		if PlayerLearns(g, "smelledBouquet") {
			g.The("evil fish").Says(`"Mmm-mm," says the fish . Damn, I sure wish I had olfactory abilities. Hey, if I did, I might be even better at noticing the presence or absence of FOOD."`)
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/puttingItOn.go:35
	"eba487ee4a161c0ed3a629161b612e": func(g G.Play) { ReflectWithContext(g, "report placed") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"ec0f84674c17bbc66902dfe2b2ab284": func(g G.Play) {
		what := g.The("action.Target") // FIX: The("object") is player, not very nice...
		if what == g.The("evil fish") {
			// "increment the score"
			story := g.Our("A Day For Fresh Sushi")
			story.SetNum("score", story.Num("score")+1)
			// story trailer text
			g.Say("Triumphantly, you dump the remaining contents of the canister of fish food into the tank. It floats on the surface like scum, but the fish for once stops jawing and starts eating. Like a normal fish. Blub, blub.")
			g.Say("*** TWO HOURS LATER ***")
			g.Our("Britney").Says(`"So," Britney says, tucking a strand of hair behind your ear, "where shall we go for dinner? Since I made the big bucks on this trip, it's my treat. Anywhere you like."`)
			g.Our("Player").Says(`"I've had a hankering all day," you admit, as the two of you turn from the shuttle platform and head toward the bank of taxis. "I could really go for some sashimi right now."`)
			// "end the story finally"
			story.Go("end the story")
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/core.go:180
	"ec75951720bb443a92e928bf812c3": func(g G.Play) {
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
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/searching.go:13
	"ed97f386ce027a2ff676ef7a8a302e1": func(g G.Play) { ReflectToTarget(g, "report search") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/core.go:180
	"eecdc7b5555289906c2b415fb201cb9b": func(g G.Play) {
		source, actor := g.The("action.Source"), g.The("action.Target")
		if g.The("player") == actor {
			g.Say("You find nothing of interest.")
		} else {
			g.Say(actor.Text("Name"), "looks under the", source.Text("Name"), ".")
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/devices.go:36
	"ef7b600f685b87fcd10d15cb21a3": func(g G.Play) {
		device, actor := g.The("device"), g.The("actor")
		if device.Is("switched off") {
			device.Go("report already off", actor)
		} else {
			device.IsNow("switched off")
			device.Go("report now off", actor)
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/debugging.go:13
	"efb794467ebba90fb9d54e3f0332632e": func(g G.Play) {
		room := g.The("player").Object("whereabouts")
		contents := room.ObjectList("contents")
		g.Say("debugging contents of", room.Text("name"))
		for _, v := range contents {
			g.Say(v.Id().String())
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/extensions/conversation.go:107
	"efb9c6008baa5f297da22b6b88a69ae": func(g G.Play) { ReflectToTarget(g, "report comment") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"f12f24b77e1c05dfffb8569d4a8cb088": func(g G.Play) {
		prop := g.The("action.Target")
		if prop != g.The("bouquet") {
			g.The("evil fish").Says(`"Okay, so, what were you, raised in a barn? Normal folks like to use that for flowers. Or so I've observed."`)
			g.StopHere()
		} else {
			if PlayerLearns(g, "insertedFlowers") {
				g.Say("You settle the flowers into the vase and arrange them so that they look sprightly.")
				g.The("evil fish").Says(`"Oooh," says the fish. "No one ever changes the plant life in HERE. It's the same seaw--"`)
				g.The("player").Says(`"Cut me a break and cork it," you reply tartly.`)
				// FIX: report inserted?
				g.Go(Insert("the bouquet").Into("vase"))
				g.StopHere()
			}
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/devices.go:36
	"f13522698fb406b905407ccdfec39": func(g G.Play) {
		device := g.The("device")
		g.Say("Now the", device.Text("Name"), "is on.")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/movement.go:96
	"f147fff0a14762b6527a8f71240cdf00": func(g G.Play) { ReflectToTarget(g, "be passed through") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/extensions/conversation.go:107
	"f1487317dd37aafc7fdc5f81a5296e88": func(g G.Play) {
		talker, quip := g.The("actor"), g.The("quip")
		comment := quip.Text("comment")
		if Debugging {
			g.Log("!", talker, "commenting", quip, "'"+comment+"'")
		}

		// the player wants to speak: probably has chosen a line of dialog from the menu
		if comment != "" {
			talker.Says(comment)
		}

		// an actor will reply to the comment.
		// they will do this via Converse() at the end of the turn.
		con := TheConversation(g)
		con.Queue().SetNextQuip(quip)
		// moved history push into the npc's discuss
		// which should happen (right) after this.
		// the conversation choices are determined by what the npc says...
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/openingClosing.go:25
	"f15ec508566495ef319bda0247c3ba": func(g G.Play) {
		g.Say("It's already opened.")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"f18bf4192989ca36ad6671e40b253546": func(g G.Play) {
		if PlayerLearns(g, "examinedGravel") {
			g.Say("The fish notices your gaze; makes a pathetic mime of trying to find little flakes of remaining food amongst the gravel.")
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/smelling.go:12
	"f1eb81a12f8e22e8c9964742f7a703": func(g G.Play) { ReflectToTarget(g, "report smell") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"f383a1073f44b75c602b1e7b404e": func(g G.Play) {
		g.The("evil fish").Says(`"HelLLLOOO," screams the fish. "Whatever happened to FEEDING MEEE?"`)
		g.StopHere()
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"f386909950f2484579711506a8baff3f": func(g G.Play) {
		g.Say(
			`You're on the run. You've got a million errands to do -- your apartment to get cleaned up, the fish to feed, lingerie to buy, Britney's shuttle to meet--
The fish. You almost forgot. And it's in the studio, halfway across town from anywhere else you have to do. Oh well, you'll just zip over, take care of it, and hop back on the El. This'll be over in no time.
Don't you just hate days where you wake up the wrong color?`)
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/core.go:180
	"f3a039df8a36cf28033c549d1d": func(g G.Play) {
		if text := ArticleName(g, "object", NameFullStop); len(text) > 0 {
			text = lang.Capitalize(text)
			g.Say(text)
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/examining.go:17
	"f3d54a0c2d53c10de008fb6ab276670f": func(g G.Play) { ReflectToTarget(g, "be examined") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/eating.go:12
	"f4d2a3f2bc4ff9069889a9299b": func(g G.Play) {
		if actor := g.The("actor"); g.The("player") == actor {
			g.Say("That's not something you can eat.")
		} else {
			actor.Go("impress")
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/examining.go:17
	"f5982891e7ebdc570e12599f257d70ea": func(g G.Play) {
		object := g.The("object")
		object.Go("print details")
		object.Go("print contents")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"f5a8b9fe48a05ba88d82829dc5e83": func(g G.Play) {
		g.The("evil fish").Says(`"So what's it of?" asks the fish, as you turn away. "She never asks if I want to see them, you know?"`)
		g.The("player").Says(`"Her mother," you respond without thinking.`)
		g.The("evil fish").Says(`"Yeah? Man. I never knew my mother. Eggs, that's the way to go."`)
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/stories.go:37
	"f6b153d8dfbb91bec15a99f842d2": func(g G.Play) {
		player := g.The("action.Target")
		room := g.The("action.Context")
		player.Set("whereabouts", room) // Now("player's whereabouts is $room")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"f6f9049ec302b58175ad90021534b902": func(g G.Play) {
		if PlayerLearns(g, "examinedBagOnce") {
			//204-208
			g.The("evil fish").Says(`"What's in THERE?" asks the fish. "Didja bring me take-out? I don't mind Chinese. They eat a lot of carp, but what do I care? I'm not a carp. Live and let live is what I s--"`)
			g.The("player").Says(`"It's NOT take-out." You stare the fish down and for once he actually backstrokes a stroke or two. "It's PRIVATE."`)
		} else if PlayerLearns(g, "examinedBagTwice") {
			// 209-211
			g.The("evil fish").Says(`"If it's not take-out, I don't see the relevance!" shouts the fish. "Food is what you want in this situation. Food for MEEEE."`)
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/_examples/stories/fishy.go:457
	"f7de12be85b890ff8a18141f91c5c976": func(g G.Play) {
		g.The("evil fish").Says(`"That there is MY PA," says the fish, pointing at the scaley triton figure with one fin.`)
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/banner.go:17
	"f80bca4a7a24cc9ee42067363ecda08": func(g G.Play) {
		story := g.The("story")
		name := story.Text("name")
		headline := story.Text("headline")
		if headline == "" {
			headline = "An Interactive fiction" // FIX: default for headline in class.
		}
		author := story.Text("author")
		g.Say(name)
		g.Say(headline, "by", author)
		g.Say(VersionString)
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/giving.go:48
	"f820558419c997630a2044117eacd4b9": func(g G.Play) {
		g.The("action.Context").Go("impress")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/kissing.go:13
	"f9baafd0c928bb023c1ae00466c0464f": func(g G.Play) {
		source, target := g.The("action.Source"), g.The("action.Target")
		if source == target {
			g.Say(source.Text("Name"), "didn't get much from that.")
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/core.go:180
	"fa0bacce6442c0445e80e1fea52eb2": func(g G.Play) {
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
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/listening.go:12
	"fa673c29457c2a9ac47c49658685": func(g G.Play) { ReflectToTarget(g, "report listen") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/core.go:180
	"fbe785a6bfa575e678ea90da0211546": func(g G.Play) {
		text := ArticleName(g, "door", func(obj G.IObject) (status string) {
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
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/listening.go:12
	"fce05a152a90e9075438ca19ce7c6": func(g G.Play) { ReflectToLocation(g, "report listen") },

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/jumping.go:12
	"fd97981f521bc27a06faf20b0d5053da": func(g G.Play) {
		actor := g.The("action.Target")
		// FIX? inform often, but not always, tests for trying silently,
		// "if the action is not silent" ...
		// seems... strange. why report if if its silent?
		if g.The("player") == actor {
			g.Say("You jump on the spot.")
		} else {
			g.Say(actor.Text("Name"), "jumps on the spot.")
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/movement.go:96
	"fdb1dc013d4b6490c6586654cf3": func(g G.Play) {
		actor := g.The("actor")
		actor.Says("It's closed.")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/devices.go:36
	"fdcb93cc0324bd63a08f6415d5f62e3": func(g G.Play) {
		g.Say("It's inoperable.")
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/smelling.go:12
	"fe0b565e4a5208b37419400dbbaea13": func(g G.Play) {
		actor := g.The("action.Target")
		if g.The("player") == actor {
			g.Say("You smell nothing unexpected.")
		} else {
			g.Say(actor.Text("Name"), "sniffs.")
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/movement.go:96
	"fe407e25d67b995683a4497a99da6e5e": func(g G.Play) {
		departingDoor, actor := g.The("door"), g.The("actor")
		if dest := departingDoor.Object("destination"); !dest.Exists() {
			if Debugging {
				g.Log("couldnt find destination")
			}
		} else if room := dest.Object("whereabouts"); !room.Exists() {
			if Debugging {
				g.Log("couldnt find whereabouts")
			}
		} else {
			if Debugging {
				g.Log("moving ", actor, " to ", room)
			}
			if departingDoor.Is("closed") {
				if departingDoor.Is("locked") {
					departingDoor.Go("report locked", actor)
				} else {
					departingDoor.Go("report currently closed", actor)
				}
			} else {
				// FIX: player property change?
				// at the very least a move action.
				g.Go(MoveThe(actor).ToThe(room))
				room.Go("report the view")
			}
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/puttingItOn.go:35
	"fe5e3c26bd18f5891ae038584f8d": func(g G.Play) {
		prop := g.The("action.Context")
		if worn := prop.Object("wearer"); worn.Exists() {
			g.Say("You can't put worn clothing.")
			// FIX: try taking off the noun
			g.StopHere()
		}
	},

	// /Users/ionous/Dev/go/src/github.com/ionous/sashimi/standard/insertingInto.go:49
	"ff8957b596f3ee3ccd61c07445b3e2c8": func(g G.Play) {
		container := g.The("container")
		if container.Is("closed") {
			g.Say(ArticleName(g, "container", nil), "is closed.")
			g.StopHere()
		}
	},
}
