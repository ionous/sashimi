package extract

import (
	"bytes"
	"fmt"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/util/ident"
	"go/format"
	"io"
	"strconv"
	"text/template"
)

type ModelConfig struct {
	snippet ModelSnippet
}

type ModelSnippet struct {
	Model   *M.Model
	PkgName string
}

func NewModelConfig(m *M.Model) *ModelConfig {
	return &ModelConfig{ModelSnippet{Model: m}}
}

func (cfg *ModelConfig) PackageName(s string) *ModelConfig {
	cfg.snippet.PkgName = s
	return cfg
}

func (cfg *ModelConfig) GetSnippet() ModelSnippet {
	return cfg.snippet
}

// cool, but unless we can deep copy the output, this is useless
func WriteModel(w io.Writer, s ModelSnippet) (err error) {
	b := new(bytes.Buffer)
	if e := dataTemplate.Execute(b, s); e != nil {
		err = e
	} else {
		snippets := b.Bytes()
		if p, e := format.Source(snippets); e != nil {
			w.Write(snippets)
			panic(e)
		} else {
			_, err = w.Write(p)
		}
	}
	return err
}

// Enums are stored as int;
// Numbers as float64;
// Pointers as ident.Id;
// Text as string.
var funcMap = template.FuncMap{
	"eval": func(v interface{}) string {
		switch v := v.(type) {
		case int:
			return strconv.Itoa(v)
		case float64:
			return strconv.FormatFloat(v, 'g', -1, 64)
		case ident.Id:
			return `"` + string(v) + `"`
		case string:
			return fmt.Sprintf("%q", v)
		default:
			panic(fmt.Sprintf("unknown type %T", v))
		}
	},
}

//
var dataTemplate = template.Must(template.New("dataTemplate").Funcs(funcMap).Parse(` 
{{ template "DataModel" . }}

{{ define "ValueList" }} {{ range $i, $el := . }} {{ eval $el }}, {{ end }} {{ end }}

{{ define "Actions" }}
	Actions: Actions{ {{range $id, $n := .Actions }}
		{{ eval $id  }}: &ActionModel{
			Id : {{ eval $n.Id  }},
			Name: {{ eval $n.Name }},
			EventId : {{ eval $n.EventId }},
			NounTypes:   []ident.Id{ {{ template "ValueList" $n.NounTypes }} },
	   },{{ end }} 
	},{{ end }}

{{ define "props" }} []PropertyModel{ {{ range $i, $el := . }} PropertyModel{
 		Id : {{ eval $el.Id }},
		Type: {{ $el.Type }},
		Name: {{ eval $el.Name }},
		Relates : {{ eval $el.Relates }},
		Relation : {{ eval $el.Relation }},
		IsMany : {{ print $el.IsMany }},
 	},{{ end }}
},{{ end }}

{{ define "Classes" }}
	Classes: Classes{ {{ range $id, $n := .Classes }}
			{{ eval $id  }}: &ClassModel{
			Id : {{ eval $n.Id }},
			Parents:  []ident.Id{ {{ template "ValueList"  $n.Parents }} },
			Plural : {{ eval $n.Plural }},
			Singular : {{ eval $n.Singular }},
			Properties: {{ template "props"  $n.Properties }}
		},{{ end }}
},{{ end }}

{{ define "Enumerations" }}
	Enumerations: Enumerations{ {{ range $id, $n := .Enumerations }}
			{{ eval $id }}: &EnumModel {
				Choices:  []ident.Id{ {{ template "ValueList"  $n.Choices }} },
		},{{ end }}
},{{ end }}

{{ define "Callbacks" }} []ListenerModel{  
{{ range $i, $el := . }} ListenerModel{
		Instance:{{ eval $el.Instance}},
		Class: {{ eval $el.Class }},
		Callback: {{ eval $el.Callback }},
		Options:ListenerOptions( {{ $el.Options }} ),
	},{{ end }} 
},{{ end }}

{{ define "Events" }}
Events: Events{  {{ range $id, $n := .Events }}
		{{ eval $id }}: &EventModel{
		Id : {{ eval $n.Id }},
		Name: {{ eval $n.Name }},
		Capture: {{ template "Callbacks"  $n.Capture }}
		Bubble: {{ template "Callbacks"  $n.Bubble }}
	},{{ end }}
},{{ end }}

{{ define "Values" }} Values{ {{ range $id, $n := . }}
		{{ eval $id }}: {{eval $n}},{{end}}
},{{ end }}

{{ define "Instances" }}
Instances: Instances{ 
		{{ range $id, $n := .Instances }}
			{{ eval $id }}: &InstanceModel{
			Id : {{ eval $n.Id }},
			Class: {{ eval $n.Class }},
			Name: {{ eval $n.Name }},
			Values: {{ template "Values" $n.Values }}
		},{{ end }}
},{{ end }}

{{ define "Aliases" }}
Aliases: Aliases{ {{ range $key, $n := .Aliases }} 
		{{ eval $key }}:  []ident.Id{ {{ template "ValueList" $n }} },{{ end }}
},{{ end }}

{{ define "ParserActions" }}
ParserActions: []ParserAction{
	{{ range $idx, $n := .ParserActions }}  ParserAction {
		Action:  {{ eval $n.Action }},
		Commands: []string{ {{ template "ValueList" $n.Commands }}  },
	},{{ end }}
},{{ end }}

{{ define "Relations" }}
Relations: Relations{ 
	{{ range $id, $n := .Relations }}
		{{ eval $id }}: &RelationModel{
		Id : {{ eval $n.Id }},
		Name:  {{ eval $n.Name }},
		Style: {{ $n.Style }},
		Source: {{ eval $n.Source }},
		Target: {{ eval $n.Target }},		
	},{{ end }}
},{{ end }}

{{ define "SingleToPlural" }}
SingleToPlural: SingleToPlural{
	{{ range $id, $n := .SingleToPlural }}
		{{ eval $id }} : {{ eval $n }},{{ end }}
},{{ end }}

{{ define "DataModel" }}
	package {{.PkgName}}
	import ( 
	  "github.com/ionous/sashimi/util/ident"
	. "github.com/ionous/sashimi/compiler/model"
	)
	
	var data = &Model { 
		{{ template "Actions" .Model }}
		{{ template "Classes" .Model }}
		{{ template "Enumerations" .Model }}
		{{ template "Events" .Model }}
		{{ template "Instances"  .Model }}
		{{ template "Aliases"        .Model }}
		{{ template "ParserActions"  .Model }}
		{{ template "Relations"      .Model }}
		{{ template "SingleToPlural" .Model }}
	}
{{ end }}		
`))
