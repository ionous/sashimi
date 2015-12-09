package tests

import (
	. "github.com/ionous/sashimi/extension/facts"
	. "github.com/ionous/sashimi/extension/quips"
	quips "github.com/ionous/sashimi/extension/quips/native"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestQuipVisit to ensure that we can visit the dialog and see their properties.

func TestQuipVisit(t *testing.T) {
	s := TalkScript()
	if test, err := NewTestGame(t, s); assert.NoError(t, err) {
		total, comments, repeats := 3, 2, 1
		g := test.Game.NewAdapter()
		for i, quips := 0, g.List("quips"); i < quips.Len(); i++ {
			q := quips.Get(i).Object()
			if q.Text("comment") != "" {
				comments--
			}
			if q.Is("repeatable") {
				repeats--
			}
			total--
		}
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
		g := test.Game.NewAdapter()
		qh := quips.Converse(g).History()
		// FIX:  saying hello.
		qh.PushQuip(g.Our("WhatsTheMatter"))
		//
		lastQuip := qh.MostRecent()
		require.True(t, lastQuip.Exists(), "found last quip")
		//
		npc := lastQuip.Object("subject")
		require.EqualValues(t, "AlienBoy", npc.Id().String())
		//
		repeats := lastQuip.Is("one time")
		require.True(t, repeats, "repeats")
	}
}

// TestQuipDiscuss for the discuss event and the conversation queue.
func TestQuipDiscuss(t *testing.T) {
	standard.Debugging = true
	s := TalkScript()
	if test, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := test.Game.NewAdapter()
		if boy := g.The("alien boy"); assert.True(t, boy.Exists(), "found boy") {
			if player := g.The("player"); assert.True(t, player.Exists(), "found player") {
				if con := quips.Converse(g); assert.True(t, con.Exists(), "found conversation") {
					if greeting := boy.Object("greeting"); assert.True(t, greeting.Exists(), "found greeting") {
						player.Go("greet", boy)
						// clear the test, and make sure the queue is empty.
						if lines, err := test.FlushOutput(); assert.NoError(t, err) {
							require.Len(t, lines, 1)
							require.Equal(t, `Alien Boy: "You wouldn't happen to have a matter disrupter?"`, lines[0])
						}
					}
				}
			}
		}
	}
}

// TestQuipTalkQuips to test player quip generation.
// FIX: this will only be interesting if there are multiple possiblities
// both related to the alien boy, but not this convo, and to other npcs
func TestQuipTalkQuips(t *testing.T) {
	s := TalkScript()
	if test, e := NewTestGame(t, s); assert.NoError(t, e) {
		g := test.Game.NewAdapter()
		if boy := g.The("alien boy"); assert.True(t, boy.Exists(), "found boy") {
			if player := g.The("player"); assert.True(t, player.Exists(), "found player") {
				if con := quips.Converse(g); assert.True(t, con.Exists(), "found conversation") {
					// pretend we just said whats the matter.
					require.False(t, con.Conversing())
					con.Actor().SetObject(boy)
					con.History().PushQuip(g.Our("WhatsTheMatter"))
					require.True(t, con.Conversing())

					// we should have two comments in return.
					list := quips.PlayerQuips(g)
					require.Len(t, list, 2)
					require.Contains(t, list, g.Our("Later"))

					// test the actual converation choices printed
					player.Go("print conversation choices")
					if lines, e := test.FlushOutput(); assert.NoError(t, e) {
						require.True(t, len(lines) > 2)
						if stories := g.List("stories"); assert.Equal(t, 1, stories.Len()) {
							story := stories.Get(0).Object()

							story.Get("player input").SetText("2")
							if act, e := test.Game.QueueAction("parse player input", story.Id()); assert.NoError(t, e) {
								if lines, e := test.FlushOutput(); assert.NoError(t, e) {
									require.True(t, act.Cancelled())

									// test the selection
									// if lines, err := test.RunInput("2"); assert.NoError(t, err, "handling menu") {
									require.Len(t, lines, 1)
									require.Contains(t, lines, `player: "Oh, sorry," Alice says. "I'll be back."`)
								}
							}
						}
					}
				}
			}
		}
	}
}

