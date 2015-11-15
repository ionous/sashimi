package internal

import (
	G "github.com/ionous/sashimi/game"
)

type nullValue PropertyPath

func (_ nullValue) Set(value G.IValue)  {}
func (_ nullValue) Num() (ret float32)  { return }
func (_ nullValue) SetNum(float32)      {}
func (n nullValue) Object() G.IObject   { return NullObject("") }
func (_ nullValue) SetObject(G.IObject) {}
func (_ nullValue) Text() (ret string)  { return }
func (_ nullValue) SetText(string)      {}
