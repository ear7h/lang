package parser

import "fmt"

func ExpectString(s string) Parser {
	return ParserFunc(func(c *Cursor) (interface{}, bool) {
		for _, v := range s {
			if v != c.ReadRune() {
				return nil, false
			}
		}

		return s, true
	})
}

func MustParseString(p Parser, s string) interface{} {
	v, ok, err := DoParseString(p, s, "MustParseString")
	if !ok {
		panic("parse !ok")
	}
	if err != nil {
		panic(err)
	}

	return v
}

func DoParseString(p Parser, s string, name string) (v interface{}, ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			if rerr, ok := r.(error); ok {
				err = rerr
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	v, ok = p.Parse(NewCursorString(s, name))
	return v, ok, err
}

func DoParseStringForTest(p Parser, s string, name string) (initCur Cursor, v interface{}, ok bool, err error) {
	/*
	defer func() {
		if r := recover(); r != nil {
			if rerr, ok := r.(error); ok {
				err = rerr
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	*/

	c := NewCursorString(s, name)
	initCur = *c

	v, ok = p.Parse(c)
	return initCur, v, ok, err
}