// TestQuipDirectFollows to test player quip generation.
func TestQuipDirectFollows(t *testing.T) {
	s := InitScripts()
	Describe_Facts(s)
	Describe_Quips(s)

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

	// FIFIIFSIFSDFISI
	s.The("quip", Called("Automat Goodbye"),
		Has("subject", "Alien Boy"),
		Has("comment", "I'm going to look around some more."))

	standard.Debugging = true
	if test, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := test.Game.NewAdapter()
		if boy := g.The("alien boy"); assert.True(t, boy.Exists(), "found boy") {
			if player := g.The("player"); assert.True(t, player.Exists(), "found player") {
				player.Go("greet", boy)
				test.FlushOutput()
				if con := quips.Converse(g); assert.True(t, con.Conversing(), "found conversation") {
					if latest := con.History().MostRecent(); assert.True(t, latest.Exists(), "should have most recent quip") {
						require.True(t, latest.Is("restrictive"), "should be restrictive")
						if quips := quips.PlayerQuips(g); assert.Len(t, quips, 1) {
							q := quips[0]
							require.EqualValues(t, "WeNeedHelp", q.Id().String())
						}
					}
				}
			}
		}
	}
}

// TestQuipFactFinding to verify facts and their associated conversation rules.
func TestQuipFactFinding(t *testing.T) {
	s := InitScripts()
	Describe_Facts(s)
	Describe_Quips(s)

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
		Table("fact", "permitted", "quip").
			Has("red", "permitted", "retort").
			And("red", "prohibited", "smoothie request").
			And("green", "permitted", "retort").
			And("yellow", "prohibited", "retort"))

	if test, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := test.Game.NewAdapter()
		qm := quips.PlayerMemory(g)

		assert.True(t, qm.IsQuipAllowed(g.The("orbital wonder")))
		assert.True(t, qm.IsQuipAllowed(g.The("smoothie request")))
		assert.False(t, qm.IsQuipAllowed(g.The("retort")), "the retort should not be allowed")

		red := g.The("red")
		require.True(t, red.Exists())
		require.False(t, qm.Recollects(red), "before learning, doesnt contain red")
		qm.Learns(red)
		require.True(t, qm.Recollects(red), "after learning contains red")

		assert.True(t, qm.IsQuipAllowed(g.The("orbital wonder")))
		assert.False(t, qm.IsQuipAllowed(g.The("smoothie request")))
		assert.False(t, qm.IsQuipAllowed(g.The("retort")))
		green := g.The("green")
		qm.Learns(green)
		assert.True(t, qm.IsQuipAllowed(g.The("orbital wonder")))
		assert.False(t, qm.IsQuipAllowed(g.The("smoothie request")))
		assert.True(t, qm.IsQuipAllowed(g.The("retort")))
		yellow := g.The("yellow")
		qm.Learns(yellow)
		assert.True(t, qm.IsQuipAllowed(g.The("orbital wonder")))
		assert.False(t, qm.IsQuipAllowed(g.The("smoothie request")))
		assert.False(t, qm.IsQuipAllowed(g.The("retort")))
	}
}

// TalkScript creates some dialog to test.
func TalkScript() *Script {
	s := InitScripts()
	Describe_Facts(s)
	Describe_Quips(s)
	s.The("story", Called("testing"), Exists())
	s.The("actor", Called("The Alien Boy"), Exists())
	s.The("alien boy", Has("greeting", "WhatsTheMatter"))

	s.The("quip", Called("WhatsTheMatter"),
		Has("subject", "alien boy"),
		Is("one time"),
		Has("reply", `"You wouldn't happen to have a matter disrupter?"`))

	s.The("quip",
		Called("DoesAnybody"),
		Has("subject", "Alien Boy"),
		Is("one time"),
		DirectlyFollows("WhatsTheMatter"),
		Has("comment", `"Does anybody?"`),
		Has("reply", `"Or,"asks the Alien Boy, "maybe a ray gun?"`))

	s.The("quip", Called("Later"),
		Has("subject", "alien boy"),
		Has("comment", `"Oh, sorry," Alice says. "I'll be back."`))
	return s
}
