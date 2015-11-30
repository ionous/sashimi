package extract

import (
	"bytes"
	"encoding/json"
	"fmt"
	M "github.com/ionous/sashimi/compiler/model"
	"go/format"
	"io"
)

func WriteJsonModel(w io.Writer, name string, m *M.Model) (err error) {
	b := new(bytes.Buffer)
	if json, e := json.MarshalIndent(m, "", " "); e != nil {
		err = e
	} else {
		fmt.Fprintln(b, "package ", name)
		fmt.Fprintf(b, "var data=`%s`", json)
		if p, e := format.Source(b.Bytes()); e != nil {
			err = e
		} else {
			_, err = w.Write(p)
		}
	}
	return
}
