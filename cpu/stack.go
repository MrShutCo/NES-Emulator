package cpu

import "fmt"

func Stack() {
	newInst(0x48, "PHA", "", 3)
	newInst(0x08, "PHP", "", 3)
	newInst(0x68, "PLA", "", 4)
	newInst(0x28, "PLP", "", 4)
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

func printStack() {
	for i := 0x100; i <= 0x01FF; i++ {
		fmt.Printf("%02X,", RAM[i])
	}
	fmt.Print("\n")
}

func PrintPage(page byte) {
	start := uint16(page) << 8
	for y := uint16(0); y < 0x10; y++ {
		fmt.Printf("%04X:  ", start+y*0x10)
		for x := uint16(0); x < 0x10; x++ {
			fmt.Printf("%02X,", RAM[start+y*0x10+x])
		}
		fmt.Printf("\n")
	}
}
