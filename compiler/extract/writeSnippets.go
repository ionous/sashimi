package extract

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"text/template"
)

func WriteSnippets(w io.Writer, cx *CallExtractor, packages ...string) (err error) {
	if t, e := template.New("pkg").Parse(templateText); e != nil {
		err = fmt.Errorf("error parsing template: %s", e)
	} else {
		b := new(bytes.Buffer)
		s := templateData{cx.pkgname, packages, cx.snippets}
		if e := t.Execute(b, s); e != nil {
			err = fmt.Errorf("error executing template: %s", e)
		} else {
			snippets := b.Bytes()
			if p, e := format.Source(snippets); e != nil {
				//err = fmt.Errorf("error formatting source: %s", e)
				w.Write(snippets)
				panic(e)
			} else {
				_, err = w.Write(p)
			}
		}
	}
	return err
}

//
var templateText string = `
package {{.PkgName}}
import ( 
G "github.com/ionous/sashimi/game"
. "github.com/ionous/sashimi/standard"
. "github.com/ionous/sashimi/extensions"
"github.com/ionous/sashimi/util/lang"
"github.com/ionous/sashimi/util/ident"
"fmt"
"strings"{{ range .Packages }}
"{{.}}"{{ end }}
)

// from script...
func Lines(a ...string) string {
	return strings.Join(a, "\n")
}

var code = map[ident.Id]G.Callback {
	{{ range $id, $snippet := .Snippets }}
		// {{$snippet.File}}:{{$snippet.Line}}
		"{{ $id }}" : {{$snippet.Content}},
	{{end }} 
}
`

type templateData struct {
	PkgName  string
	Packages []string
	Snippets map[string]Snippet
}
