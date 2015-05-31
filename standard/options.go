package standard

import (
	"flag"
	C "github.com/ionous/sashimi/console"
)

//
// options and command line parsing for terminal style playback.
//
type Options struct {
	verbose, text, dump, hasConsole bool
	cons                            C.IConsole
}

// create options by reading the command line
func ParseCommandLine() Options {
	verbose := flag.Bool("verbose", false, "prints log output when true.")
	text := flag.Bool("text", false, "uses the simpler text console when true.")
	dump := flag.Bool("dump", false, "dump the model.")
	flag.Parse()
	return Options{verbose: *verbose, text: *text, dump: *dump}
}

func (this Options) SetVerbose(okay bool) Options {
	this.verbose = okay
	return this
}

func (this Options) UseTextConsole(okay bool) Options {
	this.text = okay
	return this
}

func (this Options) DumpModelAndExit(okay bool) Options {
	this.dump = okay
	return this
}

// override console with some external instance
func (this Options) SetConsole(cons C.IConsole) Options {
	this.cons = cons
	this.hasConsole = true
	return this
}
