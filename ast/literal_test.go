package ast_test

import (
	"testing"

	"github.com/ear7h/lang/ast"
	"github.com/ear7h/lang/ast/parser"
)

func TestParseNumberLiteral(t *testing.T) {
	type tcase struct {
		str string
		ok  bool
		err error
		out *ast.NumberLiteral
	}

	fn := func(tc tcase) func(t *testing.T) {
		return func(t *testing.T) {
			initCur, v, ok, err :=
				parser.DoParseStringForTest(
					&ast.NumberLiteral{}, tc.str, "test")

			assertErrIs(t, tc.err, err)
			assertEq(t, tc.ok, ok)
			if !ok || err != nil {
				return
			}

			n := v.(*ast.NumberLiteral)

			assertEq(t, initCur.FileInfo(), n.FileInfo())
			assertEq(t, tc.out.Orig, n.Orig)
			assertEq(t, tc.out.Parsed, n.Parsed)
		}
	}

	tcases := map[string]tcase{
		"pos": tcase{
			str: `123`,
			ok:  true,
			out: &ast.NumberLiteral{
				Orig:   `123`,
				Parsed: 123,
			},
		},
	}

	for k, v := range tcases {
		t.Run(k, fn(v))
	}

}

func TestParseStringLiteral(t *testing.T) {
	type tcase struct {
		str string
		ok  bool
		err error
		out *ast.StringLiteral
	}

	fn := func(tc tcase) func(t *testing.T) {
		return func(t *testing.T) {
			initCur, v, ok, err :=
				parser.DoParseStringForTest(
					&ast.StringLiteral{}, tc.str, "test")

			assertErrIs(t, tc.err, err)
			assertEq(t, ok, tc.ok)
			if !ok || err != nil {
				return
			}

			n := v.(*ast.StringLiteral)

			assertEq(t, initCur.FileInfo(), n.FileInfo())
			assertEq(t, tc.out.Orig, n.Orig)
			assertEq(t, tc.out.Parsed, n.Parsed)
		}
	}

	tcases := map[string]tcase{
		"simple": tcase{
			str: `"asd"`,
			ok:  true,
			out: &ast.StringLiteral{
				Orig:   `"asd"`,
				Parsed: "asd",
			},
		},
		"newline": tcase{
			str: `"asd\nqwe"`,
			ok:  true,
			out: &ast.StringLiteral{
				Orig:   `"asd\nqwe"`,
				Parsed: "asd\nqwe",
			},
		},
	}

	for k, v := range tcases {
		t.Run(k, fn(v))
	}

}

func TestParseLiteral(t *testing.T) {
	type tcase struct {
		str string
		ok  bool
		err error
		out interface{}
	}

	fn := func(tc tcase) func(t *testing.T) {
		return func(t *testing.T) {
			_, v, ok, err :=
				parser.DoParseStringForTest(
					&ast.LiteralParser{}, tc.str, "test")

			assertErrIs(t, tc.err, err)
			assertEq(t, tc.ok, ok)
			if !ok || err != nil {
				return
			}

			assertEq(t, tc.out, v)
		}
	}

	tcases := map[string]tcase{
		"simple": tcase{
			str: `"asd"`,
			ok:  true,
			out: &ast.StringLiteral{
				Orig:   `"asd"`,
				Parsed: "asd",
			},
		},
		"newline": tcase{
			str: `"asd\nqwe"`,
			ok:  true,
			out: &ast.StringLiteral{
				Orig:   `"asd\nqwe"`,
				Parsed: "asd\nqwe",
			},
		},
		"pos": tcase{
			str: `123`,
			ok:  true,
			out: &ast.NumberLiteral{
				Orig:   `123`,
				Parsed: 123,
			},
		},
	}

	for k, v := range tcases {
		t.Run(k, fn(v))
	}

}
