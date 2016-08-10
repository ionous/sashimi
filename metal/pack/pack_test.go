package pack

import (
	"encoding/json"
	"github.com/ionous/sashimi/compiler/model/modeltest"
	"github.com/ionous/sashimi/meta/metatest"
	"github.com/ionous/sashimi/metal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func match(t *testing.T, a, b metal.ObjectValueMap) {
	for k, v := range a {
		assert.Equal(t, v, b[k], k)
	}
}

func TestRawPack(t *testing.T) {
	src := metal.ObjectValueMap{"StatusBar.StatusBarInstancesLeft": "The Automat", "StatusBar.StatusBarInstancesRight": "\"Alice and the Galactic Traveller\" by everMany games", "AutomatHallDoor.OpenersOpen": "closed", "AutomatHallDoor.OpenersUnlocked": "locked", "Automat.RoomsVisited": "visited", "AliceDemo.StoriesPlaying": "playing"}

	//ObjectsAreEqualValues
	var pack ObjectValuePack
	if assert.NotPanics(t, func() {
		pack = Pack(src)
	}) {
		t.Log(pack)
		restored, err := Unpack(pack)
		if assert.NoError(t, err, "unpacking") {
			match(t, src, restored)
			match(t, restored, src)
		}
	}
}

func TestMetalPack(t *testing.T) {
	src := make(metal.ObjectValueMap)
	// src["a.b"] = []float32{0, 1}
	m := metal.NewMetal(modeltest.NewModel(), src)
	metatest.ApiTest(t, m, modeltest.TestInstance)

	var pack ObjectValuePack
	if assert.NotPanics(t, func() {
		pack = Pack(src)
	}) {
		restored, err := Unpack(pack)
		if assert.NoError(t, err, "unpacking") {
			match(t, src, restored)
			match(t, restored, src)
			text, _ := json.Marshal(src)
			packed, _ := json.Marshal(pack)
			textLen, packLen := len(text), len(packed)
			t.Log("textLen", textLen, "packLen", packLen)
			t.Log(string(text))
			t.Log(string(packed))
			assert.True(t, packLen < textLen, "yikes")

		}
	}
}
