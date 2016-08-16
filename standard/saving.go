package standard

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
	. "github.com/ionous/sashimi/standard/live"
)

func init() {
	AddScript(func(s *Script) {
		// future: maybe a named or enum for require? ( tho an arbitrary string would be  unsafe )
		// maybe an optional require?
		s.The("globals", Called("save-settings"), Exist())
		s.The("save-setting", Called("auto-save"), Exist())
		s.The("save-setting", Called("normal-save"), Exist())

		s.The("actors",
			Can("save via input").And("saving via input").RequiresNothing(),
			To("save via input", func(g G.Play) {
				g.The("actor").Go("save it", g.The("normal-save"))
			}),
			Can("autosave via input").And("autosaving via input").RequiresNothing(),
			To("autosave via input", func(g G.Play) {
				g.The("actor").Go("save it", g.The("auto-save"))
			}),
			Can("save it").And("saving it").RequiresOne("save-setting"),
			To("save it", func(g G.Play) {
				autoSave := g.The("save-setting").Equals(g.The("auto-save"))
				if s, e := runtime.SaveGame(g, autoSave); !IsNil(e) {
					g.Log("error", e)
				} else {
					g.Say("saved", s)
				}
			}),
		)
		s.Execute("save via input", Matching("save"))
		s.Execute("autosave via input", Matching("autosave"))
	})
	// FIX: future runtime load, incl
	// listing and selection of save games.
}
