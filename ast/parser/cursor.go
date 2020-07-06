package parser

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type FileInfo struct {
	Name string
	Line int64
	Col  int64
}

func NewCursorString(s string, name string) *Cursor {
	return &Cursor{
		r:    strings.NewReader(s),
		i:    0,
		name: name,
		line: 1,
		col:  1,
	}
}

type Cursor struct {
	r    io.ReaderAt
	i    int64
	eof  bool
	name string
	line int64
	col  int64
}

func (c *Cursor) Fatalf(f string, v ...interface{}) {
	panic(fmt.Errorf(f, v...))
}

func (c Cursor) Fatal(v interface{}) {
	panic(v)
}

// reads the next rune
func (c *Cursor) readRune() rune {

	buf := make([]byte, 4)
	n, err := c.r.ReadAt(buf, c.i)
	if err != nil {
		if errors.Is(err, io.EOF) && n == 0 && c.eof {
			// reading past eof
			c.Fatal(err)
		} else if errors.Is(err, io.EOF) && n == 0 {
			// got to eof
			c.eof = true
		} else if !errors.Is(err, io.EOF) {
			// some other error
			c.Fatal(err)
		}
	}

	r, n := utf8.DecodeRune(buf)
	if r == utf8.RuneError {
		c.Fatal("rune error")
	}

	c.i += int64(n)
	if r == '\n' {
		c.line++
		c.col = 0
	} else {
		c.col++
	}

	return r
}

// reads an n length string starting at off
// this method mutates the internal cursor, so use wisely
func (c *Cursor) stringAt(off, n int64) string {
	c.eof = false
	c.i = off
	ret := ""
	for ; n > 0; n-- {
		ret += string(c.readRune())
	}
	return ret
}

func (c *Cursor) ReadRune() rune {
	return c.readRune()
}

func (c *Cursor) PeekRune() rune {
	cc := *c
	return cc.readRune()
}

func (c *Cursor) EOF() bool {
	return c.eof
}

func (c *Cursor) FileInfo() FileInfo {
	return FileInfo{
		Name: c.name,
		Line: c.line,
		Col:  c.col,
	}
}
