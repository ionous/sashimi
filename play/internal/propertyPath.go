package internal

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
	"strconv"
	"strings"
)

// PropertyPath is used for error logging to track the full path to some property
// ex. bob/inventory/items/4/name
// FIX FIX FIX: the storing of strings has overhead, and is a bit ambiguous.
// might be interesting if instead, each element ( object, gameList, gameValue )
// supported an interface for debug printing the chain ( since thats what we and nothing more )
// in most cases they could keep a source access parent pointer of the same interface.
// ( note that's the accessed parent starting from the script callback, not some absolute parent --
// ie. the name of an item could come via some class query ( all swords ) or bob's inventory
type PropertyPath struct {
	path []string
}

func NewPath(id ident.Id) PropertyPath {
	// OPT: not id.String() to help gopherjs
	return PropertyPath{[]string{string(id)}}
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
