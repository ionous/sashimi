package jsmem

import (
	"encoding/json"
	"github.com/ionous/sashimi/compiler/ingest"
	"github.com/ionous/sashimi/metal"
	"github.com/ionous/sashimi/metal/pack"
	"github.com/ionous/sashimi/net/app"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/standard/framework"
	"strings"
)

type JsMem struct {
	session *app.PartialSession
	memory  *memory
}

type SaveCallback func(Storage, bool) (slot string, err error)
type Storage pack.ObjectValuePack

// New creates a blank js session.
// main can contain the mc pairing.
func New(mc ingest.ModelCode, save SaveCallback) JsMem {
	mem := &memory{make(metal.ObjectValueMap), save}
	return mem.createGame(mc)
}

// Restore uses the passed code and data as a starting point, and restores the saved json over top of it.
func Restore(mc ingest.ModelCode, data Storage, save SaveCallback) (ret JsMem, err error) {
	if values, e := pack.Unpack(pack.ObjectValuePack(data)); e != nil {
		err = e
	} else {
		//metal.ObjectValueMap(values)
		mem := &memory{values, save}
		ret = mem.createGame(mc)
	}
	return
}

type memory struct {
	values metal.ObjectValueMap
	save   SaveCallback
}

func (mem *memory) SaveGame(autosave bool) (string, error) {
	data := pack.Pack(mem.values)
	return mem.save(Storage(data), autosave)
}

func (mem *memory) createGame(mc ingest.ModelCode) JsMem {
	meta := metal.NewMetal(mc.Model, mem.values)
	out := app.NewCommandOutput("gopherjs", meta, framework.NewStandardView(meta))
	sess, e := app.NewPartialSession(meta, mc.Code, mem, out)
	if e != nil {
		panic(e)
	}
	return JsMem{sess, mem}
}

// Get mirrors http get: retrieving data from the game
func (js *JsMem) Get(_, path string) (ret string, err error) {
	if res, e := resource.FindResource(js.session, path); e != nil {
		err = e
	} else if r, e := js.encode(res.Query()); e != nil {
		err = e
	} else {
		ret = r
	}
	return
}

// Post mirrors http post: sending commands for the turn
func (js *JsMem) Post(_, body string) (ret string, err error) {
	reader := strings.NewReader(body)
	if doc, e := js.session.Post(reader); e != nil {
		err = e
	} else if r, e := js.encode(doc); e != nil {
		err = e
	} else {
		ret = r
	}
	return
}

func (js JsMem) encode(doc resource.Document) (resp string, err error) {
	if prettyBytes, e := json.Marshal(doc); e != nil {
		err = e
	} else {
		resp = string(prettyBytes)
	}
	return
}
