package script

import (
	"github.com/ionous/sashimi/compiler"
	"github.com/ionous/sashimi/compiler/call"
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

func (s *Script) Compile(writer io.Writer) (res compiler.MemoryResult, err error) {
	if s.err != nil {
		err = s.err
	} else {
		res, err = compiler.Compile(writer, s.blocks.Statements())
	}
	return
}

// Compile the current script into a model which can be used by the runtime.
func (s *Script) CompileCalls(writer io.Writer, calls call.Compiler) (res *M.Model, err error) {
	if s.err != nil {
		err = s.err
	} else {
		cfg := compiler.Config{Calls: calls, Output: writer}
		res, err = cfg.Compile(s.blocks.Statements())
	}
	return res, err
}

// The main script function used to asset the existence of a class, instance, property, etc.
// Returns a placeholder variable, and an error -- both of which are intended to help distinguish it from The() used in callbacks.
// ex. The("example", s.Has("...") )
func (s *Script) The(key string, fragments ...IFragment) (int, error) {
	return s.the(key, fragments...)
}

// A=>The alias
func (s *Script) A(key string, fragments ...IFragment) (int, error) {
	return s.the(key, fragments...)
}

// Our=>The alias
func (s *Script) Our(key string, fragments ...IFragment) (int, error) {
	return s.the(key, fragments...)
}

func (s *Script) the(key string, fragments ...IFragment) (int, error) {
	b := SubjectBlock{s, key, findSubject(key, fragments), &s.blocks}
	for _, frag := range fragments {
		if e := frag.MakeStatement(b); e != nil {
			s.err = errutil.Append(s.err, e)
		}
	}
	return 0, s.err
}

// Execute( action, Matching(regexp).Or(other regexp) )
func (s *Script) Execute(what string, as Parsing) {
	origin := NewOrigin(1)
	alias := S.AliasFields{what, as.phrases}
	s.blocks.NewAlias(alias, origin.Code())
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
