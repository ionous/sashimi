package standard

import (
	//	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

//
func init() {
	AddScript(func(s *Script) {
		s.The("globals",
			Called("status bar instances"),
			Have("left", "text"),
			Have("right", "text"))

		s.The("status bar instance",
			Called("status bar"),
			Exists())
	})
}
