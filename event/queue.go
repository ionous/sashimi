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
// helper for creating a message and enqueuing
func (this Queue) QueueEvent(target ITarget, name string, data interface{}) {
	msg := Message{Name: name, Data: data}
	this.Enqueue(target, msg)
}

//
func (this Queue) Enqueue(target ITarget, msg Message) {
	qd := Queued{target, &msg}
	this.PushBack(qd)
}

//
func (this Queue) Empty() bool {
	return this.Front() == nil
}

//
func (this Queue) Pop() (tgt ITarget, msg *Message) {
	if el := this.Front(); el != nil {
		qd := this.Remove(el).(Queued)
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
