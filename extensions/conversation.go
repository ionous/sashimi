package extensions

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard"
	"reflect"
)

var Debugging bool = true

type Conversation struct {
	Interlocutor GosNilInterfacesAreAnnoying
	History      QuipHistory
	Memory       QuipMemory
	Queue        QuipQueue
}

func discuss(how, other string) IFragment {
	return NewFunctionFragment(func(b SubjectBlock) error {
		b.The("following quips",
			Table("following", "indirectly-following-property", "leading").Contains(
				b.Subject(), how, other))
		return nil
	})
}
func DirectlyFollows(other string) IFragment {
	return discuss("directly following", other)
}
func IndirectlyFollows(other string) IFragment {
	return discuss("directly following", other)
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
				con := g.Global("conversation").(*Conversation)
				greeter, greeted := g.The("action.Source"), g.The("action.Target")
				if greeter == g.The("player") && greeted.Exists() {
					if npc, ok := con.Interlocutor.Get(); !ok {
						if greeting := greeted.Object("greeting"); greeting.Exists() {
							con.Interlocutor.Set(greeted)
							greeted.Go("discuss", greeting)
						} else if npc == greeted {
							g.Say("You're already speaking to them!")
						} else {
							g.Say("You're already speaking to someone!")
						}
					}
				}
			}))
		s.Execute("greet", Matching("talk to {{something}}").Or("t|talk|greet|ask {{something}}"))

		s.Generate("conversation", reflect.TypeOf(Conversation{}))

		s.The("actors",
			Can("depart").And("departing").RequiresNothing(),
			To("depart", func(g G.Play) {
				con := g.Global("conversation").(*Conversation)
				if Debugging {
					fmt.Println("!", g.The("actor"), "departing", con.Interlocutor)
				}
				con.History.ClearQuips()
				con.Queue.ResetQuipQueue()
				if npc, ok := con.Interlocutor.Get(); ok {
					npc.Object("next quip").Remove()
					con.Interlocutor.Clear()
				}
			}))

		s.The("stories",
			When("ending the turn").Always(func(g G.Play) {
				con := g.Global("conversation").(*Conversation)
				if npc, ok := con.Interlocutor.Get(); ok {
					Converse(g)
					g.The("player").Go("print conversation choices", npc)
				}
			}))

		s.The("actors",
			Can("discuss").And("discussing").RequiresOne("quip"),
			To("discuss", standard.ReflectToTarget("report discuss")))

		s.The("quips",
			Can("report discuss").And("reporting discuss").RequiresOne("actor"),
			To("report discuss", func(g G.Play) {
				player, talker, quip := g.The("player"), g.The("actor"), g.The("quip")
				if Debugging {
					fmt.Println("!", talker, "discussing", quip)
				}
				// the player wants to speak: probably has chosen a line of dialog from the menu
				if talker == player {
					comment := quip.Text("comment")
					player.Says(comment)

				}
				// an actor wants to reply to the quip that was discussed.
				// they will do this via Converse() at the end of the turn.
				con := g.Global("conversation").(*Conversation)
				con.History.PushQuip(quip) // FIX: when to advance this...?
				con.Queue.QueueQuip(quip)
			}))

		s.The("actors",
			Can("print conversation choices").And("printing conversation choices").RequiresOne("actor"),
			To("print conversation choices", func(g G.Play) {
				player, talker, talkedTo := g.The("player"), g.The("action.Source"), g.The("action.Target")
				if player == talker {
					quips := GetPlayerQuips(g)
					if Debugging {
						fmt.Println(talker, "printing", talkedTo, quips)
					}
					for i, quip := range quips {
						cmt := quip.Text("comment")
						text := fmt.Sprintf("%d: %s", i+1, cmt) // FIX? template instead of fmt
						g.Say(text)                             // FIX FIX: CAN "SAY" TEXT BE SCOPED TO THE EVENT IN THE CMD OUTPUT.
					}

					// time to hack the parser
					standard.TheParser.CaptureInput(func(input string) (err error) {
						var choice int
						if _, e := fmt.Sscan(input, &choice); e != nil {
							err = fmt.Errorf("Please choose a number from the menu; input: %s", input)
						} else if choice < 1 || choice > len(quips) {
							err = fmt.Errorf("Please choose a number from the menu; number: %d of %d", choice, len(quips))
						} else {
							quip := quips[choice-1]
							if Debugging {
								fmt.Println("!", player, "chose", quip)
							}
							player.Go("discuss", quip)
						}
						return err
					})
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
			// then the immediate optional conversation doesn't cleared; it sticks in there until the player chooses some unrestrictive quip.
			// but: it's difficult to get immediate conversation assigned to a person who isnt the current interlocutor
			// because the shortcuts always refer to the current interlocutor.
			// the gap is likely an oversight.
			AreEither("planned").Or("casual"),
		)
		s.The("actors",
			// FIX? with pointers, it wouldnt be too difficult to have parts now
			// an auto-created association.
			Have("next quip", "next quip"))
	})
}
