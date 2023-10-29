Package nmfmt wraps fmt.Xprintf functions providing $name style placeholders.

[![](https://godoc.org/github.com/shu-go/nmfmt?status.svg)](https://godoc.org/github.com/shu-go/nmfmt)
[![Go Report Card](https://goreportcard.com/badge/github.com/shu-go/nmfmt)](https://goreportcard.com/report/github.com/shu-go/nmfmt)
![MIT License](https://img.shields.io/badge/License-MIT-blue)

# nmfmt

Each nmfmt.Xprintf function has a signature like:

 func Printf(format string, m map[string]any) (int, error)

And prints according to a format that contains $name style placeholders.

## Examples

```go
func Example() {
	nmfmt.Printf("$name is $age years old.\n", nmfmt.Named("name", "Kim", "age", 22))

	nmfmt.Printf("$name ${ name } $name:q ${name:q}aaa\n", nmfmt.Named("name", "Kim", "age", 22))

	// Output:
	// Kim is 22 years old.
	// Kim Kim "Kim" "Kim"aaa
}

func ExampleStruct() {
	nmfmt.Printf("$Name is $Age years old.\n", nmfmt.Struct(struct {
		Name string
		Age  int
	}{Name: "Kim", Age: 22}))

	// Output:
	// Kim is 22 years old.
}
```

## Placeholders

Example: `variable1:q`

Each placeholder in the format is like $name, ${name}, $name:verb or ${name:verb}.
The names are keys of the map m. (case sensitive)
And their values are to be embedded.

### Name

Must match \w.

See `Named()` and `Struct()` in the doc.

### Verb

Verb is with `:`.
Must match \w.

Defaults to `v`.

See https://pkg.go.dev/fmt.

## Performance

nmfmt (nm) V.S. fmt (std)

About 2 times slower.

```
BenchmarkStruct/nm-16    3499479               333.2 ns/op           448 B/op          8 allocs/op
BenchmarkFprintf/std-16                 16273215                71.12 ns/op            8 B/op          0 allocs/op
BenchmarkFprintf/nm-16                   8108190               141.6 ns/op             8 B/op          0 allocs/op
BenchmarkSprintf/std-16                 13886348                84.93 ns/op           56 B/op          2 allocs/op
BenchmarkSprintf/nm-16                   7585075               156.1 ns/op            56 B/op          2 allocs/op
PASS
```

### Code (Fprintf)

```go
func BenchmarkFprintf(b *testing.B) {
	b.Run("std", func(b *testing.B) {
		buf := &bytes.Buffer{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			fmt.Fprintf(buf,
				"%s's age is %d, and has %s",
				"Player", i, "Posion")
		}
	})

	b.Run("nm", func(b *testing.B) {
		buf := &bytes.Buffer{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			nmfmt.Fprintf(buf,
				"$Name's age is $Age, and has $Item",
				nmfmt.Named("Name", "Player", "Age", i, "Item", "Potion"),
			)
		}
	})
}
```


----

Copyright 2023 Shuhei Kubota

<!--  vim: set et ft=markdown sts=4 sw=4 ts=4 tw=0 : -->
