package metatest

import (
	"fmt"
	"github.com/ionous/sashimi/meta"
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
func ApiTest(t *testing.T, mdl meta.Model, instId ident.Id) {
	require.NotNil(t, mdl)
	_, noExist := mdl.GetInstance(ident.MakeUniqueId())
	require.False(t, noExist)
	// find inst by index
	var inst meta.Instance
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
	// value := inst.PropertyNum(0)
	// if cls := inst.GetParentClass(); assert.NotNil(t, cls) {
	// 	if plookup, ok := cls.GetProperty(value.GetId()); assert.True(t, ok, "find instance property in class") {
	// 		require.True(t, plookup.GetId() == value.GetId())
	// 	}
	// }

	type TestValue struct {
		MethodMaker
		original, value interface{}
	}

	methods := []TestValue{
		{MethodMaker("Num"), float32(0), float32(32)},
		{MethodMaker("Text"), "", "text"},
		{MethodMaker("State"), ident.MakeId("no"), ident.MakeId("yes")},
		{MethodMaker("Object"), ident.Empty(), inst.GetId()},
	}

	testValue := func(v meta.Value, test TestValue, eval func(TestValue, reflect.Value)) {
		if value := reflect.ValueOf(v); assert.True(t, value.IsValid()) {
			// test getting the vaule of the appropriate type succeeds
			var prev interface{}
			require.NotPanics(t, func() {
				prev = test.GetFrom(value)
			}, fmt.Sprintf("get default value: %s", test.String()))
			require.EqualValues(t, test.original, prev,
				fmt.Sprintf("original value from %s", test))

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
					var err error
					require.Panics(t, func() {
						err = m.SetTo(value, m.value)
					}, fmt.Sprintf("trying to set %s to %s", m, test))
					assert.NoError(t, err)
				}
			}
		}
	}

	testMethods := func(src meta.Prototype, eval func(TestValue, reflect.Value)) {
		// test class values;
		// test instance values
		for _, test := range methods {
			// get the property
			if p, ok := src.FindProperty(test.String()); assert.True(t, ok, test.String()) {
				// request the value from the property
				var v meta.Value
				require.NotPanics(t, func() { v = p.GetValue() })
				require.Panics(t, func() { p.GetValues() })
				testValue(v, test, eval)
			}
		}
	}
	_ = testMethods
	// instances can get and set values
	// value is reflect valueOf meta.Value
	testInstanceValue := func(test TestValue, value reflect.Value) {
		var err error
		require.NotPanics(t, func() {
			err = test.SetTo(value, test.value)
		}, fmt.Sprintf("instance set value panic: %s", test))
		require.NoError(t, err, fmt.Sprintf("instance set value error: %s", test))

		var next interface{}
		require.NotPanics(t, func() {
			next = test.GetFrom(value)
		}, fmt.Sprintf("get changed value: %s", test))
		require.EqualValues(t, next, test.value, fmt.Sprintf("%v:%T should be %v:%T", next, next, test.value, test.value))
	}
	testMethods(inst, testInstanceValue)

	// classes disallow set values
	// value is reflect valueOf meta.Value
	testClassValue := func(test TestValue, value reflect.Value) {
		var err error
		require.Panics(t, func() {
			err = test.SetTo(value, test.value)
		}, fmt.Sprintf("class set value: %s", test))
		require.NoError(t, err, fmt.Sprintf("class set value error: %s", test))
	}
	_ = testClassValue
	testMethods(inst.GetParentClass(), testClassValue)

	testArrays := func(src meta.Prototype, eval func(TestValue, meta.Values)) {
		// test class values;
		// test instance values
		for _, test := range methods {
			// no state arrays right now
			if test.String() == "State" {
				continue
			}
			// for every property type: num, text, state
			name := test.String() + "s"
			// get the property
			if p, ok := src.FindProperty(name); assert.True(t, ok, fmt.Sprintf("trying to get property %s.%s", src, name)) {
				// request the value from the property
				var vs meta.Values
				require.Panics(t, func() { p.GetValue() })
				require.NotPanics(t, func() { vs = p.GetValues() })
				var cnt int
				require.NotPanics(t, func() {
					cnt = vs.NumValue()
				}, fmt.Sprintf("trying to num value %s.%s", src, name))
				require.Equal(t, 0, cnt)

				eval(test, vs)

				// test alld other appends fails
				value := reflect.ValueOf(vs)
				for _, m := range methods {
					if m != test {
						// panic all other methods:
						require.Panics(t, func() {
							m.Append(value, test.value)
						}, fmt.Sprintf("trying to append %s from %s", m, test))
					}
				}
			}
		}
	}

	testArrays(inst, func(test TestValue, vs meta.Values) {
		value := reflect.ValueOf(vs)

		// append
		for i := 0; i < 3; i++ {
			// instances can get and set values
			require.NotPanics(t, func() {
				test.Append(value, test.original)
			}, fmt.Sprintf("instance set value: %s", test))
			cnt := vs.NumValue()
			require.EqualValues(t, i+1, cnt, fmt.Sprintf("cnt(%d) should be bigger", cnt))
			next := vs.ValueNum(i)
			testValue(next, test, testInstanceValue)
		}

		// clear
		require.NotPanics(t, func() {
			vs.ClearValues()
		}, fmt.Sprintf("instance clear values: %s", test))

		// back to zero
		cnt := vs.NumValue()
		require.EqualValues(t, 0, cnt, fmt.Sprintf("cnt(%d) should be zero", cnt))

		test.Append(value, test.value)
	})

	testArrays(inst.GetParentClass(), func(test TestValue, vs meta.Values) {
		// classes disallow set values
		value := reflect.ValueOf(vs)
		require.Panics(t, func() {
			test.Append(value, test.value)
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

func (m MethodMaker) SetTo(value reflect.Value, v interface{}) (err error) {
	name := "Set" + m.String()
	if n := value.MethodByName(name); !n.IsValid() {
		panic(fmt.Sprintf("method didnt exist %s", name))
	} else {
		ret := n.Call([]reflect.Value{reflect.ValueOf(v)})
		if v := ret[0]; !v.IsNil() {
			err = v.Interface().(error)
		}
	}
	return
}

func (m MethodMaker) Append(value reflect.Value, v interface{}) (err error) {
	name := "Append" + m.String()
	if n := value.MethodByName(name); !n.IsValid() {
		panic(fmt.Sprintf("method didnt exist %s", name))
	} else {
		ret := n.Call([]reflect.Value{reflect.ValueOf(v)})
		if v := ret[0]; !v.IsNil() {
			err = v.Interface().(error)
		}
	}
	return
}
