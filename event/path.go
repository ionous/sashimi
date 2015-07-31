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
func (path *PathList) addPath(n ITarget) {
	if p, ok := n.Parent(); ok {
		path.addPath(p)
	}
	path.PushBack(n)
}

//
func (path PathList) String() string {
	arr := make([]string, 0, path.Len())
	for it := path.Front(); it != nil; it = it.Next() {
		arr = append(arr, fmt.Sprintf("`%s`", it.Value))
	}
	return "[" + strings.Join(arr, ",") + "]"
}
