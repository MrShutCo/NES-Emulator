package cpu

func Flag() {
	newInst(0x18, "CLC", "", 2)
	newInst(0xD8, "CLD", "", 2)
	newInst(0x58, "CLI", "", 2)
	newInst(0xB8, "CLV", "", 2)
	newInst(0x78, "SEI", "", 2)
	newInst(0x38, "SEC", "", 2)
	newInst(0xF8, "SED", "", 2)
	// CLC
	FuncMap[0x18] = func() {
		setCarryFlag(false)
		PC++
	}
	// CLD
	FuncMap[0xD8] = func() {
		setDecimalFlag(false)
		PC++
	}
	// CLI
	FuncMap[0x58] = func() {
		setInterruptFlag(false)
		PC++
	}
	// CLV
	FuncMap[0xB8] = func() {
		setOverflowFlag(false)
		PC++
	}

	// SEI
	FuncMap[0x78] = func() {
		setInterruptFlag(true)
		PC++
	}
	// SEC
	FuncMap[0x38] = func() {
		setCarryFlag(true)
		PC++
	}
	// SED
	FuncMap[0xF8] = func() {
		setDecimalFlag(true)
		PC++
	}
}
