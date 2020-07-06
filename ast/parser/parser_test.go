package parser_test

import (
	"testing"
	"github.com/ear7h/lang/ast/parser"
)

func TestParsers(t *testing.T) {
	type tcase struct {
		str string
		p   parser.Parser
		ok  bool
		err error
		out interface{}
	}

	fn := func(tc tcase) func(t *testing.T) {
		return func(t *testing.T) {
			_, v, ok, err :=
				parser.DoParseStringForTest(tc.p, tc.str, "test")

			assertErrIs(t, tc.err, err)
			assertEq(t, tc.ok, ok)
			if !ok || err != nil {
				return
			}

			assertEq(t, tc.out, v)
		}
	}

	tcases := map[string]tcase{
		"ExpectString": tcase{
			str: `asd`,
			ok:  true,
			p:   parser.ExpectString("asd"),
			out: "asd",
		},
		"ExpectString1": tcase{
			str: `asdf`,
			ok:  true,
			p:   parser.ExpectString("asd"),
			out: "asd",
		},
		"FirstString": tcase{
			str: "asd",
			ok: true,
			p: parser.FirstString("qwe", "asd"),
			out: "asd",
		},
		"FirstString1": tcase{
			str: "asdf",
			ok: true,
			p: parser.FirstString("qwe", "asd"),
			out: "asd",
		},
		"First": tcase{
			str: "asd",
			ok: true,
			p: parser.First(
				parser.ExpectString("qwe"),
				parser.ExpectString("asd"),
			),
			out: "asd",
		},
		"All": tcase{
			str: "qweasd",
			ok: true,
			p: parser.All(
				parser.ExpectString("qwe"),
				parser.ExpectString("asd"),
			),
			out: []interface{}{
				"qwe",
				"asd",
			},
		},
	}

	for k, v := range tcases {
		t.Run(k, fn(v))
	}

}
