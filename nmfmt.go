// Package nmfmt wraps fmt.Xprintf functions providing $name style placeholders.
//
// Each nmfmt.Xprintf function has a signature like:
//
//	func Printf(format string, m map[string]any) (int, error)
//
// And prints according to a format that contains $name style placeholders.
//
// # Placeholders
//
// Example: `variable1:q`
//
// Each placeholder in the format is like $name, ${name}, $name:verb or ${name:verb}.
// The names are keys of the map m. (case sensitive)
// And their values are to be embedded.
//
// # Name
//
// Must match \w.
//
// See [Named], [Struct]
//
// # Verb
//
// Verb is with `:`.
// Must match \w.
//
// Defaults to `v`.
//
// See https://pkg.go.dev/fmt.
//
// # debug notation
//
// If a placeholder starts with `$=`, the output starts with the name of the placeholder followed by `=`.
//
// `$=name` -> `name=NAME_VALUE`
package nmfmt

import (
	"io"
)

var f Formatter = New()

type M map[string]any

func Printf(format string, a ...any) (int, error) {
	return f.Printf(format, a...)
}

func Fprintf(w io.Writer, format string, a ...any) (int, error) {
	return f.Fprintf(w, format, a...)
}

func Sprintf(format string, a ...any) string {
	return f.Sprintf(format, a...)
}

func Errorf(format string, a ...any) error {
	return f.Errorf(format, a...)
}
