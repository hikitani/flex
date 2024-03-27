package flex

import (
	"reflect"
	"testing"
)

type fooTest struct {
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
}

func TestXxx(t *testing.T) {
	foo := fooTest{
		f1:  1,
		f2:  ptrOf(2),
		f3:  "hello",
		f4:  ptrOf("world"),
		f5:  [3]*int{ptrOf(1), ptrOf(2)},
		f6:  []byte{1, 2, 3, 4},
		f7:  func() {},
		f8:  "string",
		f9:  make(chan int),
		f10: map[string]any{"foo": "bar"},
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
		"f1":  foo.f1,
		"f2":  foo.f2,
		"f3":  foo.f3,
		"f4":  foo.f4,
		"f5":  foo.f5,
		"f6":  foo.f6,
		"f7":  foo.f7,
		"f8":  foo.f8,
		"f9":  foo.f9,
		"f10": foo.f10,
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
