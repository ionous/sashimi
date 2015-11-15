package internal

import (
	G "github.com/ionous/sashimi/game"
)

type nullList PropertyPath

func (_ nullList) Len() int                  { return 0 }
func (n nullList) Get(int) G.IValue          { return nullValue(n) }
func (n nullList) Contains(interface{}) bool { return false }
func (n nullList) AppendNum(float32)         {}
func (n nullList) AppendText(string)         {}
func (n nullList) AppendObject(G.IObject)    {}
func (n nullList) Reset()                    {}
