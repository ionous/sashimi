package standard

import (
	"github.com/ionous/sashimi/parser"
	R "github.com/ionous/sashimi/runtime"
)

// STORE-FIX: but, currently there are other global variables as well.
// not perfectly sure how to fix this -- user data on the game ( and accessible via Play );
// an event watcher in standard game, and the ability to parameterize events with user data.
// currently the only way to support multiple games simultaneously would be multiple processes.
var TheParser *StandardParser

// CapturedInputCallback can override the standard game's processing of user input.
type CapturedInputCallback func(string) error

// Parser exists to inject a CapturedInputCallback into the path of the normal parsing.
type StandardParser struct {
	ObjectParser   *R.ObjectParser
	capturingInput CapturedInputCallback
}

func NewStandardParser(game *R.Game) (ret *StandardParser, err error) {
	// globals behave badly during test.
	// if TheParser != nil {
	// 	err = fmt.Errorf("multiple parsers registered")
	// } else
	if p, e := R.NewObjectParser(game); e != nil {
		err = e
	} else {
		ret = &StandardParser{p, nil}
		TheParser = ret
	}
	return ret, err
}

func (sp *StandardParser) ParseInput(input string) (ret parser.Matched, err error) {
	if capture := sp.capturingInput; capture == nil {
		ret, err = sp.ObjectParser.Parser.ParseInputString(input)
	} else {
		ret = parser.Matched{nil, func() (err error) {
			if e := capture(input); e != nil {
				err = e
			} else {
				sp.capturingInput = nil
			}
			return err
		}}
	}
	return ret, err
}

// CaptureInput hacks in menu support ( ex. conversations )
// FUTURE: push to inject, pop and hold on detect, the return code to excise.
func (sp *StandardParser) CaptureInput(cb CapturedInputCallback) CapturedInputCallback {
	old := sp.capturingInput
	sp.capturingInput = cb
	return old
}
