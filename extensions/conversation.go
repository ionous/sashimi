package extensions

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/util/lang"
)

type Conversation struct {
	Interlocutor GosNilInterfacesAreAnnoying
	History      QuipHistory
	Memory       QuipMemory
	Queue        QuipQueue
}

// FIX: replace  with player, go learn
// ALSO: if this were in the "fact" package, it could be: fact.Learn
// and maybe prop.Give?
func Learn(fact string) FactPhrase {
	return FactPhrase(fact)
}
func LearnThe(fact G.IObject) FactPhrase {
	return FactPhrase(fact.Id().String())
}

type FactPhrase string

func (fact FactPhrase) Execute(g G.Play) {
	if con, ok := g.Global("conversation"); ok {
		con := con.(*Conversation)
		con.Memory.Learn(g.The(string(fact)))
	}
}

func PlayerLearns(g G.Play, fact string) (newlyLearned bool) {
	if con, ok := g.Global("conversation"); ok {
		con := con.(*Conversation)
		quip := g.The(string(fact))
		if recollects := con.Memory.Recollects(quip); !recollects {
			con.Memory.Learn(quip)
			newlyLearned = true
		}
	}
	return newlyLearned
}

func PlayerRecollects(g G.Play, fact string) (okay bool) {
	if con, ok := g.Global("conversation"); ok {
		con := con.(*Conversation)
		okay = con.Memory.Recollects(g.The(fact))
	}
	return okay
}

func DirectlyFollows(other string) IFragment {
	return discuss("directly following", other)
}
func IndirectlyFollows(other string) IFragment {
	return discuss("indirectly following", other)
}

func IsPermittedBy(fact string) IFragment {
	return requires("permitted", fact)
}

func IsProhibitedBy(fact string) IFragment {
	return requires("prohibited", fact)
}

func discuss(how, other string) IFragment {
	// FIX: a way to change the orgin?
	return NewFunctionFragment(func(b SubjectBlock) error {
		b.The("following quips",
			Table("following", "indirectly-following-property", "leading").Has(
				b.Subject(), how, other))
		return nil
	})
}

func requires(how, fact string) IFragment {
	return NewFunctionFragment(func(b SubjectBlock) error {
		b.The("quip requirements",
			Table("fact", "permitted-property", "quip").
				Has(fact, how, b.Subject()))
		return nil
	})
}

func init() {
	AddScript(func(s *Script) {
		s.The("kinds",
			Called("quips"),
			//Comment comes first in dialog, and is said by the player.
			//If there is no comment then it is considered “npc-directed”.
			//For instance, a greeting when the player selects an NPC.
			Have("comment", "text"),
			Have("subject", "actor"),
			Have("reply", "text"),
			//Have("hook", "text"), // displayed on the menu
			//performative, informative, questioning: used for ask about, tell about, or simply state the quip name
			AreEither("repeatable").Or("one time").Usually("one time"),
			AreEither("restrictive").Or("unrestricted").Usually("unrestricted"),
		//really important, unimportant, ...: from my extension to add priority sorting
		)

		// FIX: data not kinds.
		s.The("kinds",
			Called("following quips"),
			Have("leading", "quip"),
			Have("following", "quip"),
			AreEither("indirectly following").Or("directly following"),
		)
		s.The("kinds",
			Called("pending quips"),
			Have("subject", "actor"),
			AreEither("immediate").Or("postponed"),
			AreEither("obligatory").Or("optional"),
		)

		s.The("actors",
			Have("greeting", "quip"))

		s.The("actors",
			Can("greet").And("greeting").RequiresOne("actor"),
			To("greet", func(g G.Play) {
				g.Go(Introduce("action.Source").To("action.Target").WithDefault())
			}))
		s.Execute("greet", Matching("talk to {{something}}").Or("t|talk|greet|ask {{something}}"))

		//		s.Generate("conversation", reflect.TypeOf(Conversation{}))

		s.The("actors",
			Can("depart").And("departing").RequiresNothing(),
			To("depart", func(g G.Play) {
				if con, ok := g.Global("conversation"); ok {
					con := con.(*Conversation)
					if npc := con.Depart(); npc.Exists() {
						if Debugging {
							g.Log("!", g.The("actor"), "departing", npc)
						}
						g.Say("(", lang.Capitalize(DefiniteName(g, "actor", nil)), "says goodbye.", ")")
					}
				}
			}))

		s.The("stories",
			When("ending the turn").Always(func(g G.Play) {
				if con, ok := g.Global("conversation"); ok {
					con := con.(*Conversation)
					con.Converse(g)
				}
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
				if con, ok := g.Global("conversation"); ok {
					con := con.(*Conversation)
					con.Queue.SetNextQuip(g, quip)
				}
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
				if con, ok := g.Global("conversation"); ok {
					con := con.(*Conversation)
					con.Memory.Learn(quip)
					con.History.PushQuip(quip)
				}
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

						// time to hack the parser
						panic("use player state instead")
						//TheParser.CaptureInput
						_ = func(input string) (err error) {
							var choice int
							if _, e := fmt.Sscan(input, &choice); e != nil {
								err = fmt.Errorf("Please choose a number from the menu; input: %s", input)
							} else if choice < 1 || choice > len(quips) {
								err = fmt.Errorf("Please choose a number from the menu; number: %d of %d", choice, len(quips))
							} else {
								quip := quips[choice-1]
								if Debugging {
									g.Log("!", player, "chose", quip)
								}
								player.Go("comment", quip)
							}
							return err
						}
					}
				}
			}))

		s.The("kinds",
			Called("next quips"),
			Have("quip", "quip"),
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
			AreEither("planned").Or("casual"),
		)
		s.The("actors",
			// FIX? with pointers, it wouldnt be too difficult to have parts now; an auto-created association.
			Have("next quip", "next quip"))
	})
}
