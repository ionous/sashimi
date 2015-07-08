package app

import (
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/util/ident"
	"log"
)

type PropertyChangeHandler struct {
	game   *R.Game
	output *CommandOutput
}

//
// IPropertyChanged
//
func (this PropertyChangeHandler) PropertyChanged(
	objectId, propertyId ident.Id,
	prev, value interface{},
) {
	if prev != value {
		if gobj, ok := this.game.Objects[objectId]; !ok {
			log.Println("PropertyChanged, but couldnt find object", objectId, propertyId)
		} else if prop, ok := gobj.Class().PropertyById(propertyId); !ok {
			log.Println("PropertyChanged, but couldnt find property", objectId, propertyId)
		} else {
			this.output.propertyChanged(this.game, gobj, prop, prev, value)
		}
	}
}
