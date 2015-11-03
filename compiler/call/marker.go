package call

import (
	"crypto/md5"
	"encoding/json"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
	"reflect"
	"runtime"
	"strings"
)

type Marker struct {
	M.Callback
}

func MakeMarker(cb G.Callback) Marker {
	v := reflect.ValueOf(cb)
	pc := v.Pointer()
	f := runtime.FuncForPC(pc)
	file, line := f.FileLine(pc - 1)
	return Marker{M.Callback{file, line, 0}}
}

func (cfg Config) MakeMarker(cb G.Callback) Marker {
	m := MakeMarker(cb)
	if strings.HasPrefix(m.File, cfg.BasePath+"/") {
		m.File = m.File[len(cfg.BasePath)+1:]
	}
	return m
}

func (m Marker) Encode() (ret ident.Id, err error) {
	if b, e := json.Marshal(m); e != nil {
		err = e
	} else {
		ret = ident.Id(string(md5.New().Sum(b)))
	}
	return
}
