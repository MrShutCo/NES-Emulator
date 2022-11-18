package cpu

func LDA() {
	newInst(0xA9, "LDA", "immediate", 2)
	newInst(0xA5, "LDA", "zeropage", 3)
	newInst(0xB5, "LDA", "zeropage,X", 4)
	newInst(0xAD, "LDA", "absolute", 4)
	newInst(0xBD, "LDA", "absolute,X", 4)
	newInst(0xB9, "LDA", "absolute,Y", 4)
	a := []foo{
		{0xA9, func() { SetAC(immed()) }},
		{0xA5, func() { SetAC(zeropage()) }},
		{0xB5, func() { SetAC(zeropageX()) }},
		{0xAD, func() { SetAC(absolute()) }},
		{0xBD, func() { SetAC(absoluteX()) }},
		{0xB9, func() { SetAC(absoluteY()) }},
	}
	apply(a)

}

func LDX() {
	newInst(0xA2, "LDX", "immediate", 2)
	newInst(0xA6, "LDX", "zeropage", 3)
	newInst(0xB6, "LDX", "zeropage,Y", 4)
	newInst(0xAE, "LDX", "absolute", 4)
	newInst(0xBE, "LDX", "absolute,Y", 4)
	a := []foo{
		{0xA2, func() { SetX(immed()) }},
		{0xA6, func() { SetX(zeropage()) }},
		{0xB6, func() { SetX(zeropageY()) }},
		{0xAE, func() { SetX(absolute()) }},
		{0xBE, func() { SetX(absoluteY()) }},
	}
	apply(a)
}

func LDY() {
	newInst(0xA0, "LDY", "immediate", 2)
	newInst(0xA4, "LDY", "zeropage", 3)
	newInst(0xB4, "LDY", "zeropage,X", 4)
	newInst(0xAC, "LDY", "absolute", 4)
	newInst(0xBC, "LDY", "absolute,X", 4)
	a := []foo{
		{0xA0, func() { SetY(immed()) }},
		{0xA4, func() { SetY(zeropage()) }},
		{0xB4, func() { SetY(zeropageY()) }},
		{0xAC, func() { SetY(absolute()) }},
		{0xBC, func() { SetY(absoluteY()) }},
	}
	apply(a)
}

func Transfer() {
	newInst(0xAA, "TAX", "implied", 2)
	newInst(0xA8, "TAY", "implied", 2)
	newInst(0xBA, "TSX", "implied", 2)
	newInst(0x8A, "TXA", "implied", 2)
	newInst(0x9A, "TXS", "implied", 2)
	newInst(0x98, "TYA", "implied", 2)
	// TAX
	FuncMap[0xAA] = func() { SetX(AC); PC++ }
	// TAY
	FuncMap[0xA8] = func() { SetY(AC); PC++ }
	// TSX
	FuncMap[0xBA] = func() { SetX(SP); PC++ }
	// TXA
	FuncMap[0x8A] = func() { SetAC(X); PC++ }
	// TXS
	FuncMap[0x9A] = func() { SP = X; PC++ }
	// TYA
	FuncMap[0x98] = func() { SetAC(Y); PC++ }
}
