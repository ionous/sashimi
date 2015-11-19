package call

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	M "github.com/ionous/sashimi/compiler/model"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/ident"
	"io"
	"reflect"
	"runtime"
	"strings"
)

type Marker struct {
	M.CallbackModel
}

func MakeMarker(cb G.Callback) Marker {
	v := reflect.ValueOf(cb)
	pc := v.Pointer()
	f := runtime.FuncForPC(pc)
	file, line := f.FileLine(pc - 1)
	return Marker{M.CallbackModel{file, line, 0}}
}

func (cfg Config) MakeMarker(cb G.Callback) Marker {
	m := MakeMarker(cb)
	if strings.HasPrefix(m.File, cfg.BasePath+"/") {
		m.File = m.File[len(cfg.BasePath)+1:]
	}
	return m
}

func (m Marker) Encode() (ret ident.Id, err error) {
	if text, e := json.Marshal(m); e != nil {
		err = e
	} else {
		hash := md5.New()
		io.WriteString(hash, string(text))
		b := hash.Sum(nil)
		s := fmt.Sprintf("%x", b)
		if len(s) > len("e2c569be17396eca2a2e3c11578123ed") {
			panic(s)
		}
		ret = ident.MakeId(s)
	}
	return
}
