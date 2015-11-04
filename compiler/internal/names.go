package internal

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/ident"
)

// Provides collision detection for names, especially b/t enum choices and other properties,
// as well as handling for empty string names.
// FIX? In the future add value and source code?
// FIX? this would make an excellent verbose logging hook
type NameSource map[string]NameEntry
type NameEntry struct {
	value string
	ref   string
}

func NewNameSource() NameSource {
	res := make(NameSource)
	res[""] = NameEntry{value: "empty string", ref: ""}
	return res
}

type NameScope struct {
	names NameSource
	scope interface{}
}

//
func (this NameScope) addName(name string, value string) (curr ident.Id, err error) {
	return this.names.addName(this.scope, name, value, "")
}

//
func (this NameScope) addRef(name string, value string, src source.Code) (curr ident.Id, err error) {
	return this.names.addName(this.scope, name, value, string(src))
}

//
// wrap the name source in a scope
func (this NameSource) NewScope(scope interface{}) NameScope {
	ret := NameScope{this, scope}
	ret.addName("", "empty string")
	return ret
}

//
// to avoid collisions between instances and types, only one name from each can exist
func (this NameSource) addName(scope interface{}, name string, value string, code string) (id ident.Id, err error) {
	id = M.MakeStringId(name)
	//
	if scope != nil {
		name = fmt.Sprintf("%s.%s", scope, id)
	} else {
		name = id.String()
	}
	//
	if curr, ok := this[name]; !ok {

		curr = NameEntry{value, code}
		this[name] = curr
	} else if curr.value != value {
		err = fmt.Errorf("'%v' in use as '%v'('%s'); respecified as '%v'('%s').",
			name,
			curr.value, curr.ref,
			value, code)
	}
	return id, err
}
