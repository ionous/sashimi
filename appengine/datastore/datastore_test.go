package datastore

import (
	"appengine/aetest"
	"bytes"
	"fmt"
	"github.com/ionous/sashimi/compiler/model/modeltest"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/meta/metatest"
	"github.com/ionous/sashimi/metal/metaltest"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// https://golang.org/pkg/encoding/binary/#Write

func TestStoreNums(t *testing.T) {
	vals := []float64{-1.2, 3.1, 4.2}
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
	buf := bytes.NewBuffer(b)
	if res, e := readIds(buf); assert.NoError(t, e) {
		require.EqualValues(t, vals, res)
	}
}

func TestStoreCoding(t *testing.T) {
	type X struct {
		ptype meta.PropertyType
		val   interface{}
	}
	xs := []X{{
		meta.NumProperty, float64(2.3),
	}, {
		meta.TextProperty, "text",
	}, {
		meta.StateProperty, ident.MakeId("state"),
	}, {
		meta.ObjectProperty, ident.MakeId("object"),
	}, {
		meta.NumProperty | meta.ArrayProperty, []float64{1, 2, 3},
	}, {
		meta.TextProperty | meta.ArrayProperty, []string{"a", "b", "c"},
	}, {
		meta.ObjectProperty | meta.ArrayProperty, []ident.Id{
			ident.MakeId("first"),
			uuid.MakeUniqueId(),
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
	if ctx, err := aetest.NewContext(nil); assert.NoError(t, err) {
		defer ctx.Close()
		test := modeltest.NewModel()
		ds := NewModelStore(ctx, test, nil)
		//
		ctx.Infof("running api test..")
		metatest.ApiTest(t, ds.mdl, modeltest.TestInstance)
		ctx.Infof("saving...")
		err := ds.Flush()
		require.NoError(t, err, "saving db")
		ds.kvs.Reset()
		ctx.Infof("testing post conditions...")
		metaltest.VerifyPostConditions(t, ds.kvs)
	}
}
