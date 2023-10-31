package nmfmt_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/shu-go/gotwant"
	"github.com/shu-go/nmfmt"
)

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

func TestVSStd(t *testing.T) {
	cases := []struct {
		stdinput string
		stdargs  []any

		nminput string
		nmargs  map[string]any

		inconpatible bool
		desc         string
	}{
		{
			stdinput: "",
			nminput:  "",
		},
		{
			stdinput: "a",
			nminput:  "a",
		},
		{
			stdinput: "\n",
			nminput:  "\n",
		},
		{
			stdinput: "hello\n",
			nminput:  "hello\n",
		},
		{
			stdinput: "hello\nworld",
			nminput:  "hello\nworld",
		},
		{
			inconpatible: true,
			desc:         "extra",
			stdinput:     "a",
			stdargs:      []any{"hoge"},
			nminput:      "a",
			nmargs:       map[string]any{"1": "hoge"},
		},
		{
			stdinput: "hello, %s",
			stdargs:  []any{"Player"},
			nminput:  "hello, $Name",
			nmargs:   map[string]any{"Name": "Player"},
		},
		{
			stdinput: "%s, %s",
			stdargs:  []any{"Hello", "Player"},
			nminput:  "$Greeting, $Name",
			nmargs:   map[string]any{"Greeting": "Hello", "Name": "Player"},
		},
		{
			stdinput: "%[1]s, %[2]s",
			stdargs:  []any{"Hello", "Player"},
			nminput:  "$Greeting, $Name",
			nmargs:   map[string]any{"Greeting": "Hello", "Name": "Player"},
		},
		{
			desc:     "positional",
			stdinput: "%[2]s, %[1]s",
			stdargs:  []any{"Hello", "Player"},
			nminput:  "$Name, $Greeting",
			nmargs:   map[string]any{"Greeting": "Hello", "Name": "Player"},
		},
		{
			desc:     "position repeated",
			stdinput: "%[2]s, %[1]s, %[2]s",
			stdargs:  []any{"Hello", "Player"},
			nminput:  "$Name, $Greeting, $Name",
			nmargs:   map[string]any{"Greeting": "Hello", "Name": "Player"},
		},
		{
			desc:     "position repeated",
			stdinput: "%[1]s, %[1]q, %[2]d",
			stdargs:  []any{"Hello", 42},
			nminput:  "$Greeting, $Greeting:q, $ID",
			nmargs:   map[string]any{"Greeting": "Hello", "ID": 42},
		},
		{
			desc:     "verb",
			stdinput: "%[1]s, %[1]q, %[2]d",
			stdargs:  []any{"Hello", 42},
			nminput:  "${ Greeting }, ${ Greeting : q }, ${ID}",
			nmargs:   map[string]any{"Greeting": "Hello", "ID": 42},
		},
		{
			inconpatible: true,
			desc:         "missing arg",
			stdinput:     "%[1]s, %[1]q, %[2]d",
			stdargs:      []any{"Hello"},
			nminput:      "${ Greeting }, ${ Greeting : q }, ${ID}",
			nmargs:       map[string]any{"Greeting": "Hello"},
		},
	}

	// also shows how they are inconpatible
	t.Run("Fprintf", func(t *testing.T) {
		for _, c := range cases {
			stdb := &bytes.Buffer{}
			nmb := &bytes.Buffer{}

			fmt.Fprintf(stdb, c.stdinput, c.stdargs...)
			nmfmt.Fprintf(nmb, c.nminput, c.nmargs)

			if c.inconpatible {
				fmt.Fprintf(os.Stderr, "%s\nstd: %s\nnm: %s\n", c.desc, stdb.String(), nmb.String())
			} else {
				gotwant.Test(t, stdb.Bytes(), nmb.Bytes(), gotwant.Format("%q"), gotwant.Desc(c.desc))
			}
		}
	})

	t.Run("Sprintf", func(t *testing.T) {
		for _, c := range cases {

			stds := fmt.Sprintf(c.stdinput, c.stdargs...)
			nms := nmfmt.Sprintf(c.nminput, c.nmargs)

			if !c.inconpatible {
				gotwant.Test(t, stds, nms, gotwant.Format("%q"))
			}
		}
	})

	t.Run("Errorf", func(t *testing.T) {
		for _, c := range cases {

			stds := fmt.Errorf(c.stdinput, c.stdargs...)
			nms := nmfmt.Errorf(c.nminput, c.nmargs)

			if !c.inconpatible {
				gotwant.Test(t, stds, nms, gotwant.Format("%q"))
			}
		}
	})
}

func TestStruct(t *testing.T) {
	want := fmt.Sprintf(
		"%[1]v's name is %[1]q. %[1]v's age is %[2]d, and was born in %[3]d.",
		"Player",
		23,
		time.Now().AddDate(-23, 0, 0).Year())

	f := "$Name's name is $Name:q. $Name's age is $Age, and was born in $Year."
	a := nmfmt.Struct(
		struct {
			Name string
			Age  int
		}{Name: "Player", Age: 23},
		struct {
			Year, Month, Day int
		}{Year: time.Now().AddDate(-23, 0, 0).Year()},
	)

	gotwant.Test(t, nmfmt.Sprintf(f, a), want)
}

func BenchmarkStruct(b *testing.B) {
	s := nmfmt.Sprintf(
		"$Name's age is $Age, and has $Item",
		nmfmt.Struct(struct {
			Name string
			Age  int
		}{Name: "Player", Age: 123},
			struct{ Item string }{Item: "Potion"},
		))
	gotwant.Test(b, s, "Player's age is 123, and has Potion")

	b.Run("nm", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			nmfmt.Sprintf(
				"$Name's age is $Age, and has $Item",
				nmfmt.Struct(struct {
					Name string
					Age  int
				}{Name: "Player", Age: 123},
					struct{ Item string }{Item: "Potion"},
				))
		}
	})
}

func BenchmarkFprintf(b *testing.B) {
	b.Run("std", func(b *testing.B) {
		buf := &bytes.Buffer{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			fmt.Fprintf(buf,
				"%s's age is %d, and has %s",
				"Player", i, "Potion")
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

func BenchmarkSprintf(b *testing.B) {
	b.Run("std", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf(
				"%s's age is %d, and has %s",
				"Player", i, "Potion")
		}
	})

	b.Run("nm", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			nmfmt.Sprintf(
				"$Name's age is $Age, and has $Item",
				nmfmt.Named("Name", "Player", "Age", i, "Item", "Potion"),
			)
		}
	})

	b.Run("std no%", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("hello")
		}
	})

	b.Run("nm no$", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			nmfmt.Sprintf("hello", nil)
		}
	})

	b.Run("std 1", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("hello, %s", "Player")
		}
	})

	b.Run("nm 1", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			nmfmt.Sprintf("hello, $Name", nmfmt.Named("Name", "Player"))
		}
	})

}
