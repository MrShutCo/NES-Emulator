package cpu

import "fmt"

func STA() {
	newInst(0x85, "STA", "zeropage", 3, 2)
	newInst(0x95, "STA", "zeropage,X", 4, 2)
	newInst(0x8D, "STA", "absolute", 4, 3)
	newInst(0x9D, "STA", "absolute,X", 5, 3)
	newInst(0x99, "STA", "absolute,Y", 5, 3)
	newInst(0x81, "STA", "(indirect,X)", 6, 2)
	newInst(0x91, "STA", "(indirect),Y", 6, 2)
	a := []foo{
		{0x85, func() {
			oldVal := RAM[RAM[PC+1]]
			if OutputCommands {
				output = fmt.Sprintf("$%02X = %02X", RAM[PC+1], oldVal) // -1 since it was increased already
			}
			SetRAM(uint16(RAM[PC+1]), AC)
			PC += 2
		}},
		{0x95, func() {
			addr := zeropageXAddr()
			SetRAM(addr, AC)
			PC += 2
		}},
		{0x8D, func() {
			word := getNextWord()
			if OutputCommands {
				output = fmt.Sprintf("$%04X = %02X", word, RAM[word])
			}

			SetRAM(word, AC)
			PC += 3
		}},
		{0x9D, func() {
			addr := absoluteXAddr()
			SetRAM(addr, AC)
			PC += 3
		}},
		{0x99, func() {
			addr := absoluteYAddr()
			SetRAM(addr, AC)
			PC += 3
		}},
		{0x81, func() {
			addr := indirectX16()
			SetRAM(addr, AC)
			PC += 2
		}},
		{0x91, func() {
			addr := indirectYAddr()
			SetRAM(addr, AC)
			PC += 2
		}},
	}
	apply(a)
}

func STX() {
	newInst(0x86, "STX", "zeropage", 3, 2)
	newInst(0x96, "STX", "zeropage,Y", 4, 2)
	newInst(0x8E, "STX", "absolute", 4, 3)

	// STX zeropage
	FuncMap[0x86] = func() {
		addr := zeropageAddr()
		SetRAM(addr, X)
		PC += 2
	}
	// STX zeropage,Y
	FuncMap[0x96] = func() {
		addr := zeropageYAddr()
		SetRAM(addr, X)
		PC += 2
	}
	// STX absolute
	FuncMap[0x8E] = func() {
		addr := absoluteAddr()
		SetRAM(addr, X)
		PC += 3
	}
}

func STY() {
	newInst(0x84, "STY", "zeropage", 3, 2)
	newInst(0x94, "STY", "zeropage,X", 4, 2)
	newInst(0x8C, "STY", "absolute", 4, 3)
	// STY zeropage
	FuncMap[0x84] = func() {
		addr := zeropageAddr()
		SetRAM(addr, Y)
		PC += 2
	}
	// STY zeropage,X
	FuncMap[0x94] = func() {
		addr := zeropageXAddr()
		SetRAM(addr, Y)
		PC += 2
	}
	// STY absolute
	FuncMap[0x8C] = func() {
		addr := absoluteAddr()
		SetRAM(addr, Y)
		PC += 3
	}
}
