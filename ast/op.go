package ast

import (
	"fmt"

	"github.com/ear7h/lang/ast/parser"
)

var _ = fmt.Println

const (
	UnaryPos   rune = '+'
	UnaryNeg   rune = '-'
	UnaryDeref rune = '*'
	UnaryAddr  rune = '&'
	UnaryNot   rune = '!'
)

var unaryOperators = []rune{
	UnaryPos,
	UnaryNeg,
	UnaryDeref,
	UnaryAddr,
	UnaryNot,
}

type UnaryExpr struct {
	BaseNode
	Op      rune
	Operand interface{}
}

func (n *UnaryExpr) Parse(c *parser.Cursor) (interface{}, bool) {
	n.setFileInfo(c)

	r := c.PeekRune()

	found := false
	for _, v := range unaryOperators {
		if r == v {
			found = true
			n.Op = c.ReadRune()
			break
		}
	}

	if !found {
		return nil, false
	}

	var ok bool
	n.Operand, ok = ExprParser{}.Parse(c)
	if !ok {
		return nil, false
	}

	return n, true
}

const (
	// arith
	BinaryAdd = "+"
	BinarySub = "-"
	BinaryMul = "*"
	BinaryDiv = "/"
	BinaryMod = "%"

	// bits
	BinaryShr    = ">>"
	BinaryShl    = "<<"
	BinaryBitAnd = "&"
	BinaryBitOr  = "|"
	BinaryBitXor = "^"

	// bool
	BinaryBoolAnd = "&&"
	BinaryBoolOr  = "||"

	// cmp
	BinaryLt  = "<"
	BinaryGt  = ">"
	BinaryLte = "<="
	BinaryGte = ">="
	BinaryEq  = "=="
	BinaryNeq = "!="
)

// BinaryExpr represents an expression with a binary operator
// and it's two operands.
// Precedence:
//		<< >> & | ^
//		* / %
//		+ -
//		< > <= >= == !=
//		&& ||
// associativity is left to right
type BinaryExpr struct {
	BaseNode
	Op          string
	Left, Right interface{}
}

func (_ *BinaryExpr) Parse(c *parser.Cursor) (interface{}, bool) {

	precs := [][]string{
		{
			BinaryShr, BinaryShl, BinaryBitAnd, BinaryBitOr, BinaryBitXor,
		},
		{
			BinaryMul, BinaryDiv, BinaryMod,
		},
		{
			BinaryAdd, BinarySub,
		},
		{
			BinaryLt, BinaryGt, BinaryLte, BinaryGte, BinaryEq, BinaryNeq,
		},
		{
			BinaryBoolAnd, BinaryBoolOr,
		},
	}

	lower := func() parser.Parser {
		return ExprOperandParser{}
	}

	for _, v := range precs {
		vv := v
		tmp := lower
		lower = func() parser.Parser {
			return parser.First(
				BinaryExprPrecedenceGroup(tmp, vv...),
				tmp(),
			)
		}
	}

	return lower().Parse(c)
}

func BinaryExprPrecedenceGroup(lower func() parser.Parser,
	ops ...string) parser.Parser {
	return parser.ParserFunc(func(c *parser.Cursor) (interface{}, bool) {
		var (
			n  BinaryExpr
			ok bool
		)

		n.setFileInfo(c)

		v, ok := parser.All(
			lower(),
			parser.WS(),
			parser.FirstString(ops...),
			parser.WS(),
			lower(),
		).Parse(c)
		if !ok {
			return nil, false
		}

		slc := v.([]interface{})

		n.Left = slc[0]
		n.Op = slc[2].(string)
		n.Right = slc[4]

		return &n, true
	})
}
