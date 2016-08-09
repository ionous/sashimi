package mem

import (
	"encoding/json"
	"github.com/ionous/sashimi/metal"
)

// SaveHelper implements api.SaveLoad for a MemSaver interface
type SaveHelper struct {
	id     string
	values metal.ObjectValueMap
	saver  MemSaver
}

func NewSaveHelper(id string, values metal.ObjectValueMap, saver MemSaver) SaveHelper {
	return SaveHelper{id, values, saver}
}

//extract data and call
// implement runtime.api.SaveLoad, collect the model into json, push the string ( or bytes ) into memSaver
func (m SaveHelper) SaveGame(autosave bool) (ret string, err error) {
	if b, e := json.Marshal(m.values); e != nil {
		err = e
	} else {
		var saveId string
		if autosave {
			saveId = "autosave"
		}
		if r, e := m.saver.SaveBlob(saveId, b); e != nil {
			err = e
		} else {
			ret = r
		}
	}
	return
}

func LoadGame(slot string, saver MemSaver) (ret metal.ObjectValueMap, err error) {
	if data, e := saver.LoadBlob(slot); e != nil {
		err = e
	} else {
		values := make(metal.ObjectValueMap)
		if e := json.Unmarshal(data, &values); e != nil {
			err = e
		} else {
			ret = values
		}
	}
	return
}
