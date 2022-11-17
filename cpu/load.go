package cpu

import "fmt"

func LDA() {
	newInst(0xA9, "LDA", "immediate", 2)
	newInst(0xA5, "LDA", "zeropage", 3)
	newInst(0xB5, "LDA", "zeropage,X", 4)
	// LDA immeadiate
	FuncMap[0xA9] = func() {
		SetAC(RAM[PC+1])
		PC += 2
		output = fmt.Sprintf("#$%02X", AC)
	}
	// LDA zeropage
	FuncMap[0xA5] = func() {
		SetAC(RAM[RAM[PC+1]]) // Access $00XX
		PC += 2
	}
	// LDA zeropage,X
	FuncMap[0xB5] = func() {
		SetAC(RAM[RAM[PC+1]+X])
		PC += 2
	}
}

func LDX() {
	newInst(0xA2, "LDX", "immediate", 2)
	newInst(0xA6, "LDX", "zeropage", 3)
	newInst(0xB6, "LDX", "zeropage,Y", 4)
	newInst(0xAE, "LDX", "absolute", 4)
	newInst(0xBE, "LDX", "absolute,Y", 4)
	// LDX immeadiate
	FuncMap[0xA2] = func() {
		SetX(RAM[PC+1])
		PC += 2
		output = fmt.Sprintf("#$%02X", X)
	}
	// LDX zeropage
	FuncMap[0xA6] = func() {
		SetX(RAM[RAM[PC+1]]) // Access $00XX
		output = fmt.Sprintf("$%02X = %02X", RAM[PC+1], X)
		PC += 2
	}
	// LDX zeropage,Y
	FuncMap[0xB6] = func() {
		SetX(RAM[RAM[PC+1]+Y]) // Access $00XX+Y
		PC += 2
	}
	// LDX absolute
	FuncMap[0xAE] = func() {
		SetX(RAM[getNextWord()])
		PC += 3
	}
	// LDX absolute,Y
	FuncMap[0xBE] = func() {
		SetX(RAM[getNextWord()+uint16(Y)])
		PC += 3
	}
}

func LDY() {
	newInst(0xA0, "LDY", "immediate", 2)
	newInst(0xA4, "LDY", "zeropage", 3)
	newInst(0xB4, "LDY", "zeropage,X", 4)
	newInst(0xAC, "LDY", "absolute", 4)
	newInst(0xBC, "LDY", "absolute,X", 4)
	// LDY immeadiate
	FuncMap[0xA0] = func() {
		SetY(RAM[PC+1])
		PC += 2
		output = fmt.Sprintf("#$%02X", Y)
	}
	// LDY zeropage
	FuncMap[0xA4] = func() {
		SetY(RAM[RAM[PC+1]]) // Access $00YY
		PC += 2
	}
	// LDY zeropage,Y
	FuncMap[0xB4] = func() {
		SetY(RAM[RAM[PC+1]+X]) // Access $00YY+X
		PC += 2
	}
	// LDY absolute
	FuncMap[0xAC] = func() {
		SetY(RAM[getNextWord()])
		PC += 3
	}
	// LDY absolute,X
	FuncMap[0xBC] = func() {
		SetY(RAM[getNextWord()+uint16(X)])
		PC += 3
	}
}
