package stories

import (
	"github.com/ionous/mars/script"
	"github.com/ionous/sashimi/util/errutil"
)

func Select(name string) (*script.Script, bool) {
	s, okay := stories[name]
	return s, okay
}

func List() (ret []string) {
	return stories.List()
}

// The list of all known stories provided by the examples.
var stories = make(Registry)

// Registry
type Registry map[string]*script.Script

// Register  adds a story to the list of known stories.
func (r Registry) Register(name string, s script.Script) {
	if _, exists := r[name]; exists {
		panic(errutil.New("the story", name, "already exists"))
	}
	r[name] = &s
}

// List return the list all known stories.
func (r Registry) List() (ret []string) {
	for key, _ := range r {
		ret = append(ret, key)
	}
	return ret
}
