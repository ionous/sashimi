// playing with code generating go from model data.
package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/ionous/sashimi/_examples/stories"
	"github.com/ionous/sashimi/compiler/call"
	"github.com/ionous/sashimi/compiler/extract"
	_ "github.com/ionous/sashimi/extensions"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/net/app"
	"github.com/ionous/sashimi/script"
	_ "github.com/ionous/sashimi/standard"
	"github.com/ionous/sashimi/util/ident"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// maybe a template would be nicer?
var header string = `
package fishy
import ( 
G "github.com/ionous/sashimi/game"
. "github.com/ionous/sashimi/standard"
. "github.com/ionous/sashimi/extensions"
"github.com/ionous/sashimi/util/lang"
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

// implments the compiler interface to turn a callback function into a callback marker
func (ex *Extracto) CompileCallback(cb G.Callback) (ret ident.Id, err error) {
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

var test = M.Model{
	Classes: M.ClassMap{"a": &M.ClassInfo{}},
}

//- signed integers (int, int8, int16, int32 and int64),
// - bool,
// - string,
// - float32 and float64,
// - []byte (up to 1 megabyte in length)
func generateClasses(m *M.Model) (ret string, err error) {
	const modelTemplate = `
{{range $id, $cls :=  .Classes}}
    type {{className $id}} struct {
    		{{if $cls.Parent}}{{className $cls.Parent.Id}}{{end}}
    	{{range $prop :=  $cls.Properties}}
    		{{propertyName $prop}} {{propertyType $prop}}
    	{{end}}
    }
{{ end }}
`
	className := func(id ident.Id) string {
		cls := m.Classes[id]
		return ident.MakeId(cls.Singular).String()
	}
	/// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		"className": className,
		"propertyName": func(p M.IProperty) (ret string) {
			// turn all relatives into pointers
			name := p.GetId().String()
			if p, isRel := p.(M.RelativeProperty); isRel && p.IsMany {
				ret = "//" + name
			} else {
				ret = name
			}
			return ret
		},
		"propertyType": func(p M.IProperty) (ret string) {
			switch p := p.(type) {
			case M.RelativeProperty:
				ret = "*" + className(p.Relates)
			case M.TextProperty:
				ret = "string"
			case M.NumProperty:
				ret = "float32"
			case M.EnumProperty:
				ret = "int"
			case M.PointerProperty:
				ret = "*" + className(p.Class)
			default:
				panic("unhandled type")
			}
			return ret
		},
	}
	modelT := template.Must(template.New("model").Funcs(funcMap).Parse(modelTemplate))

	fmt.Println(len(m.Classes))
	var buf = new(bytes.Buffer)
	if e := modelT.Execute(buf, m); e != nil {
		err = fmt.Errorf("error: executing!", e)
	} else {
		//b := buf.Bytes()
		lines := strings.Split(buf.String(), "\n")
		filtered := lines[:0]
		for _, line := range lines {
			s := strings.TrimSpace(line)
			if empty := len(s) == 0; !empty {
				filtered = append(filtered, s)
			}
		}
		s := strings.Join(filtered, "\n")
		b := []byte(s)
		if r, e := format.Source(b); e != nil {
			err = fmt.Errorf("error: formatting!", e)
			//fmt.Println(s)
			//fmt.Println("error: formatting!", e)
		} else {
			ret = string(r)
		}
	}
	return ret, err
}

func jsonClasses(m *M.Model) (err error) {
	res := app.ClassResource(m.Classes)
	//classes := res.Query()
	//text, _ := json.MarshalIndent(classes, "", " ")
	//fmt.Print(string(text))
	for k, _ := range m.Classes {
		if class, ok := res.Find(k.String()); !ok {
			err = fmt.Errorf("couldnt find class", k)
		} else {
			doc := class.Query()
			text, _ := json.MarshalIndent(doc, "", " ")
			fmt.Print(string(text))
		}
	}
	return err
}

func modelClasses(m *M.Model) (ret string, err error) {
	const modelTemplate = `
	package gen
	import ( 
		M "github.com/ionous/sashimi/model"
		//"github.com/ionous/sashimi/util/ident"
	)
{{range $id, $cls :=  .Classes}}
    var {{$id}} = M.ClassInfo {
    		{{if $cls.Parent}}Parent:&{{$cls.Parent.Id}},{{end}}
    		Id: "{{$id}}",
    		Plural: "{{$cls.Plural}}",
    		Singular: "{{$cls.Singular}}",
    }
{{ end }}`
	// {{range $prop :=  $cls.Properties}}
	//     		{{propertyName $prop}} {{propertyType $prop}}
	//     	{{end}}
	/// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{}
	modelT := template.Must(template.New("model").Funcs(funcMap).Parse(modelTemplate))

	fmt.Println(len(m.Classes))
	var buf = new(bytes.Buffer)
	if e := modelT.Execute(buf, m); e != nil {
		err = fmt.Errorf("error executing,", e)
	} else {
		//b := buf.Bytes()
		lines := strings.Split(buf.String(), "\n")
		filtered := lines[:0]
		for _, line := range lines {
			s := strings.TrimSpace(line)
			if empty := len(s) == 0; !empty {
				filtered = append(filtered, s)
			}
		}
		s := strings.Join(filtered, "\n")
		b := []byte(s)
		if r, e := format.Source(b); e != nil {
			err = fmt.Errorf("error formatting,", e)
			fmt.Println(s)
			//fmt.Println("error: formatting!", e)
		} else {
			ret = string(r)
		}
	}
	return ret, err
}

func write(s string) {
	if f, e := os.Create("$gen/gen.go"); e != nil {
		panic(e.Error())
	} else {
		defer f.Close()
		fmt.Fprint(f, s)
	}
}

func main() {
	script.AddScript(stories.A_Day_For_Fresh_Sushi)
	i := script.ScriptCount()
	fmt.Println(i, "scripts")
	s := script.InitScripts()
	d, _ := os.Getwd()
	ex := NewExtracto(d)
	if m, e := s.CompileCalls(os.Stdout, ex); e != nil {
		fmt.Println("error:", e)
	} else {
		// if e := jsonClasses(m); e != nil {
		// 	fmt.Println("error:", e)
		// }

		// if s, e := modelClasses(m); e != nil {
		// 	fmt.Println("error:", e)
		// } else {
		// 	write(s)
		// }
		var network bytes.Buffer        // Stand-in for a network connection
		enc := gob.NewEncoder(&network) // Will write to network.
		dec := gob.NewDecoder(&network) // Will read from network.

		// Encode (send) some values.
		if e := enc.Encode(m); e != nil {
			fmt.Println("error:", e)
		} else {

			var q M.Model
			if e := dec.Decode(&q); e != nil {
				fmt.Println("error:", e)
			}
		}

		// we could store just the model data
		// } else if f, e := os.Create("$gen/gen.go"); e != nil {
		// 	fmt.Println("error:", e)
		// } else {
		// 	defer f.Close()
		// 	ex.Fprint(f)
		// }
	}
	fmt.Println("done!")
}
