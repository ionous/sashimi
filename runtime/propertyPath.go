package runtime

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
	"strconv"
	"strings"
)

type PropertyPath struct {
	path []string
}

func NewPath(id ident.Id) PropertyPath {
	return PropertyPath{[]string{id.String()}}
}
func RawPath(s string) PropertyPath {
	return PropertyPath{[]string{s}}
}

func (p PropertyPath) Index(i int) PropertyPath {
	path := append(p.path, strconv.Itoa(i))
	return PropertyPath{path}
}

func (p PropertyPath) InvalidIndex(i int) PropertyPath {
	path := append(p.path, fmt.Sprintf("!%d", i))
	return PropertyPath{path}
}

func (p PropertyPath) Add(s string) PropertyPath {
	path := append(p.path, s)
	return PropertyPath{path}
}

func (p PropertyPath) String() string {
	return strings.Join(p.path, "/")
}
