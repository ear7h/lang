package parser

import (
	"unicode"
)

type Parser interface {
	Parse(*Cursor) (ret interface{}, ok bool)
}

type ParserFunc func(*Cursor) (ret interface{}, ok bool)

func (fn ParserFunc) Parse(c *Cursor) (ret interface{}, ok bool) {
	return fn(c)
}


func ReadRune() Parser {
	return ParserFunc(func(c *Cursor) (interface{}, bool) {
		return c.readRune(), true
	})
}

func ExpectRune(expect rune) Parser {
	return ParserFunc(func(c *Cursor) (interface{}, bool) {
		r := c.readRune()
		if r != expect {
			return nil, false
		}

		return r, true
	})
}

func First(p ...Parser) Parser {
	return ParserFunc(func(c *Cursor) (interface{}, bool) {
		for _, v := range p {
			cc := *c
			ret, ok := v.Parse(&cc)
			if ok {
				*c = cc
				return ret, true
			}
		}

		return nil, false
	})
}

func FirstString(slc ...string) Parser {
	return ParserFunc(func(c *Cursor) (interface{}, bool) {
		for _, v := range slc {
			cc := *c
			_, ok := ExpectString(v).Parse(&cc)
			if ok {
				*c = cc
				return v, true
			}
		}

		return nil, false
	})
}

func All(p ...Parser) Parser {
	return ParserFunc(func(c *Cursor) (interface{}, bool) {
		ret := make([]interface{}, len(p))

		for i, v := range p {
			cc := *c
			var ok bool
			ret[i], ok = v.Parse(&cc)
			if !ok {
				return nil, false
			}

			*c = cc
		}

		return ret, true
	})
}

func AllIdx(idx int, p ...Parser) ParserFunc {
	parser := All(p...)

	return func(c *Cursor) (interface{}, bool) {
		v, ok := parser.Parse(c)
		if !ok {
			return nil, false
		}

		return v.([]interface{})[idx], true
	}
}

func Maybe(p Parser) Parser {
	return ParserFunc(func(c *Cursor) (interface{}, bool) {
		cc := *c
		v, ok := p.Parse(c)
		if !ok {
			return nil, true
		}

		*c = cc
		return v, true
	})
}

// WriteTo returns a parser that wraps another parser p
// and sets the literal string tha p matched to *dst
func WriteTo(dst *string, p Parser) Parser {
	return ParserFunc(func(c *Cursor) (interface{}, bool) {
		start := c.i

		ret, ok := p.Parse(c)
		if !ok {
			return nil, false
		}
		end := c.i
		if c.eof {
			end--
		}

		*dst = c.stringAt(start, end-start)

		return ret, true
	})
}


func KleenePred(fn func(r rune) bool) ParserFunc {
	return func(c *Cursor) (interface{}, bool) {
		ret := ""
		for !c.EOF() && fn(c.PeekRune()) {
			ret += string(c.ReadRune())
		}

		return ret, true
	}
}

func PlusPred(fn func(r rune) bool) ParserFunc {
	return func(c *Cursor) (interface{}, bool) {
		v, ok := KleenePred(fn).Parse(c)
		if !ok {
			return nil, false
		}

		if len(v.(string)) == 0 {
			return nil, false
		}

		return v, true
	}
}

// WS returns a parser that matches white space
func WS() Parser {
	return KleenePred(unicode.IsSpace)
}

func WS1() Parser {
	return PlusPred(unicode.IsSpace)
}

// HS returns a parser that matches horizontal space
func HS() Parser {
	return KleenePred(func(r rune) bool {
		return unicode.IsSpace(r) && r != '\r' && r != '\n'
	})
}

func HS1() Parser {
	return PlusPred(func(r rune) bool {
		return unicode.IsSpace(r) && r != '\r' && r != '\n'
	})
}

// EOL scanns until the end of the line
func EOL() Parser {
	atEnd := false
	return KleenePred(func(r rune) bool {
		if r == '\r' || r == '\n' {
			atEnd = atEnd || true
			return true
		} else if atEnd {
			return false
		} else {
			return true
		}
	})
}

func Braced(left, middle, right Parser) Parser {
	return AllIdx(2,
		left,
		WS(),
		middle,
		WS(),
		right,
	)
}
