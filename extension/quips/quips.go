package quip

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/util/lang"
)

func Describe_Quips(s *Script) {
	// we derive topics and quips from facts, so we can recollect, prohibit, etc. equally.
	// a quip can optionally be bound to a single actor.
	s.The("facts", Called("quips"),
		Have("subject", "actor"),
		Have("comment", "text"),
		Have("reply", "text"),
		AreEither("repeatable").Or("one time"),
		AreOneOf("important", "unimportant", "trivial").Usually("unimportant"),
		AreEither("restrictive").Or("unrestricted").Usually("unrestricted"),
	)

	// quips and topics can be used by any npc; even though, most often, its by one.
	s.The("facts", Called("topics"),
		Have("lede", "text"))

	// one topic has multiple quips; one quip can be used in multiple topics.
	s.The("kinds",
		Called("quip topics"),
		Have("topic", "topic"),
		Have("quip", "quip"))

	// one quip can have multiple following quips; one quip can follow various quips.
	s.The("kinds",
		Called("following quips"),
		Have("leading", "quip"),
		Have("following", "quip"),
		//AreEither("indirectly following").Or("directly following"),
	)

	// default greeting help determine conversation when being clicked on.
	// its not completely necessary; topic-less quips fit all conversations.
	s.The("actors",
		Have("greeting", "topic"),
		Can("greet").And("greeting").RequiresOne("actor"),
		To("greet", func(g G.Play) {
			g.Go(Introduce("action.Source").To("action.Target").WithDefault())
		}))
	s.Execute("greet", Matching("talk to {{something}}").Or("t|talk|greet|ask {{something}}"))

	// conversations track the player's current discussion.
	s.The("kinds",
		Called("conversations"),
		Have("actor", "actor"),
		Have("topic", "topic"),
		Have("quip", "quip"))

	s.The("conversation", Called("conversation"), Exists())

	// note: currently storing both quips and facts in the player's recollection.
	// FIX: many-to-many doesnt exist so we are hacking the actors to allow for the player's recollection; traversing all kinds would be heavy: so just usting a flag.
	//
	// s.The("actors",
	// 	HaveMany("recollections", "kinds").
	// 		Implying("kinds", HaveMany("recallers", "actor")))
	s.The("facts", AreEither("recollected").Or("not recollected").Usually("not recollected"))

	s.The("actors",
		Can("depart").And("departing").RequiresNothing(),
		To("depart", func(g G.Play) {
			con := TheConversation(g)
			if npc := con.Depart(); npc.Exists() {
				if Debugging {
					g.Log("!", g.The("actor"), "departing", npc)
				}
				g.Say("(", lang.Capitalize(DefiniteName(g, "actor", nil)), "says goodbye.", ")")
			}
		}))

	s.The("stories",
		When("ending the turn").Always(func(g G.Play) {
			Converse(g)
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
			g.Go(SetTheNextQuip(quip))
			// moved history push into the npc's discuss
			// which should happen (right) after this.
			// the conversation choices are determined by what the npc says...
		}),
		Can("be discussed").And("being discussed").RequiresOne("actor"),
		To("be discussed", func(g G.Play) {
			talker, quip := g.The("actor"), g.The("quip")
			if Debugging {
				g.Log("!", talker, "discussing", quip)
			}
			if reply := quip.Text("reply"); reply != "" {
				talker.Says(reply)
			}
			PlayerMemory(g).Learns(quip)
			TheConversation(g).History().PushQuip(quip)
		}))

	s.The("actors",
		Can("print conversation choices").And("printing conversation choices").RequiresOne("actor"),
		To("print conversation choices", func(g G.Play) {
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
		}))

	s.The("actors", AreEither("not inputing dialog").Or("inputing dialog"))

	s.The("stories",
		When("parsing player input").Always(func(g G.Play) {
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
		}))

	// re: planned/casual:
	// ( dont love that the quip gets changed )
	//
	// these have the same meaning as "immediate obligatory" and "immediate optional".
	// casual quips lose their relevance whenever a new casual or planned quip is set,
	// and the moment after a player has spoken ( stop any planned casual follow-ups ).
	//
	// note: there is a gap in the original logic --
	// if the current quip is restrictive and if the person isnt the current interlocutor,
	// then the immediate optional conversation doesn't clear;
	// it sticks in there until the player chooses some unrestrictive quip.
	// but: it's difficult to get immediate conversation assigned to a person who isnt the current interlocutor: the shortcuts always refer to the current interlocutor.
	// the gap is likely an oversight.
	s.The("quips",
		AreEither("planned").Or("casual"))
	s.The("actors",
		Have("next quip", "quip"))

}
