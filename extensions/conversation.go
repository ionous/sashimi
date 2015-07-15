package extensions

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard"
)

func init() {
	AddScript(func(s *Script) {
		s.The("kinds",
			Called("quips"),
			//Comment comes first in dialog, and is said by the player.
			//If there is no comment then it is considered “npc-directed”.
			//For instance, a greeting when the player selects an NPC.
			Have("comment", "text"),
			Have("speaker", "actor"),
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
			Have("speaker", "actor"),
			AreEither("immediate").Or("postponed"),
			AreEither("obligatory").Or("optional"),
		)
		s.The("kinds", Called("facts"),
			// FIX: interestingly, kinds should have names
			// having the same property as a parent class probably shouldnt be an error
			Have("summary", "text"))

		// this is over-specification --
		// we've got recollection after all.
		// FIX: we just need fast sorting.
		qh := QuipHistory{}

		s.The("actors",
			Have("greeting", "quip"))

		// we need this in case the npc doesnt have a greeting
		// ( otherwise: the current quip's speaker could always be used
		// plus or minus some issues about multiple speakers. )
		interlocutor := R.NullObject("interlocutor")

		s.The("actors",
			Can("greet").And("greeting").RequiresOne("actor"),
			To("greet", func(g G.Play) {
				greeter, greeted := g.The("action.Source"), g.The("action.Target")
				if greeter == g.The("player") {
					switch {
					case !interlocutor.Exists():
						interlocutor = greeted
						if greeting := greeted.Object("greeting"); greeting.Exists() {
							qh.PushQuip(greeting)
							QueueQuip(g, greeting)
						}
					case greeted == interlocutor:
						g.Say("You're already speaking to them!")
					default:
						g.Say("You're already speaking to someone!")
					}
				}
			}))

		s.The("actors",
			Can("depart").And("departing").RequiresNothing(),
			To("depart", func(g G.Play) {
				qh.ClearQuips()
				interlocutor = R.NullObject("interlocutor")
			}))

		s.The("stories",
			When("ending the turn").Always(func(g G.Play) {
				Converse(g, qh)
			}))

		s.The("actors",
			Can("discuss").And("discussing").RequiresOne("quip"),
			To("discuss", standard.ReflectToTarget("report discuss")))

		// FIX? This event shifting ... it makes sense in a strict way --
		// but its also verbose to type, and difficult to follow.
		// part of the issue is naming convension for sure.
		// research needed...
		s.The("quips",
			Can("report discuss").And("reporting discuss").RequiresOne("actor"),
			To("report discuss", func(g G.Play) {
				player, talker, quip := g.The("player"), g.The("actor"), g.The("quip")
				// the player wants to speak: probably has chosen a line of dialog from the menu
				if talker == player {
					comment := quip.Text("comment")
					player.Says(comment)
					qh.PushQuip(quip)
				}
				// an actor wants to reply to the quip that was discussed.
				// they will do this at the end of the turn.
				QueueQuip(g, quip)
			}))

		var displayedChoices []G.IObject
		s.The("actors",
			Can("print conversation choices").And("printing conversation choices").RequiresOne("actor"),
			To("print conversation choices", func(g G.Play) {
				player, talker, talkedTo := g.The("player"), g.The("action.Source"), g.The("action.Target")
				if player == talker {
					displayedChoices = nil
					for i, quip := range GetPlayerQuips(g, qh, talkedTo) {
						// FIX? template instead of fmt
						text := quip.Text("comment")
						// FIX FIX: CAN "SAY" TEXT BE SCOPED TO THE EVENT IN THE CMD OUTPUT.
						g.Say(fmt.Sprintf("%d: %s", i+1, text))
						displayedChoices = append(displayedChoices, quip)
					}

					// time to hack the parser
					standard.TheParser.CaptureInput(func(input string) (err error) {
						var choice int
						if _, e := fmt.Sscan(input, &choice); e != nil {
							err = fmt.Errorf("Please choose a number from the menu; input: %s", input)
						} else if choice < 1 || choice > len(displayedChoices) {
							err = fmt.Errorf("Please choose a number from the menu; number: %d of %d", choice, len(displayedChoices))
						} else {
							quip := displayedChoices[choice-1]
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
