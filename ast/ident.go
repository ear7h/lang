package ast

import (
	"fmt"
	"unicode"

	"github.com/ear7h/lang/ast/parser"
)

type Ident struct {
	BaseNode

	IsExported bool
	Name       string
}

func (n *Ident) Parse(c *parser.Cursor) (interface{}, bool) {
	n.setFileInfo(c)

	r := c.PeekRune()


	if unicode.In(r, unicode.Ll) || r == '_' {
		n.IsExported = false
	} else if unicode.In(r, unicode.Lu) {
		n.IsExported = true
	} else {
		return nil, false
	}

	str := string(c.ReadRune())

	isIdentTail := func(r rune) bool {
		return unicode.In(r, unicode.Ll, unicode.Lu) ||
			unicode.IsDigit(r) ||
			r == '_'
	}

	for isIdentTail(c.PeekRune()) {
		str += string(c.ReadRune())
	}

	n.Name = str

	return n, true
}

var _ = fmt.Println

const (
	ObjField = iota
	// TODO
	ObjIdx
	ObjCall
)

type ObjExpr struct {
	BaseNode
	Object interface{} // root object
	Op int
	Arg interface{}
	Right interface{}
}

type ObjExprRightParser struct {}


func (_ ObjExprRightParser) Parse(c *parser.Cursor) (interface{}, bool) {
	return parser.First(
		ObjFieldRightParser{},
	).Parse(c)
}

type ObjFieldRightParser struct {}


func (_ ObjFieldRightParser) Parse(c *parser.Cursor) (interface{}, bool) {
	var (
		n ObjExpr
		ok bool
	)

	n.setFileInfo(c)

	_, ok = parser.ExpectString(".").Parse(c)
	if !ok {
		return nil, false
	}

	n.Arg, ok = (&Ident{}).Parse(c)
	if !ok {
		return nil, false
	}

	return &n, true
}
