package cpu

import "fmt"

func STA() {
	newInst(0x85, "STA", "zeropage", 2)
	newInst(0x95, "STA", "zeropage,Y", 3)
	newInst(0x8D, "STA", "absolute", 3)

	// STA zeropage
	FuncMap[0x85] = func() {
		RAM[RAM[PC+1]] = AC
		output = fmt.Sprintf("$%02X = %02X", RAM[PC+1], AC)
		PC += 2
	}
	// STA zeropage,X
	FuncMap[0x95] = func() {
		RAM[RAM[PC+1]+X] = AC
		PC += 2
	}
	// STA absolute
	FuncMap[0x8D] = func() {
		RAM[bytesToInt16(RAM[PC+2], RAM[PC+1])] = AC
		PC += 3
	}
}

func STX() {
	newInst(0x86, "STX", "zeropage", 2)
	newInst(0x96, "STX", "zeropage,Y", 3)
	newInst(0x8E, "STX", "absolute", 3)

	/*x := map[byte]func() byte{
		0x86: zeropage,
		0x96: zeropageY,
		0x8E: absolute,
	}
	apply(stx, x)
	*/
	// STX zeropage
	FuncMap[0x86] = func() {
		RAM[RAM[PC+1]] = X
		output = fmt.Sprintf("$%02X = %02X", RAM[PC+1], X)
		PC += 2
	}
	// STX zeropage,Y
	FuncMap[0x96] = func() {
		RAM[RAM[PC+1]+Y] = X
		PC += 2
	}
	// STX absolute
	FuncMap[0x8E] = func() {
		RAM[bytesToInt16(RAM[PC+2], RAM[PC+1])] = X
		PC += 3
	}
}

func stx(f func() byte) {
	oldX := X
	X = f()
	output += fmt.Sprintf("%02X", oldX)
}

func STY() {
	newInst(0x84, "STY", "zeropage", 2)
	newInst(0x94, "STY", "zeropage,X", 3)
	newInst(0x8C, "STY", "absolute", 3)
	// STY zeropage
	FuncMap[0x84] = func() {
		RAM[RAM[PC+1]] = Y
		output = fmt.Sprintf("$%02X = %02X", RAM[PC+1], Y)
		PC += 2
	}
	// STY zeropage,X
	FuncMap[0x94] = func() {
		RAM[RAM[PC+1]+Y] = Y
		PC += 2
	}
	// STY absolute
	FuncMap[0x8C] = func() {
		RAM[bytesToInt16(RAM[PC+2], RAM[PC+1])] = Y
		PC += 3
	}
}
