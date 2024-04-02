package flex

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type embedTest struct {
	f1 string
	f2 string
	f3 int
}

type embedTest2 struct{}

type fooTest struct {
	embedTest
	*embedTest2

	f1  int
	f2  *int
	f3  string
	f4  *string
	f5  [3]*int
	f6  []byte
	f7  func()
	f8  any
	f9  chan int
	f10 map[string]any
	f11 struct {
		f12 int
		f13 *int
		f14 struct {
			f15 int
		}
	}
	f12 *struct{ f13 int }
	f13 *fooTest
	f14 *fooTest
	_   int
	_   string
}

func TestStructToMap(t *testing.T) {
	foo := fooTest{
		embedTest: embedTest{
			f1: "f1",
			f2: "f2",
			f3: 3,
		},
		embedTest2: nil,
		f1:         1,
		f2:         ptrOf(2),
		f3:         "hello",
		f4:         ptrOf("world"),
		f5:         [3]*int{ptrOf(1), ptrOf(2)},
		f6:         []byte{1, 2, 3, 4},
		f7:         func() {},
		f8:         "string",
		f9:         make(chan int),
		f10:        map[string]any{"foo": "bar"},
		f11: struct {
			f12 int
			f13 *int
			f14 struct{ f15 int }
		}{
			f12: 12,
			f13: nil,
			f14: struct{ f15 int }{
				f15: 15,
			},
		},
		f12: &struct{ f13 int }{
			f13: 123,
		},
	}
	foo.f13 = &foo
	foo.f14 = &foo
	m, err := StructToMap(foo)
	if err != nil {
		t.Fatal(err)
	}

	assertEqualMaps(t, map[string]any{
		"embedTest": map[string]any{
			"f1": "f1",
			"f2": "f2",
			"f3": 3,
		},
		"embedTest2": (*embedTest2)(nil),
		"f1":         foo.f1,
		"f2":         foo.f2,
		"f3":         foo.f3,
		"f4":         foo.f4,
		"f5":         foo.f5,
		"f6":         foo.f6,
		"f7":         foo.f7,
		"f8":         foo.f8,
		"f9":         foo.f9,
		"f10":        foo.f10,
		"f11": map[string]any{
			"f12": foo.f11.f12,
			"f13": foo.f11.f13,
			"f14": map[string]any{
				"f15": foo.f11.f14.f15,
			},
		},
		"f12": map[string]any{
			"f13": foo.f12.f13,
		},
		"f13": &foo,
		"f14": &foo,
	}, m)
}

func TestFieldValue(t *testing.T) {
	type f3 struct {
		f4 string
	}

	type f2 struct {
		f3 f3
	}

	type f1 struct {
		f2 f2
	}

	type composite struct {
		f1 f1
	}

	testCases := []struct {
		Value         composite
		Key           string
		ExpectedValue any
		ExpectedOk    bool
	}{

		{
			Value:         composite{f1: f1{f2: f2{f3: f3{f4: "hello"}}}},
			Key:           "",
			ExpectedValue: nil,
			ExpectedOk:    false,
		},
		{
			Value:         composite{f1: f1{f2: f2{f3: f3{f4: "hello"}}}},
			Key:           "f1.",
			ExpectedValue: nil,
			ExpectedOk:    false,
		},
		{
			Value:         composite{f1: f1{f2: f2{f3: f3{f4: "hello"}}}},
			Key:           ".f1",
			ExpectedValue: nil,
			ExpectedOk:    false,
		},
		{
			Value:         composite{f1: f1{f2: f2{f3: f3{f4: "hello"}}}},
			Key:           ".",
			ExpectedValue: nil,
			ExpectedOk:    false,
		},
		{
			Value:         composite{f1: f1{f2: f2{f3: f3{f4: "hello"}}}},
			Key:           "f4",
			ExpectedValue: nil,
			ExpectedOk:    false,
		},
		{
			Value:         composite{f1: f1{f2: f2{f3: f3{f4: "hello"}}}},
			Key:           "f1",
			ExpectedValue: f1{f2: f2{f3: f3{f4: "hello"}}},
			ExpectedOk:    true,
		},
		{
			Value:         composite{f1: f1{f2: f2{f3: f3{f4: "hello"}}}},
			Key:           "f1.f2",
			ExpectedValue: f2{f3: f3{f4: "hello"}},
			ExpectedOk:    true,
		},
		{
			Value:         composite{f1: f1{f2: f2{f3: f3{f4: "hello"}}}},
			Key:           "f1.f2.f3",
			ExpectedValue: f3{f4: "hello"},
			ExpectedOk:    true,
		},
		{
			Value:         composite{f1: f1{f2: f2{f3: f3{f4: "hello"}}}},
			Key:           "f1.f2.f3.f4",
			ExpectedValue: "hello",
			ExpectedOk:    true,
		},
	}

	for _, testCase := range testCases {
		val, ok := FieldValue(testCase.Value, testCase.Key)
		assert.Equal(t, testCase.ExpectedValue, val)
		assert.Equal(t, testCase.ExpectedOk, ok)
	}
}

