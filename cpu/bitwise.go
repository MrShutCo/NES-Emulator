package cpu

type foo struct {
	opcode byte
	f      func()
}

func wrap(f func()) func() {
	return func() { f() }
}

func Bitwise() {
	newInst(0x29, "AND", "immediate", 2)
	newInst(0x25, "AND", "zeropage", 3)
	newInst(0x35, "AND", "zeropage,X", 4)
	newInst(0x2D, "AND", "absolute", 4)
	newInst(0x3D, "AND", "absolute,X", 4)
	newInst(0x39, "AND", "absolute,Y", 4)
	a := []foo{
		{0x29, func() { and(immed) }},
		{0x25, func() { and(zeropage) }},
		{0x35, func() { and(zeropageX) }},
		{0x2D, func() { and(absolute) }},
		{0x3D, func() { and(absoluteX) }},
		{0x39, func() { and(absoluteY) }},
	}
	apply(a)

	newInst(0x09, "ORA", "immediate", 2)
	newInst(0x05, "ORA", "zeropage", 3)
	newInst(0x15, "ORA", "zeropage,X", 4)
	newInst(0x0D, "ORA", "absolute", 4)
	newInst(0x1D, "ORA", "absolute,X", 4)
	newInst(0x19, "ORA", "absolute,Y", 4)
	e := []foo{
		{0x09, func() { ora(immed) }},
		{0x05, func() { ora(zeropage) }},
		{0x15, func() { ora(zeropageX) }},
		{0x0D, func() { ora(absolute) }},
		{0x1D, func() { ora(absoluteX) }},
		{0x19, func() { ora(absoluteY) }},
	}
	apply(e)

	newInst(0x49, "EOR", "immediate", 2)
	newInst(0x45, "EOR", "zeropage", 3)
	newInst(0x55, "EOR", "zeropage,X", 4)
	newInst(0x4D, "EOR", "absolute", 4)
	newInst(0x5D, "EOR", "absolute,X", 4)
	newInst(0x59, "EOR", "absolute,Y", 4)
	eo := []foo{
		{0x49, func() { eor(immed) }},
		{0x45, func() { eor(zeropage) }},
		{0x55, func() { eor(zeropageX) }},
		{0x4D, func() { eor(absolute) }},
		{0x5D, func() { eor(absoluteX) }},
		{0x59, func() { eor(absoluteY) }},
	}
	apply(eo)
}

func and(f func() byte) {
	SetAC(AC & f())
}

func ora(f func() byte) {
	SetAC(AC | f())
}

func eor(f func() byte) {
	SetAC(AC ^ f())
}
