package tests

import (
	"fmt"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

// tests some api things
// requires at least one instance and one class
// the class valueerties should include "num", "text, "state", "object" of corresponding types
// "state" should contain choices "yes" and "no"
func ApiTest(t *testing.T, mdl api.Model, instId ident.Id) {
	require.NotNil(t, mdl)
	_, noExist := mdl.GetInstance(ident.MakeUniqueId())
	require.False(t, noExist)
	// find inst by index
	var inst api.Instance
	require.True(t, mdl.NumInstance() > 0, "need instances to test")
	for i := 0; i < mdl.NumInstance(); i++ {
		indirect := mdl.InstanceNum(i)
		if instId == indirect.GetId() {
			assert.Nil(t, inst, "instance id should only exist once")
			inst = indirect
		}
	}
	require.NotNil(t, inst, "should have found instance by index")
	// find inst by id
	if direct, ok := mdl.GetInstance(instId); assert.True(t, ok, "get instance by id") {
		require.True(t, direct == inst, "equality is necessary for the sake of game object adapter")
		require.EqualValues(t, instId, direct.GetId(), "self id test")
	}
	//
	require.True(t, inst.NumProperty() > 0, "need properties to test")
	value := inst.PropertyNum(0)
	if cls := inst.GetParentClass(); assert.NotNil(t, cls) {
		if plookup, ok := cls.GetProperty(value.GetId()); assert.True(t, ok, "find instance property in class") {
			require.True(t, plookup.GetId() == value.GetId())
		}
	}

	type TestValue struct {
		MethodMaker
		original, value interface{}
	}

	methods := []TestValue{
		{MethodMaker("Num"), float32(0), float32(32)},
		{MethodMaker("Text"), "", "text"},
		{MethodMaker("State"), ident.Id("No"), ident.Id("Yes")},
		{MethodMaker("Object"), ident.Empty(), inst.GetId()},
	}
	testMethods := func(src api.Prototype, eval func(TestValue, reflect.Value)) {
		// test class values;
		// test instance values
		for _, test := range methods {
			// for every property type: num, text, state
			pid := ident.MakeId(test.String())
			// get the property
			if p, ok := src.GetProperty(pid); assert.True(t, ok, test.String()) {
				// request the value from the property
				if value := reflect.ValueOf(p.GetValue()); assert.True(t, value.IsValid()) {
					// test getting the vaule of the appropriate type succeeds
					require.NotPanics(t, func() {
						prev := test.GetFrom(value)
						require.EqualValues(t, prev, test.original, fmt.Sprintf("original value %s, %v", test, prev))
					}, fmt.Sprintf("get default value: %s", pid))

					// custom testing for instances and classes
					eval(test, value)

					// test every other property type fails
					for _, m := range methods {
						if m != test {
							// panic all other methods:
							var v interface{}
							require.Panics(t, func() {
								v = m.GetFrom(value)
							}, fmt.Sprintf("trying to get %s from %s; returned %v", m, test, v))
							require.Panics(t, func() {
								m.SetTo(value, m.value)
							}, fmt.Sprintf("trying to set %s to %s", m, test))
						}
					}
				}
			}
		}
	}
	testMethods(inst, func(test TestValue, value reflect.Value) {
		// instances can get and set values
		require.NotPanics(t, func() {
			test.SetTo(value, test.value)
		}, fmt.Sprintf("instance set value: %s", test))
		require.NotPanics(t, func() {
			next := test.GetFrom(value)
			require.EqualValues(t, next, test.value, fmt.Sprintf("%v:%T should be %v:%T", next, next, test.value, test.value))
		}, fmt.Sprintf("get changed value: %s", test))
	})

	testMethods(inst.GetParentClass(), func(test TestValue, value reflect.Value) {
		// classes disallow set values
		require.Panics(t, func() {
			test.SetTo(value, test.value)
		}, fmt.Sprintf("class set value: %s", test))
	})
}

type MethodMaker string

func (m MethodMaker) String() string {
	return string(m)
}
func (m MethodMaker) GetFrom(value reflect.Value) (ret interface{}) {
	name := "Get" + m.String()
	if n := value.MethodByName(name); !n.IsValid() {
		panic(fmt.Sprintf("method didnt exist %s", name))
	} else {
		ret = n.Call([]reflect.Value{})[0].Interface()
	}
	return
}

func (m MethodMaker) SetTo(value reflect.Value, v interface{}) {
	name := "Set" + m.String()
	if n := value.MethodByName(name); !n.IsValid() {
		panic(fmt.Sprintf("method didnt exist %s", name))
	} else {
		n.Call([]reflect.Value{reflect.ValueOf(v)})
	}
}