func TestValuesOf(t *testing.T) {
	t.Run("flat", valuesOfTest(
		struct {
			f1 string
			f2 string
			f3 int
		}{f1: "hello", f2: "world", f3: 111},
		[]string{"hello", "world"}, nil,
	))

	t.Run("repeated", valuesOfTest(
		struct {
			f1 string
			f2 string
			f3 string
		}{f1: "hello", f2: "hello", f3: "world"},
		[]string{"hello", "hello", "world"}, nil,
	))

	t.Run("nested", valuesOfTest(
		struct {
			n1 struct {
				f1 string
			}
			n2 struct {
				f1 string
			}
		}{n1: struct{ f1 string }{f1: "hello"}, n2: struct{ f1 string }{f1: "world"}},
		[]string{"hello", "world"}, nil,
	))

	type foo struct {
		f1 string
	}

	t.Run("any", valuesOfTest(
		struct {
			f1 any
			f2 any
			f3 any
			f4 any
		}{f1: 1, f2: "hello", f3: 2.2, f4: foo{f1: "world"}}, []any{1, "hello", 2.2, foo{f1: "world"}}, nil,
	))

	t.Run("nostruct", valuesOfTest[any]("nostruct", nil, ErrValueIsNotStruct))
}

func valuesOfTest[Target, From any](val From, expected []Target, expectedErr error) func(t *testing.T) {
	intersection := func(v1, v2 []Target) int {
		m := map[any]int{}
		for _, v := range v1 {
			c := m[v]
			m[v] = c + 1
		}

		res := 0
		for _, v := range v2 {
			c := m[v]
			if c == 0 {
				continue
			}

			res++
			m[v] = c - 1
		}

		return res
	}

	return func(t *testing.T) {
		values, err := ValuesOf[Target](val)
		if expectedErr == nil {
			assert.NoError(t, err)
		} else {
			assert.ErrorIs(t, err, expectedErr)
		}
		assert.Equal(t, len(expected), intersection(values, expected))
	}
}

func assertEqualMaps(t *testing.T, expected, actual map[string]any) {
	if len(expected) != len(actual) {
		t.Fatal("expect same size of maps")
	}

	diffValues, onlyInM1, onlyInM2 := diffKeys("$", expected, actual)
	if len(onlyInM1) != 0 {
		t.Fail()
		t.Logf("actual has no keys from expected: %v", onlyInM1)
	}

	if len(onlyInM2) != 0 {
		t.Fail()
		t.Logf("expected has no keys from actual: %v", onlyInM2)
	}

	if len(diffValues) != 0 {
		t.Fail()
		t.Log("found keys which have different values")
		for k, diff := range diffValues {
			t.Logf("key: %s; expected: %+v; actual: %+v;", k, diff.left, diff.right)
		}
	}
}

type diffValue struct {
	left  any
	right any
}

func diffKeys(rootKey string, m1, m2 map[string]any) (map[string]diffValue, []string, []string) {
	diffValues := map[string]diffValue{}
	onlyInM1 := []string{}
	onlyInM2 := []string{}

	for k, v1 := range m1 {
		fullKey := rootKey + "." + k
		v2, ok := m2[k]
		if !ok {
			onlyInM1 = append(onlyInM1, fullKey)
			continue
		}

		vm1, ok1 := v1.(map[string]any)
		vm2, ok2 := v2.(map[string]any)
		if ok1 && ok2 {
			innerDiffValues, innerOnlyInM1, innerOnlyInM2 := diffKeys(fullKey, vm1, vm2)
			for innerK, innerV := range innerDiffValues {
				diffValues[innerK] = innerV
			}
			onlyInM1 = append(onlyInM1, innerOnlyInM1...)
			onlyInM2 = append(onlyInM2, innerOnlyInM2...)
			continue
		}

		refV1 := reflect.ValueOf(v1)
		refV2 := reflect.ValueOf(v2)

		if refV1.Kind() != refV2.Kind() {
			diffValues[fullKey] = diffValue{
				left:  v1,
				right: v2,
			}
			continue
		}

		if refV1.Kind() == reflect.Func {
			if refV1.UnsafePointer() == refV2.UnsafePointer() {
				continue
			}

			diffValues[fullKey] = diffValue{
				left:  v1,
				right: v2,
			}
			continue
		}

		if !reflect.DeepEqual(v1, v2) {
			diffValues[fullKey] = diffValue{
				left:  v1,
				right: v2,
			}
			continue
		}
	}

	for k := range m2 {
		fullKey := rootKey + "." + k
		if _, ok := m2[k]; !ok {
			onlyInM1 = append(onlyInM2, fullKey)
		}
	}

	return diffValues, onlyInM1, onlyInM2
}

func ptrOf[T any](v T) *T {
	return &v
}
