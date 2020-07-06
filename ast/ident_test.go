package ast_test

import (
	"testing"

	"github.com/ear7h/lang/ast"
	"github.com/ear7h/lang/ast/parser"
)

func TestParseIdent(t *testing.T) {
	type tcase struct {
		str string
		ok  bool
		err error
		out *ast.Ident
	}

	fn := func(tc tcase) func(t *testing.T) {
		return func(t *testing.T) {
			initCur, v, ok, err :=
				parser.DoParseStringForTest(&ast.Ident{}, tc.str, "test")

			assertErrIs(t, tc.err, err)
			assertEq(t, tc.ok, ok)
			if !ok || err != nil {
				return
			}

			n := v.(*ast.Ident)

			assertEq(t, initCur.FileInfo(), n.FileInfo())
			assertEq(t, tc.out.Name, n.Name)
			assertEq(t, tc.out.IsExported, n.IsExported)
		}
	}

	tcases := map[string]tcase{
		"not exported": tcase{
			str: "hello",
			ok:  true,
			out: &ast.Ident{
				IsExported: false,
				Name:       "hello",
			},
		},
		"exported": tcase{
			str: "Hello",
			ok:  true,
			out: &ast.Ident{
				IsExported: true,
				Name:       "Hello",
			},
		},
		"with number": tcase{
			str: "hello1",
			ok:  true,
			out: &ast.Ident{
				IsExported: false,
				Name:       "hello1",
			},
		},
		"with underscore": tcase{
			str: "hello_",
			ok:  true,
			out: &ast.Ident{
				IsExported: false,
				Name:       "hello_",
			},
		},
		"with underscore 1": tcase{
			str: "_",
			ok:  true,
			out: &ast.Ident{
				IsExported: false,
				Name:       "_",
			},
		},
		"with underscore 2": tcase{
			str: "_asd",
			ok:  true,
			out: &ast.Ident{
				IsExported: false,
				Name:       "_asd",
			},
		},
		"fail 1": tcase{
			str: "1hello",
			ok:  false,
		},
	}

	for k, v := range tcases {
		t.Run(k, fn(v))
	}
}

func TestParseObjExpr(t *testing.T) {
	type tcase struct {
		str string
		ok  bool
		err error
		out *ast.ObjExpr
	}

	fn := func(tc tcase) func(t *testing.T) {
		return func(t *testing.T) {
			initCur, v, ok, err :=
				parser.DoParseStringForTest(
					&ast.ExprOperandParser{}, tc.str, "test")

			assertErrIs(t, tc.err, err)
			assertEq(t, tc.ok, ok)
			if !ok || err != nil {
				return
			}

			n := v.(*ast.ObjExpr)

			assertEq(t, initCur.FileInfo(), n.FileInfo())
			assertEq(t, tc.out.Object, n.Object)
			assertEq(t, tc.out.Op, n.Op)
			assertEq(t, tc.out.Arg, n.Arg)
			assertEq(t, tc.out.Right, n.Right)
		}
	}

	tcases := map[string]tcase{
		"plain": tcase{
			str: "hello.world",
			ok:  true,
			out: &ast.ObjExpr{
				Object: &ast.Ident{
					IsExported: false,
					Name:       "hello",
				},
				Op: ast.ObjField,
				Arg: &ast.Ident{
					IsExported: false,
					Name: "world",
				},
			},
		},
		"literal": tcase{
			str: "1.world",
			ok:  true,
			out: &ast.ObjExpr{
				Object: parser.MustParseString(
					ast.LiteralParser{},
					"1",
				),
				Op: ast.ObjField,
				Arg: &ast.Ident{
					IsExported: false,
					Name:       "world",
				},
			},
		},
		"expr": tcase{
			str: "(1 + 1).world",
			ok:  true,
			out: &ast.ObjExpr{
				Object: parser.MustParseString(
					ast.ExprParser{},
					"1+1",
				),
				Op: ast.ObjField,
				Arg: &ast.Ident{
					IsExported: false,
					Name:       "world",
				},
			},
		},
		"nested": tcase{
			str: "hi.hello.world",
			ok:  true,
			out: &ast.ObjExpr{
				Object: &ast.Ident{
					IsExported: false,
					Name:       "hi",
				},
				Op: ast.ObjField,
				Arg: &ast.Ident{
					IsExported: false,
					Name:       "hello",
				},
				Right: &ast.ObjExpr{
					Op: ast.ObjField,
					Arg: &ast.Ident{
						IsExported: false,
						Name: "world",
					},
				},
			},
		},
	}

	for k, v := range tcases {
		t.Run(k, fn(v))
	}

}
