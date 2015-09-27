package event

import "container/list"

//
//
type Queue struct {
	*list.List
}

//
func NewQueue() *Queue {
	return &Queue{list.New()}
}

//
func (q Queue) Enqueue(target ITarget, msg Message) {
	qd := Queued{target, &msg}
	q.PushBack(qd)
}

//
func (q Queue) Empty() bool {
	return q.Front() == nil
}

//
func (q Queue) Pop() (tgt ITarget, msg *Message) {
	if el := q.Front(); el != nil {
		qd := q.Remove(el).(Queued)
		tgt, msg = qd.target, qd.msg
	}
	return tgt, msg
}

//
// internal queued data
//
type Queued struct {
	target ITarget
	msg    *Message
}
