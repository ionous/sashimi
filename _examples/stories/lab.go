package stories

import (
	. "github.com/ionous/sashimi/script"
)

//
func The_Lab(s *Script) {
	s.The("story",
		Called("testing"),
		Has("author", "me"),
		Has("headline", "extra extra"))

	s.The("room",
		Called("the lab"),
		Has("description", "an empty room"))

	s.The("actor",
		Called("player"), Exists(),
		In("the lab"),
	)

	s.The("supporter",
		In("the lab"),
		Called("the table"),
		Is("fixed in place"),
		Supports("the glass jar"))

	s.The("prop",
		Called("the axe"),
		Has("brief", "a very nice guitar."))

	s.The("player", Possesses("the axe"))

	s.The("actor",
		Called("lab assistant"),
		Has("description", "That's Darcy. Your pretty, capable lab assistant. Although, there is a question of commas."),
		Exists(),
		In("the lab"))

	s.The("container",
		Called("the glass jar"),
		Is("transparent", "closed").And("hinged"),
		Has("brief", "beaker with a lid."),
		Contains("the eye dropper"))

	s.The("props",
		Called("droppers"),
		Have("drops", "num"))

	s.The("dropper", Called("eye dropper"), Exists(), Has("drops", 5))
}

func init() {
	stories.Register("lab", The_Lab)
}
