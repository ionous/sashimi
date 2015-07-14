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

//
func TestDiscuss(t *testing.T) {
	s := TalkScript()
	if game, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := R.NewGameAdapter(game.Game)
		boy := g.The("alien boy")
		boy.Go("discuss", boy.Object("greeting"))
		lines := game.FlushOutput()
		require.Len(t, lines, 1)
		require.Equal(t, `The Alien Boy: "You wouldn't happen to have a matter disrupter?"`, lines[0])
	}
}

//
func TestTalk(t *testing.T) {
	s := TalkScript()
	if game, err := NewTestGame(t, s); assert.NoError(t, err) {
		g := R.NewGameAdapter(game.Game)
		qh := QuipHistory{}
		// FIX:  saying hello.
		qh.PushQuip(g.Our("OldBoy # WhatsTheMatter"))
		qp := GetQuipPool(g)
		newQuip := g.Our("OldBoy # Later")
		//
		lastQuip := qh.MostRecent(g)
		require.True(t, lastQuip.Exists(), "found last quip")
		//
		npc := lastQuip.Object("Speaker")
		require.EqualValues(t, npc.Id(), "AlienBoy")
		//
		repeats := lastQuip.Is("one time")
		require.True(t, repeats, "repeats")
		//
		newrepeats := newQuip.Is("repeatable")
		require.True(t, newrepeats, "new repeats")
		//
		rank := qp.FollowsRecently(qh, newQuip)
		require.Equal(t, -1, rank, "unrelated")
		//
		after := qp.SpeakAfter(qh, newQuip)
		require.True(t, after, "speak after")
		//
		list := qp.GetPlayerQuips(qh)
		require.Len(t, list, 1)
	}
}

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
