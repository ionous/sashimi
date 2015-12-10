package internal

import (
	"container/list"
)

type ActionQueue struct {
	list *list.List
}

// FUTURE: do we need to pass Game here? could we have a future queue ( and promises ) as a completely separate system?
type Future interface {
	Run(*Game) error
}

//
func NewActionQueue() ActionQueue {
	return ActionQueue{list.New()}
}

//phrase G.RuntimePhrase
func (q ActionQueue) QueueFuture(f Future) {
	q.list.PushBack(f)
}

//
func (q ActionQueue) Empty() bool {
	return q.list.Front() == nil
}

//
func (q ActionQueue) PopFuture() (ret Future) {
	if el := q.list.Front(); el != nil {
		ret = q.list.Remove(el).(Future)
	}
	return ret
}
