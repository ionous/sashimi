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
	objectId, propertyId string,
	prev, value interface{},
) {
	if prev != value {
		if inst, e := this.game.Model.Instances.FindInstance(objectId); e != nil {
			log.Println("PropertyChanged, but couldnt find object", objectId, propertyId)
		} else if prop, ok := inst.Class().FindProperty(propertyId); !ok {
			log.Println("PropertyChanged, but couldnt find property", objectId, propertyId)
		} else {
			switch prop.(type) {
			case *M.NumProperty, *M.TextProperty:
				cmd := Dict{"obj": objectId, "prop": propertyId, "now": value}
				this.output.NewCommand("set", cmd)

			case *M.RelativeProperty:
				cmd := Dict{"obj": objectId, "now": value.(string)}
				this.output.NewCommand("relate", cmd)

			case *M.EnumProperty:
				cmd := Dict{"obj": objectId, "was": prev.(string), "now": value.(string)}
				this.output.NewCommand("choose", cmd)
			}
		}
	}
}
