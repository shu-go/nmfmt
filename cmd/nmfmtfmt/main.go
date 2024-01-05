package main

import (
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/shu-go/gli/v2"
	"github.com/shu-go/nmfmt"
	"golang.org/x/tools/go/ast/astutil"
)

type globalCmd struct {
}

func (c globalCmd) Run(args []string) error {
	if len(args) == 0 {
		return errors.New("1 arg is required")
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, args[0], nil, parser.AllErrors|parser.DeclarationErrors|parser.ParseComments)
	if err != nil {
		return err
	}

	pkgname := "nmfmt"
	if imp := find(f, func(n *ast.ImportSpec) bool {
		return strings.HasPrefix(n.Path.Value, `"github.com/shu-go/nmfmt`)
	}); imp != nil && imp.Name != nil {
		pkgname = imp.Name.Name
	}

	//nodes, err := q.Select(`//*[@type="CallExpr" and Fun[@type="SelectorExpr" and X/@Name="` + pkgname + `" and contains(Sel/@Name, "rint")]]`)
	nodes := findAll(f, func(n *ast.CallExpr) bool {
		fun := conv[ast.SelectorExpr](n.Fun)
		x := conv[ast.Ident](fun.X)
		sel := conv[ast.Ident](fun.Sel)
		return x.Name == pkgname && strings.Contains(sel.Name, "rint")
	})

	changed := false

	for _, n := range nodes {
		// filtering nmfmt.xxx("literal", nmfmt.M)
		//                     ^^^^^^^^^^^^^^^^^^
		// filtering nmfmt.xxx("literal")
		//                     ^^^^^^^^^

		if len(n.Args) == 0 || len(n.Args) > 2 {
			continue
		}
		arg0, ok0 := n.Args[0].(*ast.BasicLit)
		if !ok0 || arg0.Kind != token.STRING {
			continue
		}

		names := nmfmt.ExtractNames(arg0.Value)

		if len(names) > 0 && len(n.Args) != 2 {
			n.Args = append(n.Args, &ast.CompositeLit{
				Type: &ast.SelectorExpr{
					X: &ast.Ident{
						Name:    pkgname,
						NamePos: 1,
					},
					Sel: &ast.Ident{
						Name:    "M",
						NamePos: 1,
					},
				},
				Elts: []ast.Expr{},
			})
			changed = true
		}
		arg1, ok1 := n.Args[1].(*ast.CompositeLit)
		if !ok1 {
			continue
		}
		m, okm := arg1.Type.(*ast.SelectorExpr)
		if !okm {
			continue
		}
		if x, ok := m.X.(*ast.Ident); !ok || x.Name != pkgname || m.Sel.Name != "M" {
			continue
		}

		kvs := make([]struct {
			Key   string
			Value ast.Node
		}, len(arg1.Elts))

		for _, e := range arg1.Elts {
			kv, okkv := e.(*ast.KeyValueExpr)
			if !okkv {
				continue
			}
			key, okkey := kv.Key.(*ast.BasicLit)
			if !okkey {
				continue
			}

			kvs = append(kvs, struct {
				Key   string
				Value ast.Node
			}{
				Key:   key.Value[1 : len(key.Value)-1],
				Value: kv.Value,
			})
		}

		for name := range names {
			found := false
			for _, kv := range kvs {
				if kv.Key == name {
					found = true
					break
				}
			}

			if !found {
				arg1.Elts = append(arg1.Elts, &ast.KeyValueExpr{
					Key: &ast.BasicLit{
						Kind:  token.STRING,
						Value: `"` + name + `"`,
					},
					Value: &ast.Ident{
						Name: name,
					},
				})
				changed = true
			}
		}
		//debug
		/*
			for i, e := range arg1.Elts {
				fmt.Fprintf(os.Stderr, "%v %#v %#v\n", i, e.(*ast.KeyValueExpr).Key, e.(*ast.KeyValueExpr).Value)
			}
		*/
	}

	if !changed {
		return nil
	}

	out, err := os.Create(args[0])
	if err != nil {
		return err
	}
	defer out.Close()

	err = format.Node(out, fset, f)
	if err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func find[T any](root ast.Node, test func(n T) bool) T {
	var result T
	astutil.Apply(root, func(c *astutil.Cursor) bool {
		if n, ok := c.Node().(T); ok {
			if test(n) {
				result = n
				return false
			}
		}
		return true
	}, nil)
	return result
}

func findAll[T any](root ast.Node, test func(n T) bool) []T {
	var result []T
	astutil.Apply(root, func(c *astutil.Cursor) bool {
		if n, ok := c.Node().(T); ok {
			if test(n) {
				if result == nil {
					result = make([]T, 0, 1)
				}
				result = append(result, n)
			}
		}
		return true
	}, nil)
	return result
}

func conv[T any, P *T](n ast.Node) P {
	if c, ok := n.(P); ok {
		return c
	}

	var zero T
	return &zero
}

////////////////////////////////////////////////////////////////////////////////

// Version is app version
var Version string

func main() {
	app := gli.NewWith(&globalCmd{})
	app.Name = "nmfmtfmt"
	app.Desc = "add nmfmt.M"
	app.Version = Version
	app.Usage = ``
	app.Copyright = "(C) 2024 Shuhei Kubota"
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
