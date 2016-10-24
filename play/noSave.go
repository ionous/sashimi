package play

import (
	"github.com/ionous/sashimi/util/errutil"
)

type noSaveLoad struct{}

func (noSaveLoad) SaveGame(autosave bool) (string, error) {
	return "", errutil.New("not implemented")
}
