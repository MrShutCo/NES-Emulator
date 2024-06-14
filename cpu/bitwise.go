package cpu

type foo struct {
	opcode byte
	f      func()
}

func Bitwise() {
	newInst(0x29, "AND", "immediate", 2)
	newInst(0x25, "AND", "zeropage", 3)
	newInst(0x35, "AND", "zeropage,X", 4)
	newInst(0x2D, "AND", "absolute", 4)
	newInst(0x3D, "AND", "absolute,X", 4)
	newInst(0x39, "AND", "absolute,Y", 4)
	newInst(0x21, "AND", "(indirect,X)", 6)
	newInst(0x31, "AND", "(indirect),Y", 5)
	a := []foo{
		{0x29, func() { and(immed) }},
		{0x25, func() { and(zeropage) }},
		{0x35, func() { and(zeropageX) }},
		{0x2D, func() { and(absolute) }},
		{0x3D, func() { and(absoluteX) }},
		{0x39, func() { and(absoluteY) }},
		{0x21, func() { and(indirectX) }},
		{0x31, func() { and(indirectY) }},
	}
	apply(a)

	newInst(0x09, "ORA", "immediate", 2)
	newInst(0x05, "ORA", "zeropage", 3)
	newInst(0x15, "ORA", "zeropage,X", 4)
	newInst(0x0D, "ORA", "absolute", 4)
	newInst(0x1D, "ORA", "absolute,X", 4)
	newInst(0x19, "ORA", "absolute,Y", 4)
	newInst(0x01, "ORA", "indirect,X", 6)
	newInst(0x11, "ORA", "(indirect),Y", 5)
	e := []foo{
		{0x09, func() { ora(immed) }},
		{0x05, func() { ora(zeropage) }},
		{0x15, func() { ora(zeropageX) }},
		{0x0D, func() { ora(absolute) }},
		{0x1D, func() { ora(absoluteX) }},
		{0x19, func() { ora(absoluteY) }},
		{0x01, func() { ora(indirectX) }},
		{0x11, func() { ora(indirectY) }},
	}
	apply(e)

	newInst(0x49, "EOR", "immediate", 2)
	newInst(0x45, "EOR", "zeropage", 3)
	newInst(0x55, "EOR", "zeropage,X", 4)
	newInst(0x4D, "EOR", "absolute", 4)
	newInst(0x5D, "EOR", "absolute,X", 4)
	newInst(0x59, "EOR", "absolute,Y", 4)
	newInst(0x41, "EOR", "indirect,X", 6)
	newInst(0x51, "EOR", "indirect,Y", 5)
	eo := []foo{
		{0x49, func() { eor(immed) }},
		{0x45, func() { eor(zeropage) }},
		{0x55, func() { eor(zeropageX) }},
		{0x4D, func() { eor(absolute) }},
		{0x5D, func() { eor(absoluteX) }},
		{0x59, func() { eor(absoluteY) }},
		{0x41, func() { eor(indirectX) }},
		{0x51, func() { eor(indirectY) }},
	}
	apply(eo)

	newInst(0x4A, "LSR", "immediate", 2)
	newInst(0x46, "LSR", "zeropage", 5)
	newInst(0x56, "LSR", "zeropage,X", 6)
	newInst(0x4E, "LSR", "absolute", 6)
	newInst(0x5E, "LSR", "absolute,X", 7)
	ls := []foo{
		{0x4A, func() {
			output := lsr(ac())
			SetAC(output)
		}},
		{0x46, func() {
			SetRAM(zeropageAddr(), lsr(zeropage()))
		}},
		{0x56, func() {
			SetRAM(zeropageXAddr(), lsr(zeropageX()))
		}},
		{0x4E, func() {
			SetRAM(absoluteAddr(), lsr(absolute()))
		}},
		{0x5E, func() {
			SetRAM(absoluteXAddr(), lsr(absoluteX()))
		}},
	}
	apply(ls)

	newInst(0x0A, "ASL", "immediate", 2)
	newInst(0x06, "ASL", "zeropage", 5)
	newInst(0x16, "ASL", "zeropage,X", 6)
	newInst(0x0E, "ASL", "absolute", 6)
	newInst(0x1E, "ASL", "absolute,X", 7)
	as := []foo{
		{0x0A, func() { SetAC(asl(ac())) }},
		{0x06, func() { RAM[zeropageAddr()] = asl(zeropage()) }},
		{0x16, func() { RAM[zeropageXAddr()] = asl(zeropageX()) }},
		{0x0E, func() { RAM[absoluteAddr()] = asl(absolute()) }},
		{0x1E, func() { RAM[absoluteXAddr()] = asl(absoluteX()) }},
	}
	apply(as)

	newInst(0x6A, "ROR", "immediate", 2)
	newInst(0x66, "ROR", "zeropage", 5)
	newInst(0x76, "ROR", "zeropage,X", 6)
	newInst(0x6E, "ROR", "absolute", 6)
	newInst(0x7E, "ROR", "absolute,X", 7)
	ro := []foo{
		{0x6A, func() { output := ror(ac()); SetAC(output) }},
		{0x66, func() { RAM[zeropageAddr()] = ror(zeropage()) }},
		{0x76, func() { RAM[zeropageXAddr()] = ror(zeropageX()) }},
		{0x6E, func() { RAM[absoluteAddr()] = ror(absolute()) }},
		{0x7E, func() { RAM[absoluteXAddr()] = ror(absoluteX()) }},
	}
	apply(ro)

	newInst(0x2A, "ROL", "immediate", 2)
	newInst(0x26, "ROL", "zeropage", 5)
	newInst(0x36, "ROL", "zeropage,X", 6)
	newInst(0x2E, "ROL", "absolute", 6)
	newInst(0x3E, "ROL", "absolute,X", 7)
	lo := []foo{
		{0x2A, func() { SetAC(rol(ac())) }},
		{0x26, func() { RAM[zeropageAddr()] = rol(zeropage()) }},
		{0x36, func() { RAM[zeropageXAddr()] = rol(zeropageX()) }},
		{0x2E, func() { RAM[absoluteAddr()] = rol(absolute()) }},
		{0x3E, func() { RAM[absoluteXAddr()] = rol(absoluteX()) }},
	}
	apply(lo)
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

func lsr(val byte) byte {
	setCarryFlag(getBit(val, 0))
	val = val >> 1
	setZeroFlag(val == 0x0)
	setNegativeFlag(val >= 0x80)
	return val
}

func asl(val byte) byte {
	setCarryFlag(getBit(val, 7)) // bit 7 to carry
	val = (val << 1) & 0b1111_1110
	setZeroFlag(val == 0x0)
	setNegativeFlag(val >= 0x80)
	return val
}

func ror(val byte) byte {
	orginalVal := val
	val = val >> 1
	if isCarrySet() { // set bit 7 to 1
		val = val | 0x80
	} else { // set bit 7 to 0
		val = val & 0b0111_1111
	}
	setCarryFlag(getBit(orginalVal, 0)) // bit 0 to carry
	setZeroFlag(val == 0x0)
	setNegativeFlag(val >= 0x80)
	return val
}

func rol(val byte) byte {
	orginalVal := val
	val = val << 1
	if isCarrySet() { // set bit 0 to 1
		val = val | 0x01
	} else { // set bit 0 to 0
		val = val & 0b1111_1110
	}
	setCarryFlag(getBit(orginalVal, 7)) // bit 7 to carry
	setZeroFlag(val == 0x0)
	setNegativeFlag(val >= 0x80)
	return val
}
