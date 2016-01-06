package script

import (
	"github.com/ionous/sashimi/util/lang"
	"strings"
)

func Lines(a ...string) string {
	return strings.Join(a, lang.NewLine)
}
