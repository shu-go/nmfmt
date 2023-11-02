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

func TestNotation(t *testing.T) {
	t.Run("Boundary", func(t *testing.T) {
		gotwant.Test(t, nmfmt.Sprintf("hello, $Name.", "Name", "Hoge"), "hello, Hoge.")
		gotwant.Test(t, nmfmt.Sprintf("hello, $Name", "Name", "Hoge"), "hello, Hoge")
		gotwant.Test(t, nmfmt.Sprintf("$Name, hello", "Name", "Hoge"), "Hoge, hello")
		gotwant.Test(t, nmfmt.Sprintf("$Name, hello.", "Name", "Hoge"), "Hoge, hello.")

		gotwant.Test(t, nmfmt.Sprintf("hello, ${Name}.", "Name", "Hoge"), "hello, Hoge.")
		gotwant.Test(t, nmfmt.Sprintf("hello, ${Name}", "Name", "Hoge"), "hello, Hoge")
		gotwant.Test(t, nmfmt.Sprintf("${Name}, hello", "Name", "Hoge"), "Hoge, hello")
		gotwant.Test(t, nmfmt.Sprintf("${Name}, hello.", "Name", "Hoge"), "Hoge, hello.")
	})

	t.Run("Verb", func(t *testing.T) {
		gotwant.Test(t, nmfmt.Sprintf("$Name:q", "Name", "Hoge"), `"Hoge"`)
		gotwant.Test(t, nmfmt.Sprintf("$Name:q.", "Name", "Hoge"), `"Hoge".`)
		gotwant.Test(t, nmfmt.Sprintf("$Name:#v.", "Name", "Hoge"), `"Hoge".`)
		gotwant.Test(t, nmfmt.Sprintf("$Name:+v.", "Name", "Hoge"), `Hoge.`)
		gotwant.Test(t, nmfmt.Sprintf("$=Name:+v.", "Name", "Hoge"), `Name=Hoge.`)

		gotwant.Test(t, nmfmt.Sprintf("${Name:q}", "Name", "Hoge"), `"Hoge"`)
		gotwant.Test(t, nmfmt.Sprintf("${Name:q}.", "Name", "Hoge"), `"Hoge".`)
		gotwant.Test(t, nmfmt.Sprintf("${Name:#v}.", "Name", "Hoge"), `"Hoge".`)
		gotwant.Test(t, nmfmt.Sprintf("${Name:+v}.", "Name", "Hoge"), `Hoge.`)
		gotwant.Test(t, nmfmt.Sprintf("${=Name:+v}.", "Name", "Hoge"), `Name=Hoge.`)

		gotwant.Test(t, nmfmt.Sprintf("${ Name:q }", "Name", "Hoge"), `"Hoge"`)
		gotwant.Test(t, nmfmt.Sprintf("${ Name:q }.", "Name", "Hoge"), `"Hoge".`)
		gotwant.Test(t, nmfmt.Sprintf("${ Name:#v }.", "Name", "Hoge"), `"Hoge".`)
		gotwant.Test(t, nmfmt.Sprintf("${ Name:+v }.", "Name", "Hoge"), `Hoge.`)
		gotwant.Test(t, nmfmt.Sprintf("${ =Name:+v }.", "Name", "Hoge"), `Name=Hoge.`)
	})
}

