package support

import (
	"log"
	"net/http"
	//	"path/filepath"
)

type FilePair struct {
	pattern, dir string
}

func Dir(pattern string) FilePair {
	return FilePair{pattern, ""}
}

type ServeMux struct {
	*http.ServeMux
}

func NewServeMux() ServeMux {
	return ServeMux{http.NewServeMux()}
}

func (this ServeMux) HandleText(pattern, text string) {
	this.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != pattern {
			http.NotFound(w, r)
		} else {
			w.Write([]byte(text))
		}
	})
}

func (this ServeMux) HandleFilePatterns(root string, pairs []FilePair) {
	fs := http.Dir(root)
	for _, pair := range pairs {
		this.HandleFilePattern(fs, pair)
	}
}

func (this ServeMux) HandleFilePattern(fs http.FileSystem, pair FilePair) {
	if dir := pair.dir; dir == "" {
		this.Handle(pair.pattern, http.FileServer(fs))
		log.Println("serving", pair.pattern)
	} else {
		//	http.Handle(this.pattern,
		//		http.StripPrefix(this.pattern, http.FileServer(http.Dir(path))))
		panic("needs testing")
	}
}
