package live

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/lang"
	"strings"
)

// when is the right time for functions versus callbacks?
func ListContents(g G.Play, header string, obj G.IObject) (printed bool) {
	// if something described which is not scenery is on the noun and something which is not the player is on the noun:
	// obviously a filterd callback, visitor, would be nice FilterList("contents", func() ... )
	// FIX: if something has scenery objets, they appear as contents,
	// but then the list is empty. ( ex. lab coat, but it might happen elsewhere )
	// we'd maybe need to know if something printed?
	if contents := obj.ObjectList("contents"); len(contents) > 0 {
		g.Say(header, obj.Text("Name"), "is:")
		for _, content := range contents {
			content.Go("print description")
		}
		printed = true
	}
	return printed
}
