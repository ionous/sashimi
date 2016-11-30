package change_test

import (
	"github.com/ionous/sashimi/change"
	"github.com/ionous/sashimi/compiler/model/modeltest"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/meta/metatest"
	"github.com/ionous/sashimi/metal"
	"github.com/ionous/sashimi/metal/metaltest"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/require"
	"testing"
)

type ChangeMock struct {
	t                   *testing.T
	values              metal.ObjectValue
	nums, texts, states int
}

func (p *ChangeMock) NumChange(obj meta.Instance, prop ident.Id, prev, next float64) {
	require.NotEqual(p.t, prev, next)
	p.values.SetValue(obj.GetId(), prop, next)
	p.nums++
}

func (p *ChangeMock) TextChange(obj meta.Instance, prop ident.Id, prev, next string) {
	require.NotEqual(p.t, prev, next)
	p.values.SetValue(obj.GetId(), prop, next)
	p.texts++
}

func (p *ChangeMock) StateChange(obj meta.Instance, prop ident.Id, prev, next ident.Id) {
	require.NotEqual(p.t, prev, next)
	p.values.SetValue(obj.GetId(), prop, next)
	p.states++
}

func (p *ChangeMock) ReferenceChange(obj meta.Instance, prop, other ident.Id, prev, next meta.Instance) {
	require.NotEqual(p.t, prev, next)
	var id ident.Id
	if next != nil {
		id = next.GetId()
	}
	p.values.SetValue(obj.GetId(), prop, id)
}

func TestChange(t *testing.T) {
	ch := &ChangeMock{t: t, values: make(metal.ObjectValueMap)}
	watched := change.NewModelWatcher(ch, metal.NewMetal(modeltest.NewModel(), make(metal.ObjectValueMap)))
	metatest.ApiTest(t, watched, modeltest.TestInstance)
	require.True(t, ch.nums > 0, "nums")
	require.True(t, ch.texts > 0, "texts")
	require.True(t, ch.states > 0, "states")
	// our mock writes to its own object value map so that we can use metaltest to verify the expected changes
	metaltest.VerifyPostValues(t, ch.values)
	// lists are not currently watched
	// metaltest.VerifyPostLists(t, ch.values)
}
