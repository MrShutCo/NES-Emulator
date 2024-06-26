package cpu

func Stack() {
	newInst(0x48, "PHA", "", 3, 1)
	newInst(0x08, "PHP", "", 3, 1)
	newInst(0x68, "PLA", "", 4, 1)
	newInst(0x28, "PLP", "", 4, 1)
	// PHA
	FuncMap[0x48] = func() {
		push(AC)
		PC++
	}
	// PHP
	FuncMap[0x08] = func() {
		push(SR | 0b0011_0000)
		PC++
	}
	// PLA
	FuncMap[0x68] = func() {
		SetAC(pull())
		//printStack()
		PC++
	}
	// PLP
	FuncMap[0x28] = func() {
		// Ignore bit 5
		oldSR := SR
		tempSR := pull()
		SR = tempSR
		setBreakFlag(getBit(oldSR, 4))
		setEmptyFlag(getBit(oldSR, 5))
		PC++
	}
}

func pull() byte {
	SP++
	data := RAM[STACK+uint16(SP)]
	//fmt.Printf("%02X\n", data)
	return data
}

func push(data byte) {
	RAM[STACK+uint16(SP)] = data
	SP--
}
