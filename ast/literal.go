package ast

import (
	"strconv"
	"unicode"

	"github.com/ear7h/lang/ast/parser"
)

type LiteralParser struct{}

func (_ LiteralParser) Parse(c *parser.Cursor) (interface{}, bool) {
	return parser.First(&StringLiteral{}, &NumberLiteral{}).Parse(c)
}

type StringLiteral struct {
	BaseNode

	Orig   string
	Parsed string
}

func (n *StringLiteral) Parse(c *parser.Cursor) (interface{}, bool) {
	n.setFileInfo(c)

	var orig string

	parsed, ok := parser.WriteTo(&orig,
		parser.ParserFunc(func(c *parser.Cursor) (interface{}, bool) {
			buf := ""
			if c.ReadRune() != '"' {
				return nil, false
			}

			r := c.ReadRune()

			for ; r != '"'; r = c.ReadRune() {
				switch r {
				case '\\':
					r = c.ReadRune()
					switch r {
					case '\\':
						r = '\\'
					case 'n':
						r = '\n'
					case 't':
						r = '\t'
					default:
						panic("unknwown escape sequece \\" + string(r))
					}
				}

				buf += string(r)
			}
			if r != '"' {
				return nil, false
			}

			return buf, true
		})).Parse(c)

	if !ok {
		return nil, false
	}

	n.Orig = string(orig)
	n.Parsed = string(parsed.(string))

	return n, true
}

type NumberLiteral struct {
	BaseNode

	Orig   string
	Parsed int64
}

func (n *NumberLiteral) Parse(c *parser.Cursor) (interface{}, bool) {
	n.setFileInfo(c)

	var orig string

	parsed, ok := parser.WriteTo(&orig,
		parser.ParserFunc(func(c *parser.Cursor) (interface{}, bool) {
			buf := ""

			for unicode.IsDigit(c.PeekRune()) {
				buf += string(c.ReadRune())
				if c.EOF() {
					break
				}
			}

			parsed, err := strconv.ParseInt(string(buf), 10, 64)
			if err != nil {
				// c.Fatal(err)
				return nil, false
			}

			return parsed, true
		})).Parse(c)

	if !ok {
		return nil, false
	}

	n.Orig = string(orig)
	n.Parsed = parsed.(int64)

	return n, true
}
