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
package nmfmt

import (
	"io"
)

var f Formatter = New()

func Printf(format string, m map[string]any) (int, error) {
	return f.Printf(format, m)
}

func Fprintf(w io.Writer, format string, m map[string]any) (int, error) {
	return f.Fprintf(w, format, m)
}

func Sprintf(format string, m map[string]any) string {
	return f.Sprintf(format, m)
}

func Errorf(format string, m map[string]any) error {
	return f.Errorf(format, m)
}
