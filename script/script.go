package script

import (
	C "github.com/ionous/sashimi/compiler"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"io"
)

type Script struct {
	blocks S.BuildingBlocks
	err    error
}

//
// Turn the current script into a model which can be used by the runtime.
//
func (this *Script) Compile(writer io.Writer) (res *M.Model, err error) {
	if this.err != nil {
		err = this.err
	} else {
		res, err = C.Compile(writer, this.blocks.GetStatements())
	}
	return res, err
}

//
// The main script function used to asset the existence of a class, instance, property, etc.
// The("example", this.Has("...") )
//
func (this *Script) The(key string, fragments ...IFragment) error {
	b := SubjectBlock{this, key, findSubject(key, fragments), &this.blocks}
	for _, frag := range fragments {
		if e := frag.MakeStatement(b); e != nil {
			this.err = errutil.Append(this.err, e)
		}
	}
	return this.err
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

//
// A=>The alias
//
func (this *Script) A(key string, fragments ...IFragment) {
	this.The(key, fragments...)
}

//
// Our=>The alias
//
func (this *Script) Our(key string, fragments ...IFragment) {
	this.The(key, fragments...)
}

//
// Execute( action, Matching(regexp).Or(other regexp) )
//
func (this *Script) Execute(what string, as Parsing) {
	this.blocks.NewAlias(what, as.phrases)
}