func TestVSStd(t *testing.T) {
	cases := []struct {
		stdinput string
		stdargs  []any

		nminput string
		nmargs  []any

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
			nmargs:       []any{"1", "hoge"},
		},
		{
			stdinput: "hello, %s",
			stdargs:  []any{"Player"},
			nminput:  "hello, $Name",
			nmargs:   []any{"Name", "Player"},
		},
		{
			stdinput: "%s, %s",
			stdargs:  []any{"Hello", "Player"},
			nminput:  "$Greeting, $Name",
			nmargs:   []any{"Greeting", "Hello", "Name", "Player"},
		},
		{
			stdinput: "%[1]s, %[2]s",
			stdargs:  []any{"Hello", "Player"},
			nminput:  "$Greeting, $Name",
			nmargs:   []any{"Greeting", "Hello", "Name", "Player"},
		},
		{
			desc:     "positional",
			stdinput: "%[2]s, %[1]s",
			stdargs:  []any{"Hello", "Player"},
			nminput:  "$Name, $Greeting",
			nmargs:   []any{"Greeting", "Hello", "Name", "Player"},
		},
		{
			desc:     "position repeated",
			stdinput: "%[2]s, %[1]s, %[2]s",
			stdargs:  []any{"Hello", "Player"},
			nminput:  "$Name, $Greeting, $Name",
			nmargs:   []any{"Greeting", "Hello", "Name", "Player"},
		},
		{
			desc:     "position repeated",
			stdinput: "%[1]s, %[1]q, %[2]d",
			stdargs:  []any{"Hello", 42},
			nminput:  "$Greeting, $Greeting:q, $ID",
			nmargs:   []any{"Greeting", "Hello", "ID", 42},
		},
		{
			desc:     "verb",
			stdinput: "%[1]s, %[1]q, %[2]d",
			stdargs:  []any{"Hello", 42},
			nminput:  "${ Greeting }, ${ Greeting : q }, ${ID}",
			nmargs:   []any{"Greeting", "Hello", "ID", 42},
		},
		{
			inconpatible: true,
			desc:         "missing arg",
			stdinput:     "%[1]s, %[1]q, %[2]d",
			stdargs:      []any{"Hello"},
			nminput:      "${ Greeting }, ${ Greeting : q }, ${ID}",
			nmargs:       []any{"Greeting", "Hello"},
		},
		{
			desc:     "Arg=Arg",
			stdinput: "Greeting=%[1]q",
			stdargs:  []any{"Hello"},
			nminput:  "$=Greeting:q",
			nmargs:   []any{"Greeting", "Hello"},
		},
		{
			desc:     "Arg=Arg",
			stdinput: "Greeting=%[1]q",
			stdargs:  []any{"Hello"},
			nminput:  "${=Greeting:q}",
			nmargs:   []any{"Greeting", "Hello"},
		},
	}

	// also shows how they are inconpatible
	t.Run("Fprintf", func(t *testing.T) {
		for _, c := range cases {
			stdb := &bytes.Buffer{}
			nmb := &bytes.Buffer{}

			fmt.Fprintf(stdb, c.stdinput, c.stdargs...)
			nmfmt.Fprintf(nmb, c.nminput, c.nmargs...)

			if c.inconpatible {
				fmt.Fprintf(os.Stderr, "%s\nstd: %s\nnm: %s\n", c.desc, stdb.String(), nmb.String())
			} else {
				gotwant.Test(t, nmb.Bytes(), stdb.Bytes(), gotwant.Format("%q"), gotwant.Desc(c.desc))
			}
		}
	})

	t.Run("Sprintf", func(t *testing.T) {
		for _, c := range cases {

			stds := fmt.Sprintf(c.stdinput, c.stdargs...)
			nms := nmfmt.Sprintf(c.nminput, c.nmargs...)

			if !c.inconpatible {
				gotwant.Test(t, nms, stds, gotwant.Format("%q"))
			}
		}
	})

	t.Run("Errorf", func(t *testing.T) {
		for _, c := range cases {

			stds := fmt.Errorf(c.stdinput, c.stdargs...)
			nms := nmfmt.Errorf(c.nminput, c.nmargs...)

			if !c.inconpatible {
				gotwant.Test(t, nms, stds, gotwant.Format("%q"))
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

	gotwant.Test(t, nmfmt.Sprintf(f, a...), want)
}

func BenchmarkStruct(b *testing.B) {
	s := nmfmt.Sprintf(
		"$Name's age is $Age, and has $Item",
		nmfmt.Struct(struct {
			Name string
			Age  int
		}{Name: "Player", Age: 123},
			struct{ Item string }{Item: "Potion"},
		)...)
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
				)...)
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
				"Name", "Player", "Age", i, "Item", "Potion",
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
				"Name", "Player", "Age", i, "Item", "Potion",
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
			nmfmt.Sprintf("hello, $Name", "Name", "Player")
		}
	})

}

func BenchmarkMapOrSlice(b *testing.B) {
	b.Run("Map", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			nmfmt.Sprintf(
				"$Name's age is $Age, and has $Item",
				nmfmt.M{
					"Name": "Player",
					"Age":  i,
					"Item": "Potion",
				},
			)
		}
	})

	b.Run("Slice", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			nmfmt.Sprintf(
				"$Name's age is $Age, and has $Item",
				"Name", "Player",
				"Age", i,
				"Item", "Potion",
			)
		}
	})
}
