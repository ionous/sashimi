package internal

import (
	"container/list"
)

type ActionQueue struct {
	g    *Game // sads
	list *list.List
}

// FUTURE: do we need to pass Game here? could we have a future queue ( and promises ) as a completely separate system?
type Future interface {
	Run(g *Game) error
}

//
func NewActionQueue(g *Game) *ActionQueue {
	return &ActionQueue{g, list.New()}
}

//
func (q ActionQueue) Empty() bool {
	return q.list.Front() == nil
}

func (q ActionQueue) Enqueue(f Future) {
	q.list.PushBack(f)
	// fmt.Println(fmt.Sprintf("queuing %T future, %d", f, q.list.Len()))
}

//
func (q ActionQueue) QueueFuture(f Future) {
	//q.list.PushBack(f)
	// fmt.Println(fmt.Sprintf("queuing %T future, %d", f, q.list.Len()))
	f.Run(q.g)
}

//
func (q ActionQueue) Pop() (ret Future) {
	if el := q.list.Front(); el != nil {
		ret = q.list.Remove(el).(Future)
	}
	// fmt.Println(fmt.Sprintf("popped %T future, %d", ret, q.list.Len()))
	return ret
}

// this craziness exists to help unwind the very deep callstacks the events create
func (q *ActionQueue) ProcessActions(g *Game) (err error) {
	for !q.Empty() {
		r := q.Pop()
		r.Run(g)
	}
	return err
}
