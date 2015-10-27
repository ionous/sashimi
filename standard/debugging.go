package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

func init() {
	// FIX: a special AddDebugScript? that only gets activated with special command line parameters?
	AddScript(func(s *Script) {
		s.The("actors",
			Can("print direct parent").And("printing direct parent").RequiresOne("object"),
			To("print direct parent", func(g G.Play) {
				target := g.The("action.Target")
				parent, relation := DirectParent(target)
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
			Can("print contents").And("printing contents").RequiresOne("object"),
			To("print contents", func(g G.Play) {
				target := g.The("action.Target")
				contents := target.ObjectList("contents")
				g.Say("printing contents of", target.Text("name"))
				for _, v := range contents {
					g.Say(v.Id().String())
				}
			}))
		s.The("actors",
			Can("print room contents").And("printing room contents").RequiresNothing(),
			To("print room contents", func(g G.Play) {
				room := g.The("player").Object("whereabouts")
				contents := room.ObjectList("contents")
				g.Say("printing contents of", room.Text("name"))
				for _, v := range contents {
					g.Say(v.Id().String())
				}
			}))
		s.Execute("print direct parent",
			Matching("parent of {{something}}"))
		s.Execute("print contents",
			Matching("contents of {{something}}"))

		s.Execute("print room contents",
			Matching("contents of room"))
	})
}
