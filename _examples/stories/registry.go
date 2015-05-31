//
// Package for example stories
//
package stories

import (
	"fmt"
	"github.com/ionous/sashimi/script"
)

//
// The list of all known stories provided by the examples.
//
var Stories = Registry{make(map[string]script.InitCallback), ""}

//
// see "Stories."
//
type Registry struct {
	reg      map[string]script.InitCallback
	selected string
}

//
// Add a story to the list of known stories.
//
func (this *Registry) Register(name string, cb script.InitCallback) {
	if _, exists := this.reg[name]; exists {
		panic(fmt.Sprintf("the story '%s' already exists", name))
	}
	this.reg[name] = cb
}

//
// Return the list all known stories.
//
func (this *Registry) List() (ret []string) {
	for key, _ := range this.reg {
		ret = append(ret, key)
	}
	return ret
}

//
// Add the named story to the global scripts.
//
func (this *Registry) Select(name string) bool {
	_, ok := this.reg[name]
	if ok {
		this.selected = name
	}
	return ok
}

// when the call to InitializeScript happens,
// inject the script selected via the story registry.
func init() {
	script.AddScript(func(s *script.Script) {
		if cb, ok := Stories.reg[Stories.selected]; ok {
			cb(s)
		}
	})
}
