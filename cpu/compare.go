package cpu

func Compare() {
	newInst(0xC9, "CMP", "immediate", 2)
	newInst(0xC5, "CMP", "zeropage", 3)
	newInst(0xD5, "CMP", "zeropage,X", 4)
	newInst(0xCD, "CMP", "absolute", 4)
	newInst(0xDD, "CMP", "absolute,X", 4)
	newInst(0xD9, "CMP", "absolute,Y", 4)
	newInst(0xC1, "CMP", "(indirect,X)", 6)
	newInst(0xD1, "CMP", "(indirect),Y", 5)
	cmca := []foo{
		{0xC9, func() { cmp(AC, immed()) }},
		{0xC5, func() { cmp(AC, zeropage()) }},
		{0xD5, func() { cmp(AC, zeropageX()) }},
		{0xCD, func() { cmp(AC, absolute()) }},
		{0xDD, func() { cmp(AC, absoluteX()) }},
		{0xD9, func() { cmp(AC, absoluteY()) }},
		{0xC1, func() { cmp(AC, indirectX()) }},
		{0xD1, func() { cmp(AC, indirectY()) }},
	}
	apply(cmca)

	newInst(0xE0, "CPX", "immediate", 2)
	newInst(0xE4, "CPX", "zeropage", 3)
	newInst(0xEC, "CPX", "absolute", 4)
	cmx := []foo{
		{0xE0, func() { cmp(X, immed()) }},
		{0xE4, func() { cmp(X, zeropage()) }},
		{0xEC, func() { cmp(X, absolute()) }},
	}
	apply(cmx)

	newInst(0xC0, "CPY", "immediate", 2)
	newInst(0xC4, "CPY", "zeropage", 3)
	newInst(0xCC, "CPY", "absolute", 4)
	cmy := []foo{
		{0xC0, func() { cmp(Y, immed()) }},
		{0xC4, func() { cmp(Y, zeropage()) }},
		{0xCC, func() { cmp(Y, absolute()) }},
	}
	apply(cmy)
}

func cmp(value byte, memory byte) {
	result := value - memory
	if value < memory {
		setCarryFlag(false)
		setZeroFlag(false)
		setNegativeFlag(result >= 1<<7)
	} else if value == memory {
		setCarryFlag(true)
		setZeroFlag(true)
		setNegativeFlag(false)
	} else {
		setCarryFlag(true)
		setZeroFlag(false)
		setNegativeFlag(result >= 1<<7)
	}
}
