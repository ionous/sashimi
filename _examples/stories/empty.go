package stories

import (
	. "github.com/ionous/mars/script"
	. "github.com/ionous/mars/std/script"
)

//
func An_Empty_Room() (s Script) {
	s.The("story",
		Called("testing"),
		Has("author", "me"),
		Has("headline", "extra extra"))
	s.The("room",
		Called("somewhere"),
		Has("description", "an empty room"),
	)
	s.The("player", Exists(), In("somewhere"))
	return
}
func init() {
	stories.Register("empty", An_Empty_Room())
}
