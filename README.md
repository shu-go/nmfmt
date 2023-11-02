Package nmfmt wraps fmt.Xprintf functions providing $name style placeholders.

[![](https://godoc.org/github.com/shu-go/nmfmt?status.svg)](https://godoc.org/github.com/shu-go/nmfmt)
[![Go Report Card](https://goreportcard.com/badge/github.com/shu-go/nmfmt)](https://goreportcard.com/report/github.com/shu-go/nmfmt)
![MIT License](https://img.shields.io/badge/License-MIT-blue)

# nmfmt

Each nmfmt.Xprintf function has a signature like:

 func Printf(format string, a...any) (int, error)

and prints according to a format that contains $name style placeholders.

See examples below.

## Examples

```go
func Example() {
	nmfmt.Printf("$name is $age years old.\n", "name", "Kim", "age", 22)

	nmfmt.Printf("$name ${ name } $name:q ${name:q}aaa\n", "name", "Kim", "age", 22)

	// Output:
	// Kim is 22 years old.
	// Kim Kim "Kim" "Kim"aaa
}

func Example_map() {
	nmfmt.Printf("$name is $age years old.\n",
		nmfmt.M{
			"name": "Kim",
			"age":  22,
		})

	nmfmt.Printf("$name ${ name } $name:q ${name:q}aaa\n",
		nmfmt.M{
			"name": "Kim",
			"age":  22,
		})

	// Output:
	// Kim is 22 years old.
	// Kim Kim "Kim" "Kim"aaa
}

func Example_debug() {
	nmfmt.Printf("$=greeting:q, $=name\n", "name", "Kim", "greeting", "Hello")

	// Output:
	// greeting="Hello", name=Kim
}

func ExampleStruct() {
	nmfmt.Printf("$Name is $Age years old.\n", nmfmt.Struct(struct {
		Name string
		Age  int
	}{Name: "Kim", Age: 22})...)

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

### debug notation

If a placeholder starts with `$=`, the output starts with the name of the placeholder followed by `=`.

`$=name` -> `name=NAME_VALUE`

## Performance

nmfmt (nm) V.S. fmt (std)

About 2 times slower.

```
BenchmarkFprintf/std-16                 15291902                68.96 ns/op            8 B/op          0 allocs/op
BenchmarkFprintf/nm-16                  10768531               104.7 ns/op             8 B/op          0 allocs/op
BenchmarkSprintf/std-16                 13844199                83.69 ns/op           56 B/op          2 allocs/op
BenchmarkSprintf/nm-16                   9587182               121.4 ns/op            56 B/op          2 allocs/op
BenchmarkArgType/Map-16                  4894358               246.4 ns/op           392 B/op          4 allocs/op
BenchmarkArgType/Slice-16                9449457               123.4 ns/op            56 B/op          2 allocs/op
BenchmarkArgType/Struct-16               3580132               336.6 ns/op           288 B/op         10 allocs/op
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
