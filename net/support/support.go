package support

import (
	"log"
	"net/http"
)

type FilePattern string

// Dir generates a FilePattern for use with ServeMux
func Dir(pattern string) FilePattern {
	return FilePattern(pattern)
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
func (mux ServeMux) HandleFilePatterns(root string, filepats ...FilePattern) {
	fs := http.Dir(root)
	for _, filepat := range filepats {
		mux.HandleFilePattern(fs, filepat)
	}
}

// HandleFilePattern associates the passed file system with the passed FilePattern.
func (mux ServeMux) HandleFilePattern(fs http.FileSystem, filepat FilePattern) {
	mux.Handle(string(filepat), http.FileServer(fs))
	log.Println("serving", filepat)
}
