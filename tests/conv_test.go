package tests

import (
	. "github.com/ionous/sashimi/extensions"
	G "github.com/ionous/sashimi/game"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestVisit to ensure that we can visit the dialog and see their properties.

func TestVisit(t *testing.T) {
	s := TalkScript()
	if test, err := NewTestGame(t, s); assert.NoError(t, err) {
		total, comments, repeats := 3, 2, 1
		g := R.NewGameAdapter(test.Game)
		g.Visit("quips", func(q G.IObject) (done bool) {
			total--
			if q.Text("comment") != "" {
				comments--
			}
			if q.Is("repeatable") {
				repeats--
			}
			return
		})
		assert.Equal(t, 0, total, "total quips")
		assert.Equal(t, 0, comments, "comment quips")
		assert.Equal(t, 0, repeats, "repeat quips")
	}
}

// TestQuipHistory to test history tracking.
// FIX FIX: add more quips and make sure rank and wrapping work properly.
func TestQuipHistory(t *testing.T) {
	s := TalkScript()
	if test, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := R.NewGameAdapter(test.Game)
		qh := QuipHistory{}
		// FIX:  saying hello.
		qh.PushQuip(g.Our("WhatsTheMatter"))
		//
		lastQuip := qh.MostRecent(g)
		require.True(t, lastQuip.Exists(), "found last quip")
		//
		npc := lastQuip.Object("subject")
		require.EqualValues(t, npc.Id(), "AlienBoy")
		//
		repeats := lastQuip.Is("one time")
		require.True(t, repeats, "repeats")
	}
}

// TestDiscuss for the discuss event and the conversation queue.
func xTestDiscuss(t *testing.T) {
	s := TalkScript()
	if test, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := R.NewGameAdapter(test.Game)
		if boy := g.The("alien boy"); assert.True(t, boy.Exists(), "found boy") {
			if player := g.The("player"); assert.True(t, player.Exists(), "found player") {
				if con, ok := g.Global("conversation"); assert.True(t, ok, "found conversation") {
					con := con.(*Conversation)

					// hijack the person we are trying to talk to
					g.Go(Introduce("player").To("alien boy").WithQuip(boy.Object("greeting")))
					require.Equal(t, boy.Object("next quip").Object("quip"), boy.Object("greeting"))
					con.Converse(g)
					// clear the test, and make sure the queue is empty.
					if lines, err := test.FlushOutput(); assert.NoError(t, err) {
						require.True(t, len(lines) > 1)
						require.Equal(t, `The Alien Boy: "You wouldn't happen to have a matter disrupter?"`, lines[0])
					}
				}
			}
		}
	}
}

// TestTalkQuips to test player quip generation.
// FIX: this will only be interesting if there are multiple possiblities
// both related to the alien boy, but not this convo, and to other npcs
func xTestTalkQuips(t *testing.T) {
	s := TalkScript()
	if test, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := R.NewGameAdapter(test.Game)
		if boy := g.The("alien boy"); assert.True(t, boy.Exists(), "found boy") {
			if player := g.The("player"); assert.True(t, player.Exists(), "found player") {
				if con, ok := g.Global("conversation"); assert.True(t, ok, "found conversation") {
					con := con.(*Conversation)

					con.Interlocutor.Set(boy)

					con.History.PushQuip(g.Our("WhatsTheMatter"))

					later := g.Our("Later")
					// verify that "later" should show up after "whats the matter".
					// require.True(t, qp.SpeakAfter(later))

					// verify that "anybody" should show up after "whats the matter".
					// anybody := g.Our("DoesAnybody")
					// require.True(t, qp.SpeakAfter(anybody))

					// verify it actually does
					list := GetPlayerQuips(g)
					require.Len(t, list, 2)
					require.Contains(t, list, later)

					// test the actual converation choices printed
					player.Go("print conversation choices", boy)
					if lines, err := test.FlushOutput(); assert.NoError(t, err) {
						require.True(t, len(lines) > 2)

						// test the selection
						if lines, err := test.RunInput("2"); assert.NoError(t, err, "handling menu") {
							require.Len(t, lines, 1)
							require.Contains(t, lines, `player: "Oh, sorry," Alice says. "I'll be back."`)
						}
					}
				}
			}
		}
	}
}

