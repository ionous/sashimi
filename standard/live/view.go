package live

import G "github.com/ionous/sashimi/game"

// FUTURE? interestingly, while we wouldnt be able to encode them without special work, the contents of the phrases are fixed: we could have After("reporting").Execute(Phrase). maybe "standard" phrases could put themselves in some sort of wrapper? around the model? tho not quite sure howd that work.
type ViewRoomPhrase struct {
	object string
}

func View(object string) ViewRoomPhrase {
	return ViewRoomPhrase{object}
}

func (p ViewRoomPhrase) Execute(g G.Play) {
	room := g.The("room")
	// sometines a blank like is printed without this
	// (maybe certain directions? or going through doors?)
	// not sure why, so leaving this for consistency
	g.Say(Lines("", room.Get("Name").Text()))
	g.Say(Lines(room.Get("description").Text(), ""))
	// FIX? uses 1 to exclude the player....
	// again, this happens because we dont know if print description actually did anything (re:scenery, etc.)
	if contents := room.ObjectList("contents"); len(contents) > 1 {
		for _, obj := range contents {
			obj.Go("print description")
			g.Say("")
		}
	}
}
