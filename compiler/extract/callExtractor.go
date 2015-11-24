package extract

import (
	"fmt"
	"github.com/ionous/sashimi/compiler/call"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/ident"
	"io"
	"io/ioutil"
)

type CallExtractor struct {
	pkgname string
	trace   io.Writer
	cfg     call.Config
	files   map[string]bool
	// fix, keep or make an array of structures for sorting?
	snippets map[string]Snippet
}

type Snippet struct {
	File    string
	Line    int
	Content string
}

func (s Snippet) String() string {
	return fmt.Sprintf("%s:%d", s.File, s.Line)
}

// path, optional, for more readable strings.
func NewCallExtractor(pkgname string, path string, trace io.Writer) *CallExtractor {
	return &CallExtractor{
		pkgname:  pkgname,
		trace:    trace,
		cfg:      call.Config{path},
		files:    make(map[string]bool),
		snippets: make(map[string]Snippet),
	}
}

func (cx *CallExtractor) Count() int {
	return len(cx.snippets)
}

func (cx *CallExtractor) Trace(args ...interface{}) {
	x := fmt.Sprintln(args...)
	io.WriteString(cx.trace, x)
}

// CompileCallback implments the compiler interface to turn a callback function into a callback marker.
func (cx *CallExtractor) CompileCallback(cb G.Callback) (ret ident.Id, err error) {
	m := cx.cfg.MakeMarker(cb)
	cx.Trace("compiling callback", m)
	if !cx.files[m.File] {
		cx.Trace("parsing file", m.File)
		cx.files[m.File] = true
		if bytes, e := ioutil.ReadFile(m.File); e != nil {
			err = e
			panic(e)
		} else {
			// alternatively, you could keep the ast, and dump it.
			err = ParseSource(m.File, bytes, func(file string, line int, bytes []byte) (err error) {
				x := cx.cfg.MarkFileLine(file, line)
				if id, e := x.Encode(); e != nil {
					err = e
				} else {
					//cx.Trace("snipping", name)
					str := ident.Dash(id)

					if old, exists := cx.snippets[str]; exists {
						err = fmt.Errorf("function already exists %s %s %s", x, old, bytes)
					} else {
						cx.snippets[str] = Snippet{
							File:    m.File,
							Line:    m.Line,
							Content: string(bytes),
						}
					}
				}
				return err
			})
		}
	}
	if err == nil {
		ret, err = m.Encode()
	}
	return
}
