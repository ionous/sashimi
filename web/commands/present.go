package commands

import (
	R "github.com/ionous/sashimi/runtime"
	"github.com/ionous/sashimi/standard"
	"log"
)

func present(game *R.Game, player R.ObjectAdapter, output *CommandOutput) {
	if where, ok := standard.Location(player); !ok {
		log.Println("couldnt locate the player")
	} else if obj, ok := game.FindObject(where.Name()); !ok {
		log.Println("unknown error finding location")
	} else {
		view := SerializeView(game.Model, obj.Id())
		output.NewCommand("present", view)
	}
}
