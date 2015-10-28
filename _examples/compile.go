package main

import (
	"bytes"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/compiler/call"
	"github.com/ionous/sashimi/compiler/extract"
	_ "github.com/ionous/sashimi/extensions"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/script"
	_ "github.com/ionous/sashimi/standard"
	"go/format"
	"io"
	"io/ioutil"
	"os"
)

// maybe a template would be nicer?
var header string = `
package fishy
import ( 
G "github.com/ionous/sashimi/game"
. "github.com/ionous/sashimi/standard"
. "github.com/ionous/sashimi/extensions"
"bitbucket.org/pkg/inflect"
"fmt"
"strings"
)

// from script...
func Lines(a ...string) string {
	return strings.Join(a, "\n")
}

// future: munge (md5 maybe) the key, 
// wrap the value in a struct including file,line,func.
var callbacks = map[string]G.Callback {
`
var footer string = `}`

type Extracto struct {
	cfg       call.Config
	files     map[string]bool
	functions map[string][]byte
}

func NewExtracto(path string) *Extracto {
	return &Extracto{
		call.Config{path},
		make(map[string]bool),
		make(map[string][]byte),
	}
}

func (ex Extracto) Fprint(w io.Writer) (err error) {
	b := new(bytes.Buffer)
	b.WriteString(header)
	for k, v := range ex.functions {
		b.WriteString(fmt.Sprintf(`"%v": `, k))
		b.Write(v)
		b.WriteString(fmt.Sprintln(",\n"))
	}
	b.WriteString(footer)
	if p, e := format.Source(b.Bytes()); e != nil {
		err = e
	} else {
		w.Write(p)
	}
	return err
}
func (ex Extracto) Print() {
	ex.Fprint(os.Stdout)
}

func (ex *Extracto) Compile(cb G.Callback) (ret call.Marker, err error) {
	m := ex.cfg.MakeMarker(cb)
	if !ex.files[m.File] {
		ex.files[m.File] = true
		if bytes, e := ioutil.ReadFile(m.File); e != nil {
			err = e
			panic(e)
		} else {
			// alternatively, you could keep the ast, and dump it.
			err = extract.Extract(m.File, bytes, func(file string, line int, bytes []byte) (err error) {
				name := fmt.Sprintf("%s:%d", file, line)
				if old, exists := ex.functions[name]; exists {
					err = fmt.Errorf("function already exists %s %s %s", name, old, bytes)
				} else {
					ex.functions[name] = bytes
				}
				return err
			})
		}
	}
	if err == nil {
		ret = m
	}
	return
}

func main() {
	script.AddScript(stories.A_Day_For_Fresh_Sushi)
	i := script.ScriptCount()
	fmt.Println(i, "scripts")
	s := script.InitScripts()
	d, _ := os.Getwd()
	ex := NewExtracto(d)
	if _, e := s.CompileCalls(os.Stdout, ex); e != nil {
		fmt.Println("error:", e)
	} else if f, e := os.Create("$gen/gen.go"); e != nil {
		fmt.Println("error:", e)
	} else {
		defer f.Close()
		ex.Fprint(f)
	}
	fmt.Println("done!")
}
