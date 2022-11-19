package cpu

import "fmt"

func STA() {
	newInst(0x85, "STA", "zeropage", 3)
	newInst(0x95, "STA", "zeropage,X", 4)
	newInst(0x8D, "STA", "absolute", 4)
	newInst(0x9D, "STA", "absolute,X", 5)
	newInst(0x99, "STA", "absolute,Y", 5)
	newInst(0x81, "STA", "(indirect,X)", 5)

	a := []foo{
		{0x85, func() {
			oldVal := RAM[RAM[PC+1]]
			output = fmt.Sprintf("$%02X = %02X", RAM[PC+1], oldVal) // -1 since it was increased already
			if oldVal == 0xC6 {
				PrintPage(0x00)
				fmt.Printf("value at $%02X: %02X\n", RAM[PC+1], oldVal)
				fmt.Printf("value of AC: %02X\n", AC)
				fmt.Printf("value at op: %02X\n", RAM[PC+1])
			}
			RAM[RAM[PC+1]] = AC
			PC += 2
		}},
		{0x95, func() {
			RAM[RAM[PC+1]+X] = AC
			//output = fmt.Sprintf("$%02X,X @ = %02X", RAM[PC+1], AC)
			PC += 2
		}},
		{0x8D, func() {
			word := getNextWord()
			output = fmt.Sprintf("$%04X = %02X", word, RAM[word])
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
		{0x81, func() {
			RAM[indirectX16()] = AC
			PC += 2
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
		output = fmt.Sprintf("$%02X = %02X", RAM[PC+1], RAM[RAM[PC+1]])
		RAM[RAM[PC+1]] = X
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
		output = fmt.Sprintf("$%04X = %02X", val, RAM[val])
		RAM[val] = X
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
		output = fmt.Sprintf("$%02X = %02X", RAM[PC+1], RAM[RAM[PC+1]])
		RAM[RAM[PC+1]] = Y
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
		output = fmt.Sprintf("$%04X = %02X", val, RAM[val])
		RAM[val] = Y
		PC += 3
	}
}
