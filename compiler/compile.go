package compiler

import (
	"github.com/ionous/sashimi/compiler/call"
	"io"
)

type Config struct {
	Calls  call.Compiler
	Output io.Writer
}
