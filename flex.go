package flex

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

var ErrValueIsNotStruct = errors.New("value is not struct")

func StructToMap[T any](v T) (map[string]any, error) {
	m := map[string]any{}
	if err := structToMap(&v, m, map[string]struct{}{}); err != nil {
		return nil, err
	}

	return m, nil
}

func ValuesOf[Target, From any](v From) ([]Target, error) {
	m, err := StructToMap(v)
	if err != nil {
		return nil, fmt.Errorf("struct to map: %w", err)
	}

	return findValuesOfType[Target](m), nil
}

func findValuesOfType[T any](m map[string]any) []T {
	var values []T
	for _, v := range m {
		if v, ok := v.(T); ok {
			values = append(values, v)
		}

		if v, ok := v.(map[string]any); ok {
			values = append(values, findValuesOfType[T](v)...)
		}
	}
	return values
}

func FieldValue[T any](v T, key string) (any, bool) {
	return getFieldValue(&v, key, map[string]struct{}{})
}

func getFieldValue(v any, key string, visitedTypes map[string]struct{}) (any, bool) {
	if key == "" {
		return nil, false
	}

	refVal := reflect.Indirect(reflect.ValueOf(v))
	if refVal.Kind() != reflect.Struct {
		return nil, false
	}

	typVal := refVal.Type()
	if typVal.NumField() == 0 {
		return nil, false
	}

	beginPtr := refVal.Addr().UnsafePointer()
	typId := typVal.PkgPath() + typVal.Name()

	if isAnonStruct := typId == ""; !isAnonStruct {
		visitedTypes[typId] = struct{}{}
	}

	for i := 0; i < typVal.NumField(); i++ {
		fieldInfo := typVal.Field(i)
		if fieldInfo.Name == "_" {
			continue
		}

		key, ok := strings.CutPrefix(key, fieldInfo.Name)
		if !ok {
			continue
		}

		if len(key) == 1 {
			return nil, false
		}

		if len(key) != 0 {
			if fieldInfo.Type.Kind() == reflect.Struct {
				key = key[1:]
			} else {
				return nil, false
			}
		}

		ptr := unsafe.Pointer(uintptr(beginPtr) + fieldInfo.Offset)
		fieldRefVal := reflect.NewAt(fieldInfo.Type, ptr)
		if len(key) == 0 {
			return fieldRefVal.Elem().Interface(), true
		}

		if fieldRefVal.Elem().Kind() != reflect.Struct {
			fieldRefVal = fieldRefVal.Elem()
		}

		if fieldInfo.Type.Kind() == reflect.Pointer && !fieldRefVal.IsNil() {
			elemTyp := fieldInfo.Type.Elem()

			if elemTyp.Kind() == reflect.Struct {
				typId := fieldInfo.PkgPath + fieldInfo.Type.Name()

				if _, ok := visitedTypes[typId]; !ok {
					visitedTypes[typId] = struct{}{}
					fieldInfo.Type = elemTyp
				}
			}
		}

		if fieldInfo.Type.Kind() == reflect.Struct {
			return getFieldValue(fieldRefVal.Interface(), key, visitedTypes)
		}
	}

	return nil, false
}

func structToMap(v any, m map[string]any, visitedTypes map[string]struct{}) error {
	refVal := reflect.Indirect(reflect.ValueOf(v))
	if refVal.Kind() != reflect.Struct {
		return ErrValueIsNotStruct
	}

	typVal := refVal.Type()
	if typVal.NumField() == 0 {
		return nil
	}

	beginPtr := refVal.Addr().UnsafePointer()
	typId := typVal.PkgPath() + typVal.Name()

	if isAnonStruct := typId == ""; !isAnonStruct {
		visitedTypes[typId] = struct{}{}
	}

	for i := 0; i < typVal.NumField(); i++ {
		fieldInfo := typVal.Field(i)
		if fieldInfo.Name == "_" {
			continue
		}

		ptr := unsafe.Pointer(uintptr(beginPtr) + fieldInfo.Offset)
		fieldRefVal := reflect.NewAt(fieldInfo.Type, ptr)
		if fieldRefVal.Elem().Kind() != reflect.Struct {
			fieldRefVal = fieldRefVal.Elem()
		}

		if fieldInfo.Type.Kind() == reflect.Pointer && !fieldRefVal.IsNil() {
			elemTyp := fieldInfo.Type.Elem()

			if elemTyp.Kind() == reflect.Struct {
				typId := fieldInfo.PkgPath + fieldInfo.Type.Name()

				if _, ok := visitedTypes[typId]; !ok {
					visitedTypes[typId] = struct{}{}
					fieldInfo.Type = elemTyp
				}
			}
		}

		fieldValue := fieldRefVal.Interface()

		if fieldInfo.Type.Kind() == reflect.Struct {
			innerMap := map[string]any{}
			m[fieldInfo.Name] = innerMap
			structToMap(fieldValue, innerMap, visitedTypes)
			continue
		}

		m[fieldInfo.Name] = fieldValue
	}

	return nil
}
