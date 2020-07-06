package ast

import (
	"fmt"

	"github.com/ear7h/lang/ast/parser"
)

var _ = fmt.Println

type ExprParser struct{}

func (ExprParser) Parse(c *parser.Cursor) (interface{}, bool) {
	return parser.First(
		&UnaryExpr{},
		&BinaryExpr{},
	).Parse(c)
}

type ExprOperandParser struct{}

func (ExprOperandParser) Parse(c *parser.Cursor) (interface{}, bool) {
	fi := c.FileInfo()

	v, ok :=  parser.All(
		parser.First(
			LiteralParser{},
			&Ident{},
			parser.Braced(
				parser.ExpectString("("),
				ExprParser{},
				parser.ExpectString(")"),
			),
		),
		parser.Maybe(ExprOperandParser1{}),
	).Parse(c)
	if !ok {
		return nil, false
	}

	slc := v.([]interface{})


	left := slc[0]

	if slc[1] == nil {
		return left, true
	}

	n := slc[1].(*ObjExpr)

	n.Object = left
	n.setFi(fi)

	return n, true
}

// remove left recursion
type ExprOperandParser1 struct{}

func (ExprOperandParser1) Parse(c *parser.Cursor) (interface{}, bool) {

	v, ok := parser.First(
		ObjExprRightParser{},
		/*
		parser.First(
			parser.Kleene(
				parser.All(
					ExpectString("."),
					&Ident{},
				),
			),
		),
		*/
	).Parse(c)
	if !ok {
		return nil, false
	}


	vv, ok := ExprOperandParser1{}.Parse(c)
	if !ok {
		return v,  true
	}

	objExpr := v.(*ObjExpr)
	objExpr.Right = vv

	return objExpr, true
}

type LeftRecursive interface {
	SetRight(v interface{})
	SetLeft(v interface{})
}


