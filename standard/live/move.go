package live

import (
	G "github.com/ionous/sashimi/game"
)

// FIX: like facts.Learn() convert to a game action: actor.Go("move to", dest)
func Move(obj string) MoveToPhrase {
	return MoveToPhrase{actor: obj}
}

func MoveThe(obj G.IObject) MoveToPhrase {
	return Move(string(obj.Id()))
}

func (move MoveToPhrase) ToThe(dest G.IObject) MovingPhrase {
	return move.To(string(dest.Id()))
}

func (move MoveToPhrase) To(dest string) MovingPhrase {
	move.dest = dest
	return MovingPhrase(move)
}

func (move MoveToPhrase) OutOfWorld() MovingPhrase {
	return MovingPhrase(move)
}

func (move MovingPhrase) Execute(g G.Play) {
	actor, dest := g.The(move.actor), g.The(move.dest)
	AssignTo(actor, "whereabouts", dest)
}

type moveData struct {
	actor, dest string
}
type MoveToPhrase moveData
type MovingPhrase moveData
