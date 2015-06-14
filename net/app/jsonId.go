package app

import (
	M "github.com/ionous/sashimi/model"
	"strings"
)

func jsonId(id M.StringId) string {
	return strings.Join(id.Split(), "-")
}
