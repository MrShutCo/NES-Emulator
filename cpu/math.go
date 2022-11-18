package cpu

// TODO: refactor the INC and DEC
func Math() {
	newInst(0xE6, "INC", "immediate", 5)
	newInst(0xF6, "INC", "zeropage", 6)
	newInst(0xEE, "INC", "absolute", 6)
	newInst(0xFE, "INC", "absolute,X", 7)
	newInst(0xE8, "INX", "implied", 2)
	newInst(0xC8, "INY", "implied", 2)
	// INX
	FuncMap[0xE8] = func() {
		SetX(X + 1)
		PC++
	}
	// INY
	FuncMap[0xC8] = func() {
		SetY(Y + 1)
		PC++
	}
	// INC zeropage
	FuncMap[0xE6] = func() {
		RAM[RAM[PC+1]]++
		if RAM[RAM[PC+1]] == 0 {
			setZeroFlag(true)
			setNegativeFlag(true)
		}
		PC += 2
	}
	// INC zeropage,X
	FuncMap[0xF6] = func() {
		RAM[RAM[PC+1]+X]++
		if RAM[RAM[PC+1]+X] == 0 {
			setZeroFlag(true)
			setNegativeFlag(true)
		}
		PC += 2
	}
	// INC absolute
	FuncMap[0xEE] = func() {
		RAM[getNextWord()]++
		if RAM[getNextWord()] == 0 {
			setZeroFlag(true)
			setNegativeFlag(true)
		}
		PC += 3
	}
	// INC absolute,X
	FuncMap[0xFE] = func() {
		RAM[getNextWord()+uint16(X)]++
		if RAM[getNextWord()+uint16(X)] == 0 {
			setZeroFlag(true)
			setNegativeFlag(true)
		}
		PC += 3
	}

	newInst(0xCA, "DEX", "immediate", 5)
	newInst(0x88, "DEY", "zeropage", 6)
	newInst(0xC6, "DEC", "absolute", 6)
	newInst(0xD6, "DEC", "absolute,X", 7)
	newInst(0xCE, "DEC", "implied", 2)
	newInst(0xDE, "DEC", "implied", 2)
	// DEX
	FuncMap[0xCA] = func() {
		SetX(X - 1)
		PC++
	}
	// DEY
	FuncMap[0x88] = func() {
		SetY(Y - 1)
		PC++
	}
	// DEC zeropage
	FuncMap[0xC6] = func() {
		RAM[RAM[PC+1]]--
		if RAM[RAM[PC+1]] == 0 {
			setZeroFlag(true)
			setNegativeFlag(true)
		}
		PC += 2
	}
	// DEC zeropage,X
	FuncMap[0xD6] = func() {
		RAM[RAM[PC+1]+X]--
		if RAM[RAM[PC+1]+X] == 0 {
			setZeroFlag(true)
			setNegativeFlag(true)
		}
		PC += 2
	}
	// DEC absolute
	FuncMap[0xCE] = func() {
		RAM[getNextWord()]--
		if RAM[getNextWord()] == 0 {
			setZeroFlag(true)
			setNegativeFlag(true)
		}
		PC += 3
	}
	// DEC absolute,X
	FuncMap[0xDE] = func() {
		RAM[getNextWord()+uint16(X)]--
		if RAM[getNextWord()+uint16(X)] == 0 {
			setZeroFlag(true)
			setNegativeFlag(true)
		}
		PC += 3
	}

	newInst(0x69, "ADC", "immediate", 2)
	newInst(0x65, "ADC", "zeropage", 3)
	newInst(0x75, "ADC", "zeropage,X", 4)
	newInst(0x6D, "ADC", "absolute", 4)
	newInst(0x7D, "ADC", "absolute,X", 4)
	newInst(0x79, "ADC", "absolute,Y", 4)
	ad := []foo{
		{0x69, func() { adc(immed) }},
		{0x65, func() { adc(zeropage) }},
		{0x75, func() { adc(zeropageX) }},
		{0x6D, func() { adc(absolute) }},
		{0x7D, func() { adc(absoluteX) }},
		{0x79, func() { adc(absoluteY) }},
	}
	apply(ad)

	newInst(0xE9, "SBC", "immediate", 2)
	newInst(0xE5, "SBC", "zeropage", 3)
	newInst(0xF5, "SBC", "zeropage,X", 4)
	newInst(0xED, "SBC", "absolute", 4)
	newInst(0xFD, "SBC", "absolute,X", 4)
	newInst(0xF9, "SBC", "absolute,Y", 4)
	sub := []foo{
		{0xE9, func() { sbc(immed) }},
		{0xE5, func() { sbc(zeropage) }},
		{0xF5, func() { sbc(zeropageX) }},
		{0xED, func() { sbc(absolute) }},
		{0xFD, func() { sbc(absoluteX) }},
		{0xF9, func() { sbc(absoluteY) }},
	}
	apply(sub)
}

func sbc(f func() byte) {
	val := f()
	adc(func() byte { return 255 - val })
}

func adc(f func() byte) {
	p := f()
	val := uint16(AC) + uint16(p)
	if isCarrySet() {
		val++
	}
	setCarryFlag(val > 0xFF)
	setOverflowFlag(^(AC^p)&(AC^byte(val))&0x80 != 0)
	SetAC(byte(val))
}
