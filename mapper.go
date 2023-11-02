package nmfmt

import (
	"reflect"
)

// Struct returns a slice of names and values of fields from structs.
//
// If a key is duplicated among structs, the first found element wins.
//
// Note: An unexported field results in <nil>. (NO: name string; YES: Name string)
func Struct(structs ...any) []any {
	a := make([]any, 0, 8)

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
				a = append(a, ft.Name)
				a = append(a, fv.Interface())
			}
		}
	}

	return a
}
