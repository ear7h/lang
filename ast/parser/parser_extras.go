package parser

var Test = false
var Debug = false

func (c *Cursor) PeekN(n int) string {
	if !(Debug || Test) {
		panic("!Debug")
	}

	cc := *c

	ret := ""
	for ;uint(n) > 0 && !cc.EOF(); n--{
		ret += string(cc.ReadRune())
	}

	return ret
}
