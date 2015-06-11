package commands

import (
	M "github.com/ionous/sashimi/model"
	R "github.com/ionous/sashimi/runtime"
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
	objectId, propertyId M.StringId,
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
