package framework

import (
	"flag"
	C "github.com/ionous/sashimi/console"
)

//
// options and command line parsing for terminal style playback.
//
type Options struct {
	verbose, text, dump, load, hasConsole bool
	cons                                  C.IConsole
}

// create options by reading the command line
func ParseCommandLine() Options {
	verbose := flag.Bool("verbose", false, "prints log output when true.")
	text := flag.Bool("text", false, "uses the simpler text console when true.")
	dump := flag.Bool("dump", false, "dump the model.")
	load := flag.Bool("load", false, "load the story save game.")
	flag.Parse()
	return Options{verbose: *verbose, text: *text, dump: *dump, load: *load}
}

func (opt Options) SetVerbose(okay bool) Options {
	opt.verbose = okay
	return opt
}

func (opt Options) UseTextConsole(okay bool) Options {
	opt.text = okay
	return opt
}

func (opt Options) DumpModelAndExit(okay bool) Options {
	opt.dump = okay
	return opt
}

// override console with some external instance
func (opt Options) SetConsole(cons C.IConsole) Options {
	opt.cons = cons
	opt.hasConsole = true
	return opt
}
