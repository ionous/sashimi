package call

import (
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
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
