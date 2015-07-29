package script

import "strings"

func Lines(a ...string) string {
	return strings.Join(a, "\n")
}
