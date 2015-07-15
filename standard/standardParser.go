package standard

import (
	P "github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
)

// FIX: but, currently there are other global variables as well.
// not perfectly sure how to fix this -- user data on the game ( and accessible via Play );
// an event watcher in standard game, and the ability to parameterize events with user data.
// currently the only way to support multiple games simultaneously would be multiple processes.
var TheParser *Parser

// CapturedInputCallback can override the standard game's processing of user input.
type CapturedInputCallback func(string) error

// Parser exists to inject a CapturedInputCallback into the path of the normal parsing.
type Parser struct {
	*R.ObjectParser
	capturingInput CapturedInputCallback
}

func NewParser(game *R.Game) (ret *Parser, err error) {
	// globals behave badly during test.
	// if TheParser != nil {
	// 	err = fmt.Errorf("multiple parsers registered")
	// } else
	if p, e := R.NewParser(game); e != nil {
		err = e
	} else {
		ret = &Parser{p, nil}
		TheParser = ret
	}
	return ret, err
}

func (sp *Parser) ParseInput(input string) (ret P.Matched, err error) {
	if capture := sp.capturingInput; capture == nil {
		ret, err = sp.ObjectParser.ParseInput(input)
	} else {
		ret = P.Matched{nil, func() error { return capture(input) }}
	}
	return ret, err
}

// CaptureInput hacks in menu support ( ex. conversations )
// FUTURE: push to inject, pop and hold on detect, the return code to excise.
func (sp *Parser) CaptureInput(cb CapturedInputCallback) CapturedInputCallback {
	old := sp.capturingInput
	sp.capturingInput = cb
	return old
}
