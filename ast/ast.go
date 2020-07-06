package ast

import "github.com/ear7h/lang/ast/parser"

type Node interface {
	FileInfo() parser.FileInfo
	setFi(parser.FileInfo)
}

type BaseNode struct {
	Fi parser.FileInfo
}

type fileInfoSetter struct {
	//setFiCursor(c *parser.Cursor)
	setFi(parser.FileInfo)
}

func (bn *BaseNode) FileInfo() parser.FileInfo {
	return bn.Fi
}

func (bn *BaseNode) setFi(fi parser.FileInfo) {
	bn.Fi = fi
}

func (bn *BaseNode) setFileInfo(c *parser.Cursor) {
	bn.Fi = c.FileInfo()
}

/*
func (bn *BaseNode) Parse(c *Cursor) {
	bn.Fi = c.FileInfo()
}
*/
