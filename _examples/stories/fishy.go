package stories

import (
	. "github.com/ionous/sashimi/extension/facts"
	facts "github.com/ionous/sashimi/extension/facts/native"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard"
)

// "A Day For Fresh Sushi"
// Ported from Inform7. Original story by Emily Short, used with permission.
// The original source for this story can be found in the Inform documentation
// ( http://inform7.com )
func A_Day_For_Fresh_Sushi(s *Script) {
	Describe_Facts(s)
	// 1-3
	s.The("story",
		Called("A Day For Fresh Sushi"),
		Has("author", "Emily Short"),
		Has("headline", "Your basic surreal gay fish romance"),
		Is("scored"),
	)
	//
	s.The("actor",
		Called("player"),
		Exists(),
		In("the studio"),
	)
	// 7
	s.The("room",
		Called("studio"),
		Has("printed name", "Studio"),
		When("printing details").Always(func(g G.Play) {
			if g.The("studio").Is("visited") {
				g.Say("Decorated with Britney's signature flair. It was her innate sense of style that first made you forgive her that ludicrous name. And here it is displayed to the fullest: deep-hued drapes on the walls, the windows flung open with their stunning view of old Vienna, the faint smell of coffee that clings to everything. Her easel stands over by the windows, where the light is brightest.")
			} else {
				g.Say(
					`This is Britney's studio. You haven't been around here for a while, because of how busy you've been with work, and she's made a few changes -- the aquarium in the corner, for instance. But it still brings back a certain emotional sweetness from the days when you had just met for the first time... when you used to spend hours on the sofa...
You shake your head. No time for fantasy. Must feed fish.`)
			}
			g.StopHere()
		}))

	//  11
	s.The("studio",
		When("reporting smell").Always(func(g G.Play) {
			g.The("evil fish").Says(`The evil fish notices you sniffing the air. "Vanilla Raspberry Roast," it remarks. "You really miss her, don't you."`)
			g.Say("You glance over, startled, but the fish's mouth is open in a piscine equivalent of a laugh. You stifle the urge to skewer the thing...")
			g.StopHere()
		}))
	// (16-18]
	s.The("player",
		When("jumping").Always(func(g G.Play) {
			g.The("evil fish").Says(`"Er," says the fish. "Does that, like, EVER help??"`)
			g.StopHere()
		}))
	// [19-21)
	// *** Instead of going nowhere:
	//	say "You can't leave until you've fed the fish. Otherwise, he'll complain, and you will never hear the end of it."
	// [22-23)
	s.The("container",
		Called("cabinet"), In("studio"),
		Is("hinged", "closed").And("fixed in place"),
		Has("brief", "A huge cabinet, in the guise of an armoire, stands between the windows."),
		Has("description", "Large, and with a bit of an Art Nouveau theme going on in the shape of the doors."),
	)
	// [24-25):
	s.The("cabinet", IsKnownAs("armoire"))
	// [26-28)
	s.The("cabinet",
		When("reporting look under").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "lookedUnderCabinet") {
				g.The("evil fish").Says(`"Dustbunnies," predicts the fish, with telling accuracy. It executes what for all the world looks like a fishy shudder. "Lemme tell you, one time I accidentally flopped outta the tank, and I was TWO HOURS on the floor with those things STARING ME IN THE NOSE. It was frightening."`)
				g.StopHere()
			}
		}))
	// [29-31) ">open cabinet":
	// 1st time. "There ya go," says the fish. "The girl is getting WARMER.”
	// 2nd time. "You open the cabinet, revealing some paints and a heap of cloths."
	s.The("cabinet",
		When("reporting now open").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "openedCabinet") {
				g.The("evil fish").Says(`"There ya go," says the fish. "The girl is getting WARMER."`)
				g.StopHere()
			}
		}))
	// [32-34)
	s.The("cabinet",
		When("reporting now closed").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "closedCabinet") && !g.The("fish food").Is("found") {
				g.The("evil fish").Says(`"Ooh, what do you think, Bob? I think we're going to have to dock the girl a few points. HAVE ANOTHER LOOK, sweetcakes, there's a doll."`)
				g.StopHere()
			}
		}))
	// 35a
	s.The("cabinet",
		Contains("some paints"),
		Contains("some cloths"))
	// [35b-36)
	s.The("prop", Called("paints"),
		Has("description", "A bunch of tubes of oil paint, most of them in some state of grunginess, some with the tops twisted partway off."),
		Is("plural-named"), //articles, and plurals named items should be used as hints, but arent currently.
	)
	// [37-39)
	s.The("paints",
		When("reporting take").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "tookPaints") {
				g.The("evil fish").Says(`"Boy," says the fish, apparently to himself, "I sure hope that's some food she's finding for me in there. You know, the yummy food in the ORANGE CAN."`)
			}
		}))
	// [40-42)
	s.The("paints",
		After("being examined").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "examinedPaints") {
				g.The("evil fish").Says(`"Tons of useful stuff in there," hollers in the fish, in a syncopated burble.`)
			}
		}))
	// [43-45]
	s.The("prop", Called("cloths"),
		Has("description", "Various colors of drapery that Britney uses to set up backgrounds and clothe her models. She does a lot of portraiture, so this comes in handy. It's all a big messy wad at the moment. Organized is not her middle name."),
		IsKnownAs("drapery").And("cloth"),
		Is("plural-named"), //articles, and plurals named items should be used as hints, but arent currently.
	)
	// 46: **** The indefinite article of the cloths is "a heap of".
	// [47-52)
	// [53-57)
	s.The("cloths",
		When("reporting search").
			Or("reporting look under").
			Always(func(g G.Play) {
			fishFood := g.The("fish food")
			if !fishFood.Is("found") {
				fishFood.IsNow("found")
				g.Go(GiveThe(fishFood).To("the player"))
				g.Say("Poking around the cloths reveals -- ha HA! -- a vehemently orange can of fish food.")
				g.StopHere()
			}
		}),
		// FIX: i like the fact the event filters can
		When("reporting shown").Always(func(g G.Play) {
			fish := g.The("evil fish")
			receiver := g.The("action.Context")
			if fish == receiver {
				fish.Says("What are you, some kind of sadist? I don't want to see a bunch of cloths! What kind of f'ing good, 'scuse my French, is that supposed to do me? I don't even wear pants for God's sake!")
				g.Say("He really looks upset. You start wondering whether apoplexy is an ailment common to fish.")
				g.StopHere()
			}
		}),
	)
	// [58-60)
	s.The("cloths",
		After("being examined").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "examinedCloths") {
				g.The("evil fish").Says("Whatcha looking at? I can't see through the doors, you know.")
				g.StopHere()
			}
		}))
	// 61: FIX? instances in inform can have properties without clases
	// ( The fish food can be found or hidden. )
	s.The("containers", Called("canisters"),
		AreEither("hidden").Or("found"),
		// FIX: bother. can't add new constraints properly.
		// AreUsually("opaque")
	)
	s.The("canister", Called("fish food"),
		Has("description", "A vehemently orange canister of fish food."),
		Is("hidden").And("opaque"),
		IsKnownAs("can"),
		Is("closed").And("hingeless"),
		Is("plural-named"), // some fish food
	)
	// [63-56]
	s.The("fish food",
		When("reporting gave").Always(func(g G.Play) {
			receiver := g.The("action.Context")
			fish := g.The("evil fish")
			if fish == receiver {
				fish.Says("That's the ticket, sweetie! Bring it on.")
				g.StopHere()
			}
		}))
	// [69-75)
	s.The("fish food",
		When("being opened by").Always(func(g G.Play) {
			fish, player := g.The("evil fish"), g.The("player")
			fish.Says(`"Oh, for--!" The evil fish breaks out in exasperation and hives. "Screw the screwing around with the screwtop. SHE never has to do that."`)
			player.Says(`"Well, SHE is not here," you reply. "What do you suggest?"`)
			fish.Says(`">FEED FISH<" says the fish promptly, making fishy faces and pointing at you with his fin. "Simplicity. Try it."`)
			g.StopHere()
		}),
	)
	// [76-78)
	s.The("fish food",
		When("being inserted").Always(func(g G.Play) {
			g.The("evil fish").Says(`"HelLLLOOO," screams the fish. "Whatever happened to FEEDING MEEE?"`)
			g.StopHere()
		}),
	)
	// [79a]
	s.The("supporter", Called("easel"), In("the studio"), Exists())
	//[79b-80)
	s.The("easel", Is("scenery"), Supports("the painting"))
	s.The("painting", IsKnownAs("portrait").And("image"))
	//[81-88)
	s.The("prop", Called("painting"), Has("description",
		`Only partway finished, but you can tell what it is: Britney's mother. You only met the old woman once, before she faded out of existence in a little hospice in Salzburg.

In the picture, her hands are grasping tightly at a small grey bottle, the pills to which she became addicted in her old age, and strange, gargoyle-like forms clutch at her arms and whisper in her ears.

But the disturbing thing, the truly awful thing, is the small figure of Britney herself, down in the corner, unmistakable: she is walking away. Her back turned.

You thought she'd finally talked this out, but evidently not. Still feels guilty for leaving. You only barely stop yourself from tracing, with your finger, those tiny slumped shoulders...
`))
	//[89-90)
	s.The("painting",
		When("reporting take").Always(func(g G.Play) {
			g.Say("No, you'd better leave it. It'd freak her out if you moved it.")
			g.StopHere()
		}))
	//[91-98)
	s.The("painting",
		Before("being examined").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "examinedPainting") {
				g.The("evil fish").Says(`A ferocious banging from the aquarium attracts your attention as you go to look at the painting. "Hey!" screams the fish. "She doesn't like strangers looking at her paintings before they're DOONNNE!"`)
				g.The("player").Says(`"Shut up, you," you reply casually. "I'm not a stranger." But the fish puts you off a little bit, and your heart is already in your mouth before you see the painting itself...`)
			} else {
				g.Say("Once is really enough. It's pretty much embedded in your consciousness now.")
				g.StopHere()
			}
		}))
	//[99-105)
	// FIX: doing something is nicely broad; anyway to do the same?
	// NOTE: unlike the specificity rules mentioned above, here the "After doing something to the painting" is additive. It's not perfectly clear to me why. Possibly because "examining" and "something" are different sets of events. And yet, it only appears to trigger if the event -- examining -- was succesfully completed.
	s.The("painting",
		After("being examined").Always(func(g G.Play) {
			g.The("evil fish").Says(`"So what's it of?" asks the fish, as you turn away. "She never asks if I want to see them, you know?"`)
			g.The("player").Says(`"Her mother," you respond without thinking.`)
			g.The("evil fish").Says(`"Yeah? Man. I never knew my mother. Eggs, that's the way to go."`)
		}))
	s.The("opener", Called("window"), Is("scenery"), In("the studio"))
	s.The("window", Is("hinged").And("closed"))
	s.The("window", When("printing details").Always(func(g G.Play) {
		if g.The("window").Is("open") {
			g.Say(`Through the windows you get a lovely view of the street outside. At the moment, the glass is thrown open, and a light breeze is blowing through.`)
		} else {
			g.Say(`Through the windows, you get a lovely view of the street outside -- the little fountain on the corner, the slightly dilapidated but nonetheless magnificent Jugendstil architecture of the facing building. The glass itself is shut, however.`)
		}
		g.StopHere()
	}))
	s.The("window", IsKnownAs("windows"))
	// [108-110)
	s.The("window",
		When("reporting now open").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "openedWindow") {
				g.The("evil fish").Says(`"Thank god some air," says the fish. "Man, it was getting hard to breathe in here." Two beats pass. "Oh wait."`)
				g.StopHere()
			}
		}))
	// [111-112)
	s.The("supporter", Called("table"), Is("scenery"), In("the studio"))
	s.The("container", Called("vase"), Is("open").And("hingeless"))
	// NOTE: interstingly: because the table is scenery, once we take the vase: it becomes invisible; this is just like the original story.
	s.The("table", Supports("vase"))
	// [113-114)
	s.The("table", Has("description", "A monstrosity of poor taste and bad design: made of some heavy, French-empire sort of wood, with a single pillar for a central leg, carved in the image of Poseidon surrounded by nymphs. It's all scaley, and whenever you sit down, the trident has a tendency to stab you in the knee. But Britney assures you it's worth a fortune."))
	// [115-116)
	s.The("vase", Has("description", "A huge vase -- what you saw once described in a Regency romance as an epergne, maybe -- something so big that it would block someone sitting at the table from seeing anyone else also sitting at the table. But it does function nicely as a receptacle for hugeass bouquets of flowers."))
	// [117-
	s.The("table",
		When("reporting look under").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "lookedUnderTable") {
				g.The("evil fish").Says(`"You're not going to find anything down there," whines the fish. "I mean, c'mon. It's the fricking floor. Please tell me you can see that. I can see that. I'm a myopic fish in a tank ten feet away and I can tell you there is nothing there but floor."`)
				g.StopHere()
			}
		}))
	// [120-122):
	s.The("table",
		After("being examined").Always(func(g G.Play) {
			g.The("evil fish").Says(`"That there is MY PA," says the fish, pointing at the scaley triton figure with one fin.`)
		}))
	//[123-132)
	s.The("vase",
		When("receiving insertion").Always(func(g G.Play) {
			prop := g.The("action.Target")
			if prop != g.The("bouquet") {
				g.The("evil fish").Says(`"Okay, so, what were you, raised in a barn? Normal folks like to use that for flowers. Or so I've observed."`)
				g.StopHere()
			} else {
				if facts.PlayerLearns(g, "insertedFlowers") {
					g.Say("You settle the flowers into the vase and arrange them so that they look sprightly.")
					g.The("evil fish").Says(`"Oooh," says the fish. "No one ever changes the plant life in HERE. It's the same seaw--"`)
					g.The("player").Says(`"Cut me a break and cork it," you reply tartly.`)
					// FIX: report inserted?
					g.Go(Insert("the bouquet").Into("vase"))
					g.StopHere()
				}
			}
		}))
	//[133-134)
	s.The("prop", Called("telegram"), Exists())
	s.The("prop", Called("bouquet"), Exists())
	s.The("prop", Called("lingerie bag"), Exists())
	s.The("prop", Called("chef hat"), Exists())
	s.The("player", Possesses("telegram").And("bouquet").And("lingerie bag"), Wears("chef hat"))
	// 135
	s.The("telegram", Has("description", "A telegram, apparently. And dated three days ago. TRIUMPH OURS STOP BACK SOON STOP BE SURE TO FEED FISH STOP"))
	s.The("telegram", IsKnownAs("yellow paper"))
	// [144-149)
	// [140-143);
	// NOTE: the original script says:
	//   1. "After examining the telegram for the first time:" and,
	//   2. "After examining the telegram:"
	// It appears -- unexpectedly to me  -- Inform's rules of specificity cause the second phrase only happens when the first phrase does not.
	s.The("telegram",
		After("being examined").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "examinedTelegraph") {
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
		}))
	// [150-151);
	s.The("chef hat", Has("description", "A big white chef hat of the kind worn by chefs. In this case, you. Just goes to show what a hurry you were in on the way out of the restaurant."))
	s.The("chef hat", IsKnownAs("big").And("white").And("chefs").And("chef's"))
	// 152a
	s.The("container",
		Called("aquarium"),
		Is("transparent").And("open").And("fixed in place"), // note, you can take the aquarium in the original story, but some of the fish lines dont make sense that way.
		In("the studio"))
	// 152b
	s.The("aquarium", Is("hingeless"),
		Has("brief", "In one corner of the room, a large aquarium bubbles in menacing fashion."), Has("description", "A very roomy aquarium, large enough to hold quite a variety of colorful sealife -- if any yet survived."),
		IsKnownAs("tank"),
	)
	s.The("aquarium", Contains("gravel"), Contains("seaweed"))
	// 154a
	s.The("prop", Called("gravel"), Has("description", "A lot of very small grey rocks."))
	s.The("gravel", IsKnownAs("little rocks"), Is("plural-named"))
	// 154b
	s.The("prop", Called("seaweed"), Has("description", "Fake plastic seaweed of the kind generally bought in stores for exactly this purpose."))
	s.The("seaweed", IsKnownAs("weed"), Is("plural-named"))
	// 156: The examine containers rule does nothing when examining the aquarium.
	// interestingly, in Inform this completely hides the gravel and the seaweed even from the room description.
	// I think that's got to be a bug ( maybe a change to inform since the story was written? )
	// because there's no way to see the contents of the tank at all.
	s.The("aquarium",
		When("being examined").Always(func(g G.Play) {
			g.Say(g.The("aquarium").Text("description"))
			g.StopHere()
		}))
	// [158-160)
	s.The("gravel",
		After("being examined").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "examinedGravel") {
				g.Say("The fish notices your gaze; makes a pathetic mime of trying to find little flakes of remaining food amongst the gravel.")
			}
		}))
	// [161-163)
	s.The("seaweed",
		After("being examined").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "examinedSeaweed") {
				g.The("evil fish").Says(`"Nice, hunh?" blubs the fish, taking a stabbing bite out of one just by way of demonstration. "Look so good I could eat it."`)
			}
		}))
	// 164
	s.The("animal", Called("evil fish"),
		// helper to bind to evil fish instead of "fish food"
		IsKnownAs("fish"),
		Has("description", "Even if you had had no prior experience with him, you would be able to see at a glance that this is an evil fish. From his sharkish nose to his razor fins, every inch of his compact body exudes hatred and danger."),
	)
	s.The("aquarium", Contains("evil fish"))
	s.The("evil fish", When("reporting take").Always(func(g G.Play) {
		g.Say("The fish swims adroitly out of range of your bare hand.")
		g.The("evil fish").Says(`"Hey," he says, and the bubbles of his breath brush against your fingers. "Count yourself lucky I don't bite you right now, you stinking mammal."`)
		g.StopHere()
	}))
	s.The("evil fish", When("reporting attack").Always(func(g G.Play) {
		g.Say("Oh, it's tempting. But it would get you in a world of hurt later on.")
		g.StopHere()
	}))
	s.The("evil fish", When("reporting kiss").Always(func(g G.Play) {
		g.Say("You're saving all your lovin for someone a lot cuddlier.")
		g.StopHere()
	}))
	s.The("evil fish",
		After("being examined").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "examinedFishOnce") {
				g.Say("The fish glares at you, as though to underline this point.")
			} else if facts.PlayerLearns(g, "examinedFishTwice") {
				g.The("evil fish").Says(`"If you're looking for signs of malnutrition," says the fish, "LOOK NO FURTHER!!" And it sucks in its gills until you can see its ribcage.`)
			}
		}))
	//181
	//An every turn rule:
	//choose a random row in the Table of Fish Banter;
	//say "[comment entry][paragraph break]".
	s.The("stories",
		When("ending the turn").Always(func(g G.Play) {
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
		}))
	// 193-194
	s.The("bouquet", Has("description", "Okay, so it's silly and sentimental and no doubt a waste of money, of which there is never really enough, but: you miss her. You've missed her since ten seconds after she stepped aboard the shuttle to Luna Prime, and when you saw these -- her favorites, pure golden tulips like springtime -- you had to have them."),
		IsKnownAs("flowers").And("tulip").And("tulips"))
	// 195:
	s.The("bouquet", After("being examined").Always(func(g G.Play) {
		if facts.PlayerLearns(g, "examinedBouquet") {
			g.The("evil fish").Says(`"Oh, you shouldn't have," says the fish. "For me??"`)
			g.Say("You just respond with a livid glare.")
		}
	}))
	s.The("bouquet",
		When("reporting smell").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "smelledBouquet") {
				g.The("evil fish").Says(`"Mmm-mm," says the fish . Damn, I sure wish I had olfactory abilities. Hey, if I did, I might be even better at noticing the presence or absence of FOOD."`)
				g.StopHere()
			}
		}))
	s.The("lingerie bag", Has("description", "You grant yourself the satisfaction of a little peek inside. You went with a pale, silky ivory this time -- it has that kind of sophisticated innocence, and it goes well with the purple of your skin. A small smirk of anticipation crosses your lips."))
	// 204-
	s.The("lingerie bag",
		After("being examined").Always(func(g G.Play) {
			if facts.PlayerLearns(g, "examinedBagOnce") {
				//204-208
				g.The("evil fish").Says(`"What's in THERE?" asks the fish. "Didja bring me take-out? I don't mind Chinese. They eat a lot of carp, but what do I care? I'm not a carp. Live and let live is what I s--"`)
				g.The("player").Says(`"It's NOT take-out." You stare the fish down and for once he actually backstrokes a stroke or two. "It's PRIVATE."`)
			} else if facts.PlayerLearns(g, "examinedBagTwice") {
				// 209-211
				g.The("evil fish").Says(`"If it's not take-out, I don't see the relevance!" shouts the fish. "Food is what you want in this situation. Food for MEEEE."`)
			}
		}))
	// 212
	s.The("actors",
		Can("feed it").And("feeding it").RequiresOne("object"),
		To("feed it", func(g G.Play) {
			g.Say("That doesn't make much sense.")
		}))
	s.Execute("feed it", Matching("feed {{something}}"))
	s.The("actors",
		// 218:
		Before("feeding it").Always(func(g G.Play) {
			player := g.The("player")
			if c := Carrier(g.The("fish food")); c != player {
				g.Say("You need the fish food first!")
				g.StopHere()
			}
		}),
		// 222-227
		When("feeding it").Always(func(g G.Play) {
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
		}),
		// 228-230
		When("feeding it").Always(func(g G.Play) {
			if g.The("object") == g.The("fish food") {
				g.The("actor").Go("feed it", g.The("evil fish"))
				g.StopHere()
			}
		}),
	)
	s.The("actor", Called("Britney"), Exists())
	// 231-233
	s.Our("A Day For Fresh Sushi",
		When("commencing").Always(func(g G.Play) {
			g.Say(
				`You're on the run. You've got a million errands to do -- your apartment to get cleaned up, the fish to feed, lingerie to buy, Britney's shuttle to meet--
The fish. You almost forgot. And it's in the studio, halfway across town from anywhere else you have to do. Oh well, you'll just zip over, take care of it, and hop back on the El. This'll be over in no time.
Don't you just hate days where you wake up the wrong color?`)
		}),
		Has("maximum score", 1),
	)
	// need first, x-times filter for events.
	s.The("facts", Table("name").
		Has("lookedUnderCabinet").
		And("openedCabinet").
		And("closedCabinet").
		And("tookPaints").
		And("examinedPaints").
		And("examinedCloths").
		And("examinedPainting").
		And("openedWindow").
		And("lookedUnderTable").
		And("insertedFlowers").
		And("examinedTelegraph").
		And("examinedGravel").
		And("examinedSeaweed").
		And("examinedFishOnce").
		And("examinedFishTwice").
		And("examinedBouquet").
		And("smelledBouquet").
		And("examinedBagOnce").
		And("examinedBagTwice"))
}

func init() {
	stories.Register("sushi", A_Day_For_Fresh_Sushi)
}
