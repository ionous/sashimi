package standard

import (
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
	"sort"
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
				names := make([]string, 0, len(contents))
				for _, v := range contents {
					names = append(names, v.Id().String())
				}
				sort.Strings(names)
				for _, v := range names {
					g.Say(v)
				}
			}))

		s.The("actors",
			Can("debug save").And("debugging save").RequiresNothing(),
			To("debug save", func(g G.Play) {
				g.Say("saving...")
				// name := g.List("stories").Get(0).Object().Id().String() + ".sav"
				// if f, e := os.Create(name); e != nil {
				// 	g.Log("error creating save", name, e.Error())
				// } else {
				// 	g.Log("saving", name)
				// 	defer f.Close()
				// 	if e := runtime.DebugSave(g, f); e != nil {
				// 		g.Log("error saving", name, e.Error())
				// 	}
				// }
			}),
			// FUTURE: havent tried resync on client, some sort of refresh page thing based on event is needed.
			Can("debug load").And("debugging load").RequiresNothing(),
			To("debug load", func(g G.Play) {
				g.Say("loading...")
				// name := g.List("stories").Get(0).Object().Id().String() + ".sav"
				// if f, e := os.Open(name); e != nil {
				// 	g.Log("error opening", name, e.Error())
				// } else {
				// 	defer f.Close()
				// 	if e := runtime.DebugLoad(g, f); e != nil {
				// 		g.Log("error loading", name, e.Error())
				// 	}
				// }
			}))

		s.Execute("debug direct parent",
			Matching("parent of {{something}}"))

		s.Execute("debug contents",
			Matching("contents of {{something}}"))

		s.Execute("debug room contents",
			Matching("contents of room"))

		s.Execute("debug save",
			Matching("save"))

		s.Execute("debug load",
			Matching("load"))
	})
}