// TestDirectFollows to test player quip generation.
func xTestDirectFollows(t *testing.T) {
	s := InitScripts()
	s.The("actor", Called("The Alien Boy"), Exists())
	s.The("alien boy", Has("greeting", "Don't leave!"))

	s.The("quip", Called("Don't leave!"),
		Has("subject", "Alien Boy"),
		Has("reply", "Don't leave me!"),
		Is("restrictive"))

	s.The("quip", Called("We need help!"),
		Has("subject", "Alien Boy"),
		Has("comment", "We've got to go look for help."),
		Has("reply", "I don't think we should go anywhere."),
		Is("restrictive"),
		DirectlyFollows("Don't leave!"))

	s.The("quip", Called("Who else is there?"),
		Has("subject", "Alien Boy"),
		Has("comment", "Do you see anyone else?"),
		Has("reply", "I might... should I?"),
		Is("restrictive"),
		DirectlyFollows("We need help!"))

	s.The("quip", Called("Automat Goodbye"),
		Has("subject", "Alien Boy"),
		Has("comment", "I'm going to look around some more."))

	standard.Debugging = true
	if test, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := R.NewGameAdapter(test.Game)
		if boy := g.The("alien boy"); assert.True(t, boy.Exists(), "found boy") {
			if player := g.The("player"); assert.True(t, player.Exists(), "found player") {
				//
				g.Go(Introduce("player").To("alien boy").WithDefault())
				//
				if con, ok := g.Global("conversation"); assert.True(t, ok, "found conversation") {
					con := con.(*Conversation)

					con.Converse(g)
					latest := con.History.MostRecent(g)

					require.True(t, latest.Exists(), "should have most recent quip")
					require.True(t, latest.Is("restrictive"), "should be restrictive")

					if quips := GetPlayerQuips(g); assert.Len(t, quips, 1) {
						q := quips[0]
						require.EqualValues(t, "WeNeedHelp", q.Id())
					}
				}
			}
		}
	}
}

// TestFactFinding to verify facts and their associated conversation rules.
func TestFactFinding(t *testing.T) {
	s := InitScripts()

	s.The("facts",
		Table("name", "summary").
			Has("red", "as rage").
			And("blue", "as the sea").
			And("yellow", "as the sun").
			And("green", "as fresh spirulina"))

	s.The("quips",
		Table("name", "reply").
			Has("retort", "arg!").
			And("smoothie request", "yes, please!").
			And("orbital wonder", "what were the skies like when you were young?"))

	s.The("quip requirements",
		Table("fact", "permitted-property", "quip").
			Has("red", "permitted", "retort").
			And("red", "prohibited", "smoothie request").
			And("green", "permitted", "retort").
			And("yellow", "prohibited", "retort"))

	var qm QuipMemory
	if test, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := R.NewGameAdapter(test.Game)
		assert.True(t, qm.IsQuipAllowed(g, g.The("orbital wonder")))
		assert.True(t, qm.IsQuipAllowed(g, g.The("smoothie request")))
		assert.False(t, qm.IsQuipAllowed(g, g.The("retort")), "the retort should not be allowed")
		qm.Learn(g.The("red"))
		assert.True(t, qm.IsQuipAllowed(g, g.The("orbital wonder")))
		assert.False(t, qm.IsQuipAllowed(g, g.The("smoothie request")))
		assert.False(t, qm.IsQuipAllowed(g, g.The("retort")))
		qm.Learn(g.The("green"))
		assert.True(t, qm.IsQuipAllowed(g, g.The("orbital wonder")))
		assert.False(t, qm.IsQuipAllowed(g, g.The("smoothie request")))
		assert.True(t, qm.IsQuipAllowed(g, g.The("retort")))
		qm.Learn(g.The("yellow"))
		assert.True(t, qm.IsQuipAllowed(g, g.The("orbital wonder")))
		assert.False(t, qm.IsQuipAllowed(g, g.The("smoothie request")))
		assert.False(t, qm.IsQuipAllowed(g, g.The("retort")))
	}
}

// TalkScript creates some dialog to test.
func TalkScript() *Script {
	s := InitScripts()
	s.The("actor", Called("The Alien Boy"), Exists())
	s.The("alien boy", Has("greeting", "WhatsTheMatter"))

	s.The("quip", Called("WhatsTheMatter"),
		Has("subject", "alien boy"),
		Has("reply", `"You wouldn't happen to have a matter disrupter?"`))

	s.The("quip",
		Called("DoesAnybody"),
		Has("subject", "Alien Boy"),
		DirectlyFollows("WhatsTheMatter"),
		Has("comment", `"Does anybody?"`),
		Has("reply", `"Or,"asks the Alien Boy, "maybe a ray gun?"`))

	s.The("quip", Called("Later"),
		Is("repeatable"),
		Has("subject", "alien boy"),
		Has("comment", `"Oh, sorry," Alice says. "I'll be back."`))
	return s
}
