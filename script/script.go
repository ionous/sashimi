package script

import (
	"fmt"
	C "github.com/ionous/sashimi/compiler"
	M "github.com/ionous/sashimi/model"
	S "github.com/ionous/sashimi/source"
	"log"
)

type Script struct {
	blocks S.BuildingBlocks
	err    error
}

//
// Turn the current script into a model which can be used by the runtime.
//
func (this *Script) Compile(log *log.Logger) (res *M.Model, err error) {
	if this.err != nil {
		err = this.err
	} else {
		res, err = C.Compile(log, this.blocks.GetBlocks())
	}
	return res, err
}

//
// The main script function used to asset the existence of a class, instance, property, etc.
// The("example", this.Has("...") )
//
func (this *Script) The(key string, fragments ...IFragment) {
	// ugh. suggestions are welcome:
	// scan to find the subject called
	subject, calledIndex := key, -1
	for i, f := range fragments {
		if called, okay := f.(CalledFragment); okay {
			if calledIndex >= 0 {
				e := fmt.Errorf("multiplle call statements?")
				this.err = C.AppendError(this.err, e)
				break
			}
			subject, calledIndex = called.subject, i
		}
	}

	for _, frag := range fragments {
		b := SubjectBlock{subject, key, &this.blocks}
		if e := frag.MakeStatement(b); e != nil {
			this.err = C.AppendError(this.err, e)
		}
	}
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
