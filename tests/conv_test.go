package tests

import (
	. "github.com/ionous/sashimi/extensions"
	R "github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

//
func TestTalk(t *testing.T) {
	s := InitScripts()
	Talk(s)
	if g, err := NewTestGame(t, s); assert.NoError(t, err) {
		qh := QuipHistory{}
		qh.Push(ident.MakeId("OldBoy # WhatsTheMatter"))
		qp := QuipPool(g.Objects)
		newQuip := qp[ident.MakeId("OldBoy # Later")]
		//
		lastQuip, ok := qp.MostRecent(qh)
		require.True(t, ok, "found last quip")
		//
		npcId := qp.Interlocutor(lastQuip)
		require.EqualValues(t, npcId, "AlienBoy")
		//
		repeats, _ := lastQuip.GetValue("Repeatable").(bool)
		require.False(t, repeats, "repeats")
		//
		newrepeats, _ := newQuip.GetValue("Repeatable").(bool)
		require.True(t, newrepeats, "new repeats")
		//
		rank := qp.FollowsRecently(qh, newQuip.Id())
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
	s := InitScripts()
	Talk(s)
	if g, err := NewTestGame(t, s); assert.NoError(t, err) {
		total := 2
		comments := 1
		repeats := 1
		VisitObjects(g.Objects, "Quips", func(q *R.GameObject) (done bool) {
			total--
			if q.GetValue("Comment").(string) != "" {
				comments--
			}
			if r, _ := q.GetValue("Repeatable").(bool); r {
				repeats--
			}
			return
		})
		assert.Equal(t, 0, total, "total quips")
		assert.Equal(t, 0, comments, "comment quips")
		assert.Equal(t, 0, repeats, "repeat quips")
	}
}

func Talk(s *Script) {
	s.The("actor", Called("the alien boy"), Exists())
	s.The("alien boy", Has("greeting", "OldBoy # WhatsTheMatter"))
	s.The("quip", Called("OldBoy # WhatsTheMatter"),
		Has("speaker", "alien boy"),
		Has("reply", `"You wouldn't happen to have a matter disrupter?" the Alien Boy asks you.`))

	s.The("quip", Called("OldBoy # Later"),
		Is("repeatable"),
		Has("speaker", "alien boy"),
		Has("comment", `"Oh, sorry," Alice says. "I'll be back."`))
}
