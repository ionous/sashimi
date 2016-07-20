package standard

import (
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/runtime"
	. "github.com/ionous/sashimi/script"
)

func init() {
	AddScript(func(s *Script) {
		s.The("actors",
			// future: maybe a named string for the slot?
			Can("save it").And("saving it").RequiresNothing(),
			To("save it", func(g G.Play) {
				if s, e := runtime.SaveGame(g); e != nil {
					g.Log("error", e.Error())
				} else {
					g.Say("saved", s)
				}
			}),
		)
		s.Execute("save it", Matching("save"))
	})
	// FIX: future runtime load, incl
	// listing and selection of save games.
}
