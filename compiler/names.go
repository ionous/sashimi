package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
)

// Provides collision detection for names, especially b/t enum choices and other properties,
// as well as handling for empty string names.
// FIX? In the future add value and source code?
// FIX? this would make an excellent verbose logging hook
type NameSource map[string]string

func NewNameSource() NameSource {
	res := make(NameSource)
	res[""] = "empty string"
	return res
}

type NameScope struct {
	names NameSource
	scope interface{}
}

//
func (this NameScope) addName(name string, value string) (curr M.StringId, err error) {
	return this.names.addName(this.scope, name, value)
}

//
// wrap the name source in a scope
func (this NameSource) newScope(scope interface{}) NameScope {
	ret := NameScope{this, scope}
	ret.addName("", "empty string")
	return ret
}

//
// to avoid collisions between instances and types, only one name from each can exist
func (this NameSource) addName(scope interface{}, name string, value string) (id M.StringId, err error) {
	id = M.MakeStringId(name)
	//
	if scope != nil {
		name = fmt.Sprintf("%s.%s", scope, id)
	} else {
		name = id.String()
	}
	//
	curr := this[name]
	switch curr {
	case "": // 'name' didnt exist so add it
		this[name], curr = value, value
	case value: // already setup, good to go.
		break
	default: // 'name' was set to something other than value
		err = fmt.Errorf("'%v' in use as '%v'; respecified as '%v'.", name, curr, value)
	}
	return id, err
}
