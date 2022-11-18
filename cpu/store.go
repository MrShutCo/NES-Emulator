package cpu

import "fmt"

func STA() {
	newInst(0x85, "STA", "zeropage", 3)
	newInst(0x95, "STA", "zeropage,X", 4)
	newInst(0x8D, "STA", "absolute", 4)
	newInst(0x9D, "STA", "absolute,X", 5)
	newInst(0x99, "STA", "absolute,Y", 5)
	a := []foo{
		{0x85, func() {
			RAM[RAM[PC+1]] = AC
			output = fmt.Sprintf("$%02X = %02X", RAM[PC+1], AC)
			PC += 2
		}},
		{0x95, func() {
			RAM[RAM[PC+1]+X] = AC
			//output = fmt.Sprintf("$%02X,X @ = %02X", RAM[PC+1], AC)
			PC += 2
		}},
		{0x8D, func() {
			RAM[getNextWord()] = AC
			PC += 3
		}},
		{0x9D, func() {
			RAM[getNextWord()+uint16(X)] = AC
			PC += 3
		}},
		{0x99, func() {
			RAM[getNextWord()+uint16(Y)] = AC
			PC += 3
		}},
	}
	apply(a)
}

func STX() {
	newInst(0x86, "STX", "zeropage", 2)
	newInst(0x96, "STX", "zeropage,Y", 3)
	newInst(0x8E, "STX", "absolute", 3)

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
		val := bytesToInt16(RAM[PC+2], RAM[PC+1])
		RAM[val] = X
		output = fmt.Sprintf("$%04X", val)
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
		val := bytesToInt16(RAM[PC+2], RAM[PC+1])
		RAM[val] = Y
		output = fmt.Sprintf("$%04X", val)
		PC += 3
	}
}
