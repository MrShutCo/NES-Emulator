package cpu

func LoadInstructions() {
	NewInst(0xA9, "LDA", "immediate", 2, func() { SetAC(immed()) }, 2)
	NewInst(0xA5, "LDA", "zeropage", 3, func() { SetAC(zeropage()) }, 2)
	NewInst(0xB5, "LDA", "zeropage,X", 4, func() { SetAC(zeropageX()) }, 2)
	NewInst(0xAD, "LDA", "absolute", 4, func() { SetAC(absolute()) }, 3)
	NewInst(0xBD, "LDA", "absolute,X", 4, func() { SetAC(absoluteX()) }, 3)
	NewInst(0xB9, "LDA", "absolute,Y", 4, func() { SetAC(absoluteY()) }, 3)
	NewInst(0xA1, "LDA", "indirect,X", 6, func() { SetAC(indirectX()) }, 2)
	NewInst(0xB1, "LDA", "indirect,Y", 5, func() { SetAC(indirectY()) }, 2)

	NewInst(0xA2, "LDX", "immediate", 2, func() { SetX(immed()) }, 2)
	NewInst(0xA6, "LDX", "zeropage", 3, func() { SetX(zeropage()) }, 2)
	NewInst(0xB6, "LDX", "zeropage,Y", 4, func() { SetX(zeropageY()) }, 2)
	NewInst(0xAE, "LDX", "absolute", 4, func() { SetX(absolute()) }, 3)
	NewInst(0xBE, "LDX", "absolute,Y", 4, func() { SetX(absoluteY()) }, 3)

	NewInst(0xA0, "LDY", "immediate", 2, func() { SetY(immed()) }, 2)
	NewInst(0xA4, "LDY", "zeropage", 3, func() { SetY(zeropage()) }, 2)
	NewInst(0xB4, "LDY", "zeropage,X", 4, func() { SetY(zeropageX()) }, 2)
	NewInst(0xAC, "LDY", "absolute", 4, func() { SetY(absolute()) }, 3)
	NewInst(0xBC, "LDY", "absolute,X", 4, func() { SetY(absoluteX()) }, 3)

	NewInst(0xAA, "TAX", "implied", 2, func() { SetX(AC); PC++ }, 1)
	NewInst(0xA8, "TAY", "implied", 2, func() { SetY(AC); PC++ }, 1)
	NewInst(0xBA, "TSX", "implied", 2, func() { SetX(SP); PC++ }, 1)
	NewInst(0x8A, "TXA", "implied", 2, func() { SetAC(X); PC++ }, 1)
	NewInst(0x9A, "TXS", "implied", 2, func() { SP = X; PC++ }, 1)
	NewInst(0x98, "TYA", "implied", 2, func() { SetAC(Y); PC++ }, 1)
}
