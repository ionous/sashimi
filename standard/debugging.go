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
					g.Say(target.Name(), "=>", "out of world")
				} else {
					g.Say(target.Name(), "=>", relation, parent.Name())
				}
			}))

		s.Execute("print direct parent",
			Matching("parent of {{something}}"))
	})
}
