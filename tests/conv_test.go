package tests

import (
	. "github.com/ionous/sashimi/extensions"
	G "github.com/ionous/sashimi/game"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestVisit to ensure that we can visit the dialog and see their properties.
func TestVisit(t *testing.T) {
	s := TalkScript()
	if game, err := NewTestGame(t, s); assert.NoError(t, err) {
		total := 2
		comments := 1
		repeats := 1
		g := R.NewGameAdapter(game.Game)
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
	if game, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := R.NewGameAdapter(game.Game)
		qh := QuipHistory{}
		// FIX:  saying hello.
		qh.PushQuip(g.Our("OldBoy # WhatsTheMatter"))
		//
		lastQuip := qh.MostRecent(g)
		require.True(t, lastQuip.Exists(), "found last quip")
		//
		npc := lastQuip.Object("Speaker")
		require.EqualValues(t, npc.Id(), "AlienBoy")
		//
		repeats := lastQuip.Is("one time")
		require.True(t, repeats, "repeats")
	}
}

// TestDiscuss for the discuss event and the conversation queue.
func TestDiscuss(t *testing.T) {
	s := TalkScript()
	if game, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := R.NewGameAdapter(game.Game)
		boy := g.The("alien boy")
		// push in the queue
		boy.Go("discuss", boy.Object("greeting"))
		// talk to the boy
		Converse(g, QuipHistory{})
		// clear the test, and make sure the queue is empty.
		require.Equal(t, ResetQuipQueue(), 0)
		lines := game.FlushOutput()
		require.Len(t, lines, 1)
		require.Equal(t, `The Alien Boy: "You wouldn't happen to have a matter disrupter?"`, lines[0])
	}
}

// TestTalkQuips to test player quip generation.
// FIX: this will only be interesting if there are multiple possiblities
// both related to the alien boy, but not this convo, and to other npcs
func TestTalkQuips(t *testing.T) {
	s := TalkScript()
	if game, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := R.NewGameAdapter(game.Game)
		qh := QuipHistory{}
		// FIX: saying hello.
		qh.PushQuip(g.Our("OldBoy # WhatsTheMatter"))
		qp := GetQuipPool(g)
		// verify that "later" should show up after "whats the matter".
		later := g.Our("OldBoy # Later")
		after := qp.SpeakAfter(qh, later)
		require.True(t, after, "speak after")
		// verify it actually does
		list := qp.GetPlayerQuips(qh, g.Our("Alien Boy"))
		require.Len(t, list, 1)
		require.Equal(t, later, list[0])
	}
}

// TestTalkChoices to test conversation choice generation.
func TestTalkChoices(t *testing.T) {
	s := TalkScript()
	if game, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := R.NewGameAdapter(game.Game)
		if ab := g.The("alien boy"); assert.True(t, ab.Exists(), "found boy") {
			if player := g.The("player"); assert.True(t, player.Exists(), "found player") {
				player.Go("print conversation choices", ab)
				lines := game.FlushOutput()
				require.Len(t, lines, 1)
			}
		}
	}
}

// TestTalkChoose for simple testing of the menu choices.
func TestTalkChoose(t *testing.T) {
	s := TalkScript()
	if game, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := R.NewGameAdapter(game.Game)
		if ab := g.The("alien boy"); assert.True(t, ab.Exists(), "found boy") {
			if player := g.The("player"); assert.True(t, player.Exists(), "found player") {
				player.Go("print conversation choices", ab)
				lines := game.FlushOutput()
				require.Len(t, lines, 1)
				if err := game.RunInput("1"); assert.NoError(t, err, "handling menu") {
					lines := game.FlushOutput()
					require.Len(t, lines, 1)
					require.Equal(t, `player: "Oh, sorry," Alice says. "I'll be back."`, lines[0])
				}
			}
		}
	}
}

// TalkScript creates some dialog to test.
func TalkScript() *Script {
	s := InitScripts()
	s.The("actor", Called("The Alien Boy"), Exists())
	s.The("alien boy", Has("greeting", "OldBoy # WhatsTheMatter"))
	s.The("quip", Called("OldBoy # WhatsTheMatter"),
		Has("speaker", "alien boy"),
		Has("reply", `"You wouldn't happen to have a matter disrupter?"`))

	s.The("quip", Called("OldBoy # Later"),
		Is("repeatable"),
		Has("speaker", "alien boy"),
		Has("comment", `"Oh, sorry," Alice says. "I'll be back."`))
	return s
}
