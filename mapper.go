package nmfmt

import (
	"reflect"
)

// Named returns a map[string]any {arg1: arg2, arg3: arg4, ... }
func Named(nameAndValue ...any) map[string]any {
	if len(nameAndValue) < 2 {
		return nil
	}

	m := make(map[string]any)

	var name string
	for i := 0; i < len(nameAndValue); i++ {
		if i%2 == 0 {
			name = nameAndValue[i].(string)
			continue
		}

		m[name] = nameAndValue[i]
	}

	return m
}

// Struct returns a map of names and values of fields from structs.
//
// If a key is duplicated among structs, the first found element wins.
//
// Note: An unexported field results in <nil>. (NO: name string; YES: Name string)
func Struct(structs ...any) map[string]any {
	m := make(map[string]any)

	for i := len(structs) - 1; i >= 0; i-- {
		s := structs[i]
		v := reflect.Indirect(reflect.ValueOf(s))
		if v.Type().Kind() != reflect.Struct {
			continue
		}

		for i := 0; i < v.NumField(); i++ {
			ft := v.Type().Field(i)
			fv := v.Field(i)

			if ft.IsExported() {
				m[ft.Name] = fv.Interface()
			}
		}
	}

	return m
}
