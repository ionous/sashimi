package quips

import (
	"fmt"
	facts "github.com/ionous/sashimi/extension/facts/native"
	quips "github.com/ionous/sashimi/extension/quips/native"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/util/lang"
	"strings"
)

func Describe_Quips(s *Script) {
	// we derive topics and quips from facts, so we can recollect, prohibit, etc. equally.
	// a quip can optionally be bound to a single actor.
	s.The("facts", Called("quips"),
		Have("subject", "actor"),
		Have("comment", "text"),
		Have("reply", "text"),
		Have("topic", "kind"),
		Have("slug", "text"),
		AreEither("repeatable").Or("one time"),
		AreOneOf("important", "unimportant", "trivial").Usually("unimportant"),
		AreEither("restrictive").Or("unrestricted").Usually("unrestricted"),
	)

	// quips and topics can be used by any npc; even though, most often, its by one.
	// s.The("facts", Called("topics"),
	// 	Have("lede", "text"))

	// one topic has multiple quips; one quip can be used in multiple topics.
	// s.The("kinds",
	// 	Called("quip topics"),
	// 	Have("topic", "topic"),
	// 	Have("quip", "quip"))
	// s.The("topics", Have("quips", "quip list"))
	// s.The("quips", Have("topics", "topic list"))

	// one quip can have multiple following quips; one quip can follow various quips.
	s.The("kinds",
		Called("following quips"),
		Have("leading", "quip"),
		Have("following", "quip"),
		AreEither("indirectly following").Or("directly following"))

	s.The("kinds", Called("quip requirements"),
		Have("fact", "fact"),
		AreEither("permitted").Or("prohibited"),
		Have("quip", "quip"),
	)

	s.The("quip", Called("default greeting"),
		Has("comment", ""),
		Has("reply", ""))

	// default greeting help determine conversation when being clicked on.
	// its not completely necessary; topic-less quips fit all conversations.
	s.The("actors",
		Have("greeting", "quip"), // "topic"
		Can("greet").And("greeting").RequiresOne("actor"),
		To("greet", func(g G.Play) {
			// FIX/FUTURE - this is *very* interesting that the player actions run some reusable code
			// and that code -- which raises events -- can be overriden
			// like, maybe Execute is completely wrong.
			// it's vaguley from inform, but what if phrases just had actions directly
			// the events would come from whatever those raise, if any....
			// oto -- its nice to override them completely....
			g.Go(quips.Introduce("action.source").To("action.target").WithDefault())
		}),
		Can("be greeted by").And("being greeted by").RequiresOne("actor").AndOne("quip"),
		To("be greeted by", func(g G.Play) {
			c := quips.Converse(g)
			if npc := c.Actor().Object(); npc.Exists() {
				if npc == g.The("action.source") {
					g.Say("You're already speaking to them!")
				} else {
					g.Say("You're already speaking to someone!")
				}
			} else {
				c.Actor().SetObject(g.The("action.source"))
				// FIX: doesnt raise an error of any sor when we say go("mispelling"
				g.The("action.target").Go("comment", g.The("action.context"))
			}
		}))
	s.Execute("greet", Matching("talk to {{something}}").Or("t|talk|greet|ask {{something}}"))

	// conversations track the player's current discussion.
	s.The("kinds",
		Called("global conversations"),
		Have("actor", "actor"),
		//Have("topic", "topic"),
		Have("quip", "quip"),
		Have("current", "quip"),
		Have("parent", "quip"),
		Have("grandparent", "quip"))
	s.The("global conversation", Called("conversation"), Exists())

	s.The("actors",
		Can("depart").And("departing").RequiresNothing(),
		To("depart", func(g G.Play) {
			if c := quips.Converse(g); c.Conversing() {
				c.Reset()
				g.Say("(", lang.Capitalize(DefiniteName(g, "actor", nil)), "says goodbye.", ")")
			}
		}))

	s.The("stories",
		When("ending the turn").Always(func(g G.Play) {
			g.The("player").Go("print conversation choices")
		}))

	s.The("actors",
		Can("comment").And("commenting").RequiresOne("quip"),
		To("comment", func(g G.Play) { ReflectToTarget(g, "report comment") }),
		Can("discuss").And("discussing").RequiresOne("quip"),
		To("discuss", func(g G.Play) { ReflectToTarget(g, "be discussed") }),
	)

	s.The("quips",
		Can("report comment").And("reporting comment").RequiresOne("actor"),
		To("report comment", func(g G.Play) {
			// NOTE: commenting is always the player.
			talker, quip := g.The("actor"), g.The("quip")
			if comment := quip.Text("comment"); comment != "" {
				talker.Says(comment)
			}
			quip.Go("follow up with", g.The("actor"))
		}),
		Can("follow up with").And("following up").RequiresOne("actor"),
		To("follow up with", func(g G.Play) {
			if npc := quips.Converse(g).Actor().Object(); npc.Exists() {
				npc.Go("discuss", g.The("quip"))
			}
		}),
		Can("be discussed").And("being discussed").RequiresOne("actor"),
		To("be discussed", func(g G.Play) {
			talker, quip := g.The("actor"), g.The("quip")
			if reply := quip.Text("reply"); reply != "" {
				talker.Says(reply)
			}
			quips.Converse(g).History().PushQuip(quip)
			facts.PlayerMemory(g).Learns(quip)
		}))

	s.The("actors",
		Can("print conversation choices").And("printing conversation choices").RequiresNothing(),
		To("print conversation choices", func(g G.Play) {
			if quips.Converse(g).Conversing() {
				player, talker := g.The("player"), g.The("action.Source")
				if player == talker {
					if quips := quips.PlayerQuips(g); len(quips) == 0 {
						player.Go("depart") // safety first
					} else {
						// FIX: the console should grab this to label the list, and add the header numbers.
						text := fmt.Sprintf("%s: ", player.Get("name").Text())
						g.Say(Lines("", text))
						for i, quip := range quips {
							slug := quip.Get("slug").Text()
							if slug == "" {
								slug = quip.Get("comment").Text()
								lines := strings.Split(slug, lang.NewLine)
								if len(lines) > 0 {
									slug = lines[0]
								}
							}
							text := fmt.Sprintf("%d: %s", i+1, slug) // FIX? template instead of fmt
							g.Say(text)                              // FIX FIX: CAN "SAY" TEXT BE SCOPED TO THE EVENT IN THE CMD OUTPUT.
						}
						player.IsNow("inputing dialog")
					}
				}
			}
		}))

	s.The("actors",
		AreEither("not inputing dialog").Or("inputing dialog"))

	s.The("stories",
		When("parsing player input").Always(func(g G.Play) {
			story := g.The("story")
			if player := g.The("player"); player.Is("inputing dialog") {
				player.IsNow("not inputing dialog")
				if quips := quips.PlayerQuips(g); len(quips) == 0 {
					player.Go("depart")
				} else {
					input := g.The("story").Get("player input").Text()
					var choice int
					if _, e := fmt.Sscan(input, &choice); e != nil || choice < 1 || choice > len(quips) {
						g.Say("Please choose a number from the menu.")
						story.Go("print conversation choices")
					} else {
						quip := quips[choice-1]
						if Debugging {
							g.Log("!", player, "chose", quip)
						}
						player.Go("comment", quip)
					}
					g.StopHere()
				}
			}
		}))
}
