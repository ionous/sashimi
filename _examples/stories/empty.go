package stories

import (
	. "github.com/ionous/mars/core"
	. "github.com/ionous/mars/script"
	. "github.com/ionous/mars/std/script"
)

//
func An_Empty_Room() (s Script) {
	s.The("story",
		Called("The empty room"),
		HasText("author", T("me")),
		HasText("headline", T("extra extra")))
	s.The("room",
		Called("somewhere"),
		HasText("description", T("an empty room")),
	)
	s.The("player", Exists(), In("somewhere"))
	return
}
func init() {
	stories.Register("empty", An_Empty_Room())
}
