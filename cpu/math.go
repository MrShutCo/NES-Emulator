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
}
