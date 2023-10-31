package nmfmt_test

import (
	"bytes"
	"fmt"
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
