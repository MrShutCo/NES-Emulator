package cpu

func Bitwise() {
	a := map[byte]func() byte{
		0x29: immed,
		0x25: zeropage,
		0x35: zeropageX,
		0x2D: absolute,
		0x3D: absoluteX,
		0x39: absoluteY,
	}
	apply(and, a)

	e := map[byte]func() byte{
		0x49: immed,
		0x45: zeropage,
		0x55: zeropageX,
		0x4D: absolute,
		0x5D: absoluteX,
		0x59: absoluteY,
	}
	apply(eor, e)
}

func and(f func() byte) {
	SetAC(AC & f())
}

func eor(f func() byte) {
	SetAC(AC | f())
}
