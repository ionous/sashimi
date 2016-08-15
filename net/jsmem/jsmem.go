package jsmem

import (
	"encoding/json"
	"errors"
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

type Create struct {
	mc          ingest.ModelCode
	recordInput bool
	save        SaveCallback
}

func NewCreator(mc ingest.ModelCode) *Create {
	return &Create{mc: mc}
}

func (c *Create) Save(save SaveCallback) *Create {
	c.save = save
	return c
}

func (c *Create) RecordInput() *Create {
	c.recordInput = true
	return c
}

// New creates a blank js session.
func (c *Create) NewGame() JsMem {
	mem := &memory{make(metal.ObjectValueMap), c.save, nil}
	if c.recordInput {
		mem.input = []string{}
	}
	return mem.createGame(c.mc)
}

// Restore uses the passed code and data as a starting point, and restores the saved json over top of it.
func (c *Create) RestoreGame(data Storage) (ret JsMem, err error) {
	if values, e := pack.Unpack(pack.ObjectValuePack(data)); e != nil {
		err = e
	} else {
		//metal.ObjectValueMap(values)
		mem := &memory{values, c.save, nil}
		ret = mem.createGame(c.mc)
	}
	return
}

type SaveCallback func(Storage, []string, bool) (slot string, err error)
type Storage pack.ObjectValuePack

type memory struct {
	values metal.ObjectValueMap
	save   SaveCallback
	input  []string
}

func (mem *memory) SaveGame(autosave bool) (ret string, err error) {
	if mem.save == nil {
		err = errors.New("no save method")
	} else {
		data := pack.Pack(mem.values)
		ret, err = mem.save(Storage(data), mem.input, autosave)
	}
	return
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
	hist := js.memory.input
	if hist != nil {
		js.memory.input = append(hist, body)
	}
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
