package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func init() {
	// FIX: a special AddDebugScript? that only gets activated with special command line parameters?
	AddScript(func(s *Script) {
		s.The("actors",
			Can("debug direct parent").And("debugging direct parent").RequiresOne("object"),
			To("debug direct parent", func(g G.Play) {
				target := g.The("action.Target")
				parent, relation := target.ParentRelation()
				if relation == "" {
					g.Say(target.Text("Name"), "=>", "out of world")
				} else {
					g.Say(target.Text("Name"), "=>", relation, parent.Text("Name"))
				}
			}))
		// FIX: sometimes parent of -- matches unexpected objects
		// >parent of automat
		//	hall-automat-door => whereabouts main hallway
		s.The("actors",
			Can("debug contents").And("debugging contents").RequiresOne("object"),
			To("debug contents", func(g G.Play) {
				target := g.The("action.Target")
				contents := target.ObjectList("contents")
				g.Say("debugging contents of", target.Text("name"))
				for _, v := range contents {
					g.Say(v.Id().String())
				}
			}))
		s.The("actors",
			Can("debug room contents").And("debugging room contents").RequiresNothing(),
			To("debug room contents", func(g G.Play) {
				room := g.The("player").Object("whereabouts")
				contents := room.ObjectList("contents")
				g.Say("debugging contents of", room.Text("name"))
				for _, v := range contents {
					g.Say(v.Id().String())
				}
			}))
		s.Execute("debug direct parent",
			Matching("parent of {{something}}"))
		s.Execute("debug contents",
			Matching("contents of {{something}}"))

		s.Execute("debug room contents",
			Matching("contents of room"))
	})
}
