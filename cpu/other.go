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
		push(highByte(PC + 2))
		push(lowByte(PC + 2))
		PC = bytesToInt16(RAM[PC+2], RAM[PC+1])
		output = fmt.Sprintf("$%04X", PC)
	}
	// RTS
	FuncMap[0x60] = func() {
		//printStack()
		l := pull()
		h := pull()
		PC = bytesToInt16(h, l)
		PC += 1
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
		oldSR := SR
		SR = pull()
		setEmptyFlag(getBit(oldSR, 5))
		PC_LOW := pull()
		PC_HIGH := pull()
		PC = bytesToInt16(PC_HIGH, PC_LOW)
	}

	// BIT zeropage
	FuncMap[0x24] = func() {
		val := RAM[RAM[PC+1]]
		setNegativeFlag(getBit(val, 7))
		setOverflowFlag(getBit(val, 6))
		setZeroFlag(AC&val == 0)
		output = fmt.Sprintf("$%02X = %02X", RAM[PC+1], val)
		PC += 2
	}
	// BIT absolute
	FuncMap[0x2C] = func() {
		val := RAM[getNextWord()]
		setNegativeFlag(getBit(val, 7))
		setOverflowFlag(getBit(val, 6))
		setZeroFlag(AC&val == 0)
		output = fmt.Sprintf("$%04X = %02X", getNextWord(), val)
		PC += 3
	}

	// NOP
	FuncMap[0xEA] = func() {
		PC++
	}
}
