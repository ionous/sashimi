package script

import (
	C "github.com/ionous/sashimi/compiler"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"io"
	"reflect"
)

type Script struct {
	blocks S.BuildingBlocks
	err    error
}

// Compile the current script into a model which can be used by the runtime.
func (s *Script) Compile(writer io.Writer) (res *M.Model, err error) {
	if s.err != nil {
		err = s.err
	} else {
		res, err = C.Compile(writer, s.blocks.GetStatements())
	}
	return res, err
}

// The main script function used to asset the existence of a class, instance, property, etc.
// The("example", s.Has("...") )
func (s *Script) The(key string, fragments ...IFragment) error {
	b := SubjectBlock{s, key, findSubject(key, fragments), &s.blocks}
	for _, frag := range fragments {
		if e := frag.MakeStatement(b); e != nil {
			s.err = errutil.Append(s.err, e)
		}
	}
	return s.err
}

// A=>The alias
func (s *Script) A(key string, fragments ...IFragment) {
	s.The(key, fragments...)
}

// Our=>The alias
func (s *Script) Our(key string, fragments ...IFragment) {
	s.The(key, fragments...)
}

// Execute( action, Matching(regexp).Or(other regexp) )
func (s *Script) Execute(what string, as Parsing) {
	s.blocks.NewAlias(what, as.phrases)
}

func (s *Script) Generate(what string, gen reflect.Type) {
	fields := S.GeneratorFields{what, gen}
	s.blocks.NewGlobal(fields, NewOrigin(1).Code())
}

// FIX? scan to find the subject called
func findSubject(key string, fragments []IFragment) string {
	subject := key
	for _, f := range fragments {
		if called, ok := f.(CalledFragment); ok {
			subject = called.subject
			break
		}
	}
	return subject
}
