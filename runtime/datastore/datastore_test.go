package datastore

import (
	"appengine/aetest"
	"bytes"
	"fmt"
	"github.com/ionous/sashimi/compiler/metal"
	"github.com/ionous/sashimi/compiler/metal/metaltest"
	"github.com/ionous/sashimi/compiler/model/modeltest"
	"github.com/ionous/sashimi/compiler/model/table"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/meta/metatest"
	"github.com/ionous/sashimi/runtime/datastore/dstest"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// https://golang.org/pkg/encoding/binary/#Write

func TestStoreNums(t *testing.T) {
	vals := []float32{-1.2, 3.1, 4.2}
	stream := new(bytes.Buffer)
	require.NoError(t, writeNums(stream, vals))
	b := stream.Bytes()
	t.Log(b)
	buf := bytes.NewBuffer(b)
	if res, e := readNums(buf); assert.NoError(t, e) {
		require.EqualValues(t, vals, res)
	}
}

func TestStoreStrings(t *testing.T) {
	vals := []string{"Abc", "one-two-three", "ending"}
	stream := new(bytes.Buffer)
	require.NoError(t, writeStrings(stream, vals))
	b := stream.Bytes()
	//t.Log(b)
	buf := bytes.NewBuffer(b)
	if res, e := readStrings(buf); assert.NoError(t, e) {
		require.EqualValues(t, vals, res)
	}
}

func TestStoreIds(t *testing.T) {
	vals := []ident.Id{ident.MakeId("Abc"), ident.MakeId("one-two-three"), ident.MakeId("ending")}
	stream := new(bytes.Buffer)
	require.NoError(t, writeIds(stream, vals))
	b := stream.Bytes()
	t.Log(b)
	buf := bytes.NewBuffer(b)
	if res, e := readIds(buf); assert.NoError(t, e) {
		require.EqualValues(t, vals, res)
	}
}

func TestEncodeDecode(t *testing.T) {
	type X struct {
		ptype meta.PropertyType
		val   interface{}
	}
	xs := []X{{
		meta.NumProperty, float32(2.3),
	}, {
		meta.TextProperty, "text",
	}, {
		meta.StateProperty, ident.MakeId("state"),
	}, {
		meta.ObjectProperty, ident.MakeId("object"),
	}, {
		meta.NumProperty | meta.ArrayProperty, []float32{1, 2, 3},
	}, {
		meta.TextProperty | meta.ArrayProperty, []string{"a", "b", "c"},
	}, {
		meta.ObjectProperty | meta.ArrayProperty, []ident.Id{
			ident.MakeId("first"),
			ident.MakeUniqueId(),
			ident.MakeId("last"),
		}},
	}

	for _, x := range xs {
		encoded, err := Encode(x.ptype, x.val)
		require.NoError(t, err, fmt.Sprintf("encoding %s %v", x.ptype, x.val))
		decoded, err := Decode(x.ptype, encoded)
		require.NoError(t, err, fmt.Sprintf("decoding %s %v %v", x.ptype, x.val, encoded))
		assert.EqualValues(t, x.val, decoded, "matches")
	}
}

// go test -run TestStoreData
func TestStoreData(t *testing.T) {
	kvs := &KeyValues{}
	mdl := metal.NewMetal(modeltest.NewModel(), kvs, make(table.Tables))

	if ctx, err := aetest.NewContext(nil); assert.NoError(t, err) {
		defer ctx.Close()

		// yuck! if we shadowed the meta, we could avoid this.
		kvs.mdl = mdl
		kvs.KeyGen = dstest.NewKeyGen(mdl, ctx, nil)
		kvs.ctx = ctx
		kvs.Reset()
		ctx.Infof("running api test..")
		metatest.ApiTest(t, mdl, modeltest.TestInstance)
		ctx.Infof("saving...")
		err := kvs.Save()
		require.NoError(t, err, "saving db")
		kvs.Reset()
		ctx.Infof("testing post conditions...")
		metaltest.PostConditions(t, kvs)
	}
}