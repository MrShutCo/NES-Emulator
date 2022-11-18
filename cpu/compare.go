package cpu

func Compare() {
	newInst(0xC9, "CMP", "immediate", 2)
	cmpac := []foo{
		{0xC9, func() { cmp(AC, immed()) }},
	}
	apply(cmpac)

	// CPX immeadiate
	FuncMap[0xE0] = func() {
		cmp(X, RAM[PC+1])
		PC += 2
	}
	// CPX zeropage
	FuncMap[0xE4] = func() {
		cmp(X, RAM[RAM[PC+1]])
		PC += 2
	}
	// CPX absolute
	FuncMap[0xEC] = func() {
		cmp(X, RAM[getNextWord()])
		PC += 3
	}

	// CPY immeadiate
	FuncMap[0xC0] = func() {
		cmp(Y, RAM[PC+1])
		PC += 2
	}
	// CPY zeropage
	FuncMap[0xC4] = func() {
		cmp(Y, RAM[RAM[PC+1]])
		PC += 2
	}
	// CPY absolute
	FuncMap[0xCC] = func() {
		cmp(Y, RAM[getNextWord()])
		PC += 3
	}
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
