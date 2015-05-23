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

func (this ServeMux) HandleFilePatterns(root string, pairs []FilePair) {
	fs := http.Dir(root)
	for _, pair := range pairs {
		this.HandleFilePattern(fs, pair)
	}
}

func (this ServeMux) HandleFilePattern(fs http.FileSystem, pair FilePair) {
	if dir := pair.dir; dir == "" {
		this.Handle(pair.pattern, http.FileServer(fs))
		log.Println("serving", pair)
	} else {
		//	http.Handle(this.pattern,
		//		http.StripPrefix(this.pattern, http.FileServer(http.Dir(path))))
		panic("needs testing")
	}
}

/*
		s := &http.Server{
			Addr:           ":8080",
			Handler:        myHandler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	// nil means DefaultServeMux


*/
