package jsmem

import (
	"bytes"
	"encoding/json"
	"github.com/ionous/sashimi/compiler/ingest"
	"github.com/ionous/sashimi/metal"
	"github.com/ionous/sashimi/net/app"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/standard/framework"
	"strings"
)

type JsMem struct {
	memory  metal.ObjectValueMap
	session *app.PartialSession
}

// New creates a blank js session.
// main can contain the mc pairing.
func New(mc ingest.ModelCode) JsMem {
	//	model    *M.Model
	//	calls    call.MarkerStorage
	mem := make(metal.ObjectValueMap)
	return create(mc, mem)
}

// Restore uses the passed code and data as a starting point, and restores the saved json over top of it.
func Restore(mc ingest.ModelCode, saved string) (ret JsMem, err string) {
	mem := make(metal.ObjectValueMap)
	if e := json.NewDecoder(strings.NewReader(saved)).Decode(&mem); e != nil {
		err = e.Error()
	} else {
		ret = create(mc, mem)
	}
	return
}

// Snapshot returns a json-blob of the current game state
// ( the intended use is so the client cant save )
func (js *JsMem) Snapshot() (res string, err string) {
	buf := new(bytes.Buffer)
	if e := json.NewEncoder(buf).Encode(js.memory); e != nil {
		err = e.Error()
	} else {
		res = buf.String()
	}
	return
}

// Get mirrors http get: retrieving data from the game
func (js *JsMem) Get(path string) (resp string, err string) {
	if res, e := resource.FindResource(js.session, path); e != nil {
		err = e.Error()
	} else {
		doc := res.Query()
		if r, e := js.encode(doc); e != nil {
			err = e.Error()
		} else {
			resp = r
		}
	}
	return
}

// Post mirrors http post: sending commands for the turn
func (js *JsMem) Post(body string) (resp string, err string) {
	reader := strings.NewReader(body)
	if doc, e := js.session.Post(reader); e != nil {
		err = e.Error()
	} else if r, e := js.encode(doc); e != nil {
		err = e.Error()
	} else {
		resp = r
	}
	return
}

func create(mc ingest.ModelCode, mem metal.ObjectValueMap) JsMem {
	meta := metal.NewMetal(mc.Model, mem)
	// the command output (unfortunately) needs some kind of id
	// since we only have one session at a time, it doesnt matter what
	out := app.NewCommandOutput("gopherjs", meta, framework.NewStandardView(meta))
	sess, e := app.NewPartialSession(meta, mc.Code, nil, out)
	if e != nil {
		panic(e)
	}
	return JsMem{mem, sess}
}

func (js JsMem) encode(doc resource.Document) (resp string, err error) {
	if prettyBytes, e := json.Marshal(doc); e != nil {
		err = e
	} else {
		resp = string(prettyBytes)
	}
	return
}
