package cpu

import "fmt"

func JMP() {
	newInst(0x4C, "JMP", "absolute", 3)
	newInst(0x4C, "JMP", "indirect", 5)
	newInst(0x20, "JSR", "", 6)
	newInst(0x60, "RTS", "", 6)
	// JMP absolute
	FuncMap[0x4C] = func() {
		PC = bytesToInt16(RAM[PC+2], RAM[PC+1])
		output = fmt.Sprintf("$%04X", PC)
	}
	// JMP indirect
	// TODO: finish
	FuncMap[0x6C] = func() { // JMP indirect
		//PC = RAM[bytesToInt16(RAM[PC+2], RAM[PC+1])]
	}
	// JSR
	FuncMap[0x20] = func() {
		push(highByte(PC))
		push(lowByte(PC))
		PC = bytesToInt16(RAM[PC+2], RAM[PC+1])
		output = fmt.Sprintf("$%04X", PC)
	}
	// RTS
	FuncMap[0x60] = func() {
		//printStack()
		l := pull()
		h := pull()
		PC = bytesToInt16(h, l)
		PC += 3
	}
}

func Compare() {
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

func Other() {
	newInst(0x00, "BRK", "implied", 7)
	newInst(0x40, "RTI", "implied", 6)
	newInst(0x24, "BIT", "zeropage", 3)
	newInst(0x2C, "BIT", "absolute", 4)
	newInst(0xEA, "NOP", "implied", 2)
	// BRK
	FuncMap[0x00] = func() {
		retAddr := PC + 2

		push(highByte(retAddr))
		push(lowByte(retAddr))
		push(SR)
		setBreakFlag(true)
		PC = GetWordAt(NMI_VECTOR)
	}

	// RTI
	FuncMap[0x40] = func() {
		SR = pull()
		PC_LOW := pull()
		PC_HIGH := pull()
		PC = bytesToInt16(PC_HIGH, PC_LOW)
	}

	// BIT zeropage
	FuncMap[0x24] = func() {
		val := RAM[RAM[PC+1]]
		setNegativeFlag(getBit(val, 7))
		setOverflowFlag(getBit(val, 6))
		//SR = SR ^ (RAM[RAM[PC+1]] & 0b1100_0000) // Transfer bits 6 and 7
		setZeroFlag(AC&val == 0)
		output = fmt.Sprintf("$%02X = %02X", RAM[PC+1], AC)
		PC += 2
	}
	// BIT absolute
	FuncMap[0x2C] = func() {
		mem := getNextWord()
		SR = SR ^ (RAM[mem] & 0b1100_0000) // Transfer bits 6 and 7
		setZeroFlag(AC&RAM[PC+1] == 0)
		PC += 3
	}

	// NOP
	FuncMap[0xEA] = func() {
		PC++
	}
}
