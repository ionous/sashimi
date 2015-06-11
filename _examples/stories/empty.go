package stories

import (
	. "github.com/ionous/sashimi/script"
)

//
func An_Empty_Room(s *Script) {
	s.The("story",
		Called("testing"),
		Has("author", "me"),
		Has("headline", "extra extra"))
	s.The("room",
		Called("somewhere"),
		Has("description", "an empty room"),
	)
}
func init() {
	stories.Register("empty", An_Empty_Room)
}
