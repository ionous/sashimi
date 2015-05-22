package event

import (
	"container/list"
	"fmt"
	"strings"
)

//
type PathList struct {
	*list.List
}

//
func NewPathTo(n ITarget) PathList {
	path := PathList{list.New()}
	path.addPath(n)
	return path
}

//
func (PathList) Cast(el *list.Element) ITarget {
	return el.Value.(ITarget)
}

//
func (this *PathList) addPath(n ITarget) {
	if p, ok := n.Parent(); ok {
		this.addPath(p)
	}
	this.PushBack(n)
}

//
func (this PathList) String() string {
	arr := make([]string, 0, this.Len())
	for it := this.Front(); it != nil; it = it.Next() {
		arr = append(arr, fmt.Sprintf("`%s`", it.Value))
	}
	return "[" + strings.Join(arr, ",") + "]"
}
