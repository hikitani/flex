package flex

import (
	"errors"
	"reflect"
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
		ptr := unsafe.Pointer(uintptr(beginPtr) + fieldInfo.Offset)

		val := (*int)(ptr)
		_ = val

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
