package cpu

func Math() {
	// INX
	FuncMap[0xE8] = func() {
		X++
		if X == 0 {
			setZeroFlag(true)
			setNegativeFlag(true)
		}
	}
	// INY
	FuncMap[0xC8] = func() {
		Y++
		if Y == 0 {
			setZeroFlag(true)
			setNegativeFlag(true)
		}
	}
	// INC zeropage
	FuncMap[0xE6] = func() {
		RAM[RAM[PC+1]]++
		if RAM[RAM[PC+1]] == 0 {
			setZeroFlag(true)
			setNegativeFlag(true)
		}
	}
	// INC zeropage,X
	FuncMap[0xF6] = func() {
		RAM[RAM[PC+1]+X]++
		if RAM[RAM[PC+1]+X] == 0 {
			setZeroFlag(true)
			setNegativeFlag(true)
		}
	}
	// INC absolute
	FuncMap[0xEE] = func() {
		RAM[getNextWord()]++
		if RAM[getNextWord()] == 0 {
			setZeroFlag(true)
			setNegativeFlag(true)
		}
	}
	// INC absolute,X
	FuncMap[0xFE] = func() {
		RAM[getNextWord()+uint16(X)]++
		if RAM[getNextWord()+uint16(X)] == 0 {
			setZeroFlag(true)
			setNegativeFlag(true)
		}
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
}

func adc(f func() byte) {
	p := f()
	val := uint16(AC + p)
	if isCarrySet() {
		val++
		//p++
	}

	setCarryFlag(val > 0xFF)
	setOverflowFlag(((AC^p)&0x80 == 0) && ((AC^byte(val))&0x80 != 0))
	SetAC(byte(val))
}
