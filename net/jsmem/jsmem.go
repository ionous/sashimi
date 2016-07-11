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
	id      string
	memory  metal.ObjectValueMap
	session *app.PartialSession
}

// New creates a blank js session.
// main can contain the mc pairing.
func New(id string, mc ingest.ModelCode) JsMem {
	//	model    *M.Model
	//	calls    call.MarkerStorage
	mem := make(metal.ObjectValueMap)
	return create(id, mc, mem)
}

// Restore uses the passed code and data as a starting point, and restores the saved json over top of it.
func Restore(id string, mc ingest.ModelCode, saved string) (ret JsMem, err error) {
	mem := make(metal.ObjectValueMap)
	if e := mem.Load(strings.NewReader(saved)); e != nil {
		err = e
	} else {
		ret = create(id, mc, mem)
	}
	return
}

func (js *JsMem) Id() string {
	return js.id
}

// Snapshot returns a json-blob of the current game state
// ( the intended use is so the client cant save )
func (js *JsMem) Snapshot() (res string, err error) {
	buf := new(bytes.Buffer)
	if e := js.memory.Save(buf); e != nil {
		err = e
	} else {
		res = buf.String()
	}
	return
}

// Get mirrors http get: retrieving data from the game
func (js *JsMem) Get(path string) (resp string, err error) {
	if res, e := resource.FindResource(js.session, path); e != nil {
		err = e
	} else {
		doc := res.Query()
		if r, e := js.encode(doc); e != nil {
			err = e
		} else {
			resp = r
		}
	}
	return
}

// Post mirrors http post: sending commands for the turn
func (js *JsMem) Post(path string, body string) (resp string, err error) {
	if res, e := resource.FindResource(js.session, path); e != nil {
		err = e
	} else {
		reader := strings.NewReader(body)
		if doc, e := res.Post(reader); e != nil {
			err = e
		} else if r, e := js.encode(doc); e != nil {
			err = e
		} else {
			resp = r
		}
	}
	return
}

func create(id string, mc ingest.ModelCode, mem metal.ObjectValueMap) JsMem {
	meta := metal.NewMetal(mc.Model, mem)
	out := app.NewCommandOutput(id, meta, framework.NewStandardView(meta))
	sess, e := app.NewPartialSession(meta, mc.Code, out)
	if e != nil {
		panic(e)
	}
	return JsMem{id, mem, sess}
}

func (js JsMem) encode(doc resource.Document) (resp string, err error) {
	if prettyBytes, e := json.Marshal(doc); e != nil {
		err = e
	} else {
		resp = string(prettyBytes)
	}
	return
}
