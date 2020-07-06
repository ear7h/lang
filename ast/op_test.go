package ast_test

import (
	"testing"

	"github.com/ear7h/lang/ast"
	"github.com/ear7h/lang/ast/parser"
)

func TestParseUnaryExpr(t *testing.T) {
	type tcase struct {
		str string
		ok  bool
		err error
		out *ast.UnaryExpr
	}

	fn := func(tc tcase) func(t *testing.T) {
		return func(t *testing.T) {
			initCur, v, ok, err :=
				parser.DoParseStringForTest(&ast.UnaryExpr{}, tc.str, "test")

			assertErrIs(t, tc.err, err)
			assertEq(t, tc.ok, ok)
			if !ok || err != nil {
				return
			}

			n := v.(*ast.UnaryExpr)

			assertEq(t, initCur.FileInfo(), n.FileInfo())
			assertEq(t, tc.out.Op, n.Op)
			assertEq(t, tc.out.Operand, n.Operand)
		}
	}

	tcases := map[string]tcase{
		"pos": tcase{
			str: `+123`,
			ok:  true,
			out: &ast.UnaryExpr{
				Op: '+',
				Operand: parser.MustParseString(
					&ast.NumberLiteral{},
					"123",
				),
			},
		},
		"neg": tcase{
			str: `-123`,
			ok:  true,
			out: &ast.UnaryExpr{
				Op: '-',
				Operand: parser.MustParseString(
					&ast.NumberLiteral{},
					"123",
				),
			},
		},
		"deref": tcase{
			str: `*123`,
			ok:  true,
			out: &ast.UnaryExpr{
				Op: '*',
				Operand: parser.MustParseString(
					&ast.NumberLiteral{},
					"123",
				),
			},
		},
		"addr": tcase{
			str: `&123`,
			ok:  true,
			out: &ast.UnaryExpr{
				Op: '&',
				Operand: parser.MustParseString(
					&ast.NumberLiteral{},
					"123",
				),
			},
		},
		"not": tcase{
			str: `!123`,
			ok:  true,
			out: &ast.UnaryExpr{
				Op: '!',
				Operand: parser.MustParseString(
					&ast.NumberLiteral{},
					"123",
				),
			},
		},
		"fail no value": tcase{
			str: `123`,
			ok:  false,
		},
		"nested": tcase{
			str: `+-123`,
			ok:  true,
			out: &ast.UnaryExpr{
				Op: '+',
				Operand: &ast.UnaryExpr{
					Op: '-',
					Operand: parser.MustParseString(
						&ast.NumberLiteral{},
						"123",
					),
				},
			},
		},
	}

	for k, v := range tcases {
		t.Run(k, fn(v))
	}

}

func TestParseBinaryExpr(t *testing.T) {
	type tcase struct {
		str string
		ok  bool
		err error
		out interface{}
	}

	fn := func(tc tcase) func(t *testing.T) {
		return func(t *testing.T) {
			initCur, v, ok, err :=
				parser.DoParseStringForTest(&ast.BinaryExpr{}, tc.str, "test")

			assertErrIs(t, tc.err, err)
			assertEq(t, tc.ok, ok)
			if !ok || err != nil {
				return
			}

			n, ok := v.(*ast.BinaryExpr)
			if !ok {
				// number literals by themselves
				assertEq(t, tc.out, v)
				return
			}

			tcout := tc.out.(*ast.BinaryExpr)

			assertEq(t, initCur.FileInfo(), n.FileInfo())
			assertEq(t, tcout.Op, n.Op)
			assertEq(t, tcout.Left, n.Left)
			assertEq(t, tcout.Right, n.Right)
		}
	}

	tcases := map[string]tcase{
		"plain literal": tcase{
			str: `1`,
			ok:  true,
			out: parser.MustParseString(
				&ast.NumberLiteral{},
				"1",
			),
		},
		"shr": tcase{
			str: `1>>123`,
			ok:  true,
			out: &ast.BinaryExpr{
				Op: ">>",
				Left: parser.MustParseString(
					&ast.NumberLiteral{},
					"1",
				),
				Right: parser.MustParseString(
					&ast.NumberLiteral{},
					"123",
				),
			},
		},
		"shr space": tcase{
			str: `1 >> 123`,
			ok:  true,
			out: &ast.BinaryExpr{
				Op: ">>",
				Left: parser.MustParseString(
					&ast.NumberLiteral{},
					"1",
				),
				Right: parser.MustParseString(
					&ast.NumberLiteral{},
					"123",
				),
			},
		},
		"mul": tcase{
			str: `1<<2*3<<4`,
			ok:  true,
			out: &ast.BinaryExpr{
				Op: "*",
				Left: &ast.BinaryExpr{
					Op: "<<",
					Left: parser.MustParseString(
						&ast.NumberLiteral{},
						"1",
					),
					Right: parser.MustParseString(
						&ast.NumberLiteral{},
						"2",
					),
				},
				Right: &ast.BinaryExpr{
					Op: "<<",
					Left: parser.MustParseString(
						&ast.NumberLiteral{},
						"3",
					),
					Right: parser.MustParseString(
						&ast.NumberLiteral{},
						"4",
					),
				},
			},
		},
		"mul1": tcase{
			str: `1*2`,
			ok:  true,
			out: &ast.BinaryExpr{
				Op: "*",
				Left: parser.MustParseString(
					&ast.NumberLiteral{},
					"1",
				),
				Right: parser.MustParseString(
					&ast.NumberLiteral{},
					"2",
				),
			},
		},
		"add": tcase{
			str: `1<<2*3+4`,
			ok:  true,
			out: &ast.BinaryExpr{
				Op: "+",
				Left: &ast.BinaryExpr{
					Op: "*",
					Left: &ast.BinaryExpr{
						Op: "<<",
						Left: parser.MustParseString(
							&ast.NumberLiteral{},
							"1",
						),
						Right: parser.MustParseString(
							&ast.NumberLiteral{},
							"2",
						),
					},
					Right: parser.MustParseString(
						&ast.NumberLiteral{},
						"3",
					),
				},
				Right: parser.MustParseString(
					&ast.NumberLiteral{},
					"4",
				),
			},
		},
	}

	for k, v := range tcases {
		t.Run(k, fn(v))
	}

}
