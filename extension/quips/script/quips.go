package script

import (
	"github.com/ionous/sashimi/extension/facts"
	"github.com/ionous/sashimi/extension/quips"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard/live"
	"github.com/ionous/sashimi/util/lang"
	"strings"
)

func Describe_Quips(s *Script) {
	s.The("actors", AreEither("chatty").Or("reticent"))
	// we derive topics and quips from facts, so we can recollect, prohibit, etc. equally.
	// a quip can optionally be bound to a single actor.
	s.The("facts", Called("quips"),
		Have("subject", "actor"),
		Have("comment", "text"),
		Have("reply", "text"),
		Have("topic", "kind"),
		Have("slug", "text"),
		AreOneOf("important", "unimportant", "trivial", "departing").Usually("unimportant"),
		AreEither("restrictive").Or("unrestricted").Usually("unrestricted"))

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
		HasText("comment", T("")),
		HasText("reply", T("")))

	// default greeting help determine conversation when being clicked on.
	s.The("actors",
		Have("greeting", "quip"),
		Can("greet").And("greeting").RequiresOnly("actor"),
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
		Before("being greeted by").Always(func(g G.Play) {
			if g.The("action.source").Is("reticent") {
				g.Say("There's no response.")
				g.StopHere()
			}
		}),
		To("be greeted by", func(g G.Play) {
			c := quips.Converse(g)
			if npc := c.Actor().Object(); npc.Exists() {
				if npc.Equals(g.The("action.source")) {
					g.Say("You're already speaking to them!")
				} else {
					g.Say("You're already speaking to someone!")
				}
			} else {
				c.Actor().SetObject(g.The("action.source"))
				// FIX: doesnt raise an error of any sor when we say go("mispelling"
				quip := g.The("action.context")
				g.The("action.target").Go("comment", quip)
				c.Topic().SetObject(quip.Get("topic").Object())
			}
		}))
	s.Execute("greet", Matching("talk to {{something}}").Or("t|talk|greet|ask {{something}}"))

	// conversations track the player's current discussion.
	s.The("kinds",
		Called("global conversations"),
		Have("actor", "actor"),
		Have("quip", "quip"),
		Have("current", "quip"),
		Have("parent", "quip"),
		Have("grandparent", "quip"),
		Have("topic", "kind"))
	s.The("global conversation", Called("conversation"), Exists())

	s.The("actors",
		Can("depart").And("departing").RequiresNothing(),
		To("depart", func(g G.Play) {
			if c := quips.Converse(g); c.Conversing() {
				c.Reset()
				// HACK: this looks bad in alice,
				// hiding it via event scoping doesnt work so well
				// other events inside of this event get hidden
				// FIX: before and after actions dont really fall before/after the scope, and the probably should.
				//g.Say("(", lang.Capitalize(DefiniteName(g, "actor", nil)), "says goodbye.", ")")
			}
		}))

	s.The("stories",
		When("ending the turn").Always(func(g G.Play) {
			g.The("player").Go("print conversation choices")
		}))

	s.The("actors",
		Can("comment").And("commenting").RequiresOnly("quip"),
		To("comment", g.ReflectToTarget("report comment")),
		Can("discuss").And("discussing").RequiresOnly("quip"),
		To("discuss", g.ReflectToTarget("be discussed")),
	)

	s.The("quips",
		Can("report comment").And("reporting comment").RequiresOnly("actor"),
		To("report comment", func(g G.Play) {
			// NOTE: commenting is always the player.
			talker, quip := g.The("actor"), g.The("quip")
			if comment := quip.Text("comment"); len(comment) > 0 {
				talker.Says(comment)
			}
			quip.Go("follow up with", g.The("actor"))
		}),
		Can("follow up with").And("following up with").RequiresOnly("actor"),
		To("follow up with", func(g G.Play) {
			if npc := quips.Converse(g).Actor().Object(); npc.Exists() {
				npc.Go("discuss", g.The("quip"))
			}
		}),
		Can("be discussed").And("being discussed").RequiresOnly("actor"),
		To("be discussed", func(g G.Play) {
			talker, quip := g.The("actor"), g.The("quip")
			if reply := quip.Text("reply"); len(reply) > 0 {
				talker.Says(reply)
			}
			c := quips.Converse(g)
			c.History().PushQuip(quip)
			if topic := quip.Get("topic").Object(); !topic.Equals(g.The("player")) {
				c.Get("topic").SetObject(topic)
			}
			facts.PlayerMemory(g).Learns(quip)
			if quip.Is("departing") {
				talker.Go("depart")
			}
		}))

	s.The("actors",
		Can("print conversation choices").And("printing conversation choices").RequiresNothing(),
		To("print conversation choices", func(g G.Play) {
			if quips.Converse(g).Conversing() {
				player, talker := g.The("player"), g.The("action.Source")
				if player.Equals(talker) {
					if quips := quips.PlayerQuips(g); !(len(quips) > 0) {
						player.Go("depart") // player rejected candy
					} else {
						//FIX: my theory has been, that like the graphics display, the text display should be grabbing events. theres a balance here -- do you g.say the header for the player name generically? or only in the console? how much becomes "output" specific in any given event? maybe not too much, just occasionaly things like these menus. in which case: no special header.
						for _, quip := range quips {
							quip.Go("be offered")
						}
						player.IsNow("inputing dialog")
					}
				}
			}
		}))

	s.The("actors",
		AreEither("not inputing dialog").Or("inputing dialog"))

	s.The("kinds",
		Can("be offered").And("being offered").RequiresNothing())

	s.The("quips",
		When("being offered").Always(func(g G.Play) {
			quip := g.The("quip")
			slug := quip.Get("slug").Text()
			if !(len(slug) > 0) {
				slug = quip.Get("comment").Text()
				lines := strings.Split(slug, lang.NewLine)
				if len(lines) > 0 {
					slug = lines[0]
				}
			}
			// FIX? events cant carry text, so a quip has to "say" its programmatically determined "slug" text. if there were onions -- maybe with quips snippets? -- you could pre-generate the slug.
			g.Say(slug)
		}))

	s.The("stories",
		When("parsing player input").Always(func(g G.Play) {
			story := g.The("story")
			if player := g.The("player"); player.Is("inputing dialog") {
				player.IsNow("not inputing dialog")
				input := g.The("story").Get("player input").Text()
				if quip := g.The(input); !quip.Exists() {
					g.Say("Please choose a valid quip.")
					story.Go("print conversation choices")
				} else {
					player.Go("comment", quip)
				}
				g.StopHere()
			}
		}))
}
