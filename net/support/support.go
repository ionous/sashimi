package support

import (
	"log"
	"net/http"
	"path"
)

type FilePattern struct {
	public, dir string
}

// Dir generates a FilePattern for use with ServeMux
func Dir(pattern string) FilePattern {
	return FilePattern{pattern, pattern}
}

// RenameDir generates a FilePattern for use with ServeMux
func RenameDir(pattern, directory string) FilePattern {
	return FilePattern{pattern, directory}
}

// ServeMux extends http.ServeMux to provide additional functionality.
type ServeMux struct {
	*http.ServeMux
}

// NewServeMux extends http.NewServeMux to provide additional functionality.
func NewServeMux() ServeMux {
	return ServeMux{http.NewServeMux()}
}

// HandleText associates a static block of text with the passed url pattern.
func (mux ServeMux) HandleText(pattern, text string) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != pattern {
			http.NotFound(w, r)
		} else {
			w.Write([]byte(text))
		}
	})
}

// HandleFilePatterns creates a http.FileSystem for the passed root path, and associates it with the passed FilePatterns.
func (mux ServeMux) HandleFilePatterns(root string, fpats ...FilePattern) {
	r := http.Dir(root)
	for _, fpat := range fpats {
		var handler http.Handler
		if fpat.public != fpat.dir {
			handler = http.StripPrefix(fpat.public, http.FileServer(http.Dir(path.Join(root, fpat.dir))))
			log.Println("serving", fpat.public, "with", fpat.dir)
		} else {
			handler = http.FileServer(r)
			log.Println("serving", fpat.public)
		}
		mux.Handle(fpat.public, handler)
	}
}
