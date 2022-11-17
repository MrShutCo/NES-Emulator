package cpu

import "fmt"

func Stack() {
	newInst(0x48, "PHA", "", 3)
	newInst(0x08, "PHP", "", 3)
	newInst(0x68, "PLA", "", 3)
	newInst(0x28, "PLP", "", 3)
	// PHA
	FuncMap[0x48] = func() {
		push(AC)
		PC++
	}
	// PHP
	FuncMap[0x08] = func() {
		push(SR)
		//SR = SR //| 0b0000_0000 // TODO: is this right?
		PC++
	}
	// PLA
	FuncMap[0x68] = func() {
		AC = pull()
		PC++
	}
	// PLP
	// TODO: determine what break flag does???
	FuncMap[0x28] = func() {
		// Ignore bit 5
		SR = RAM[STACK+uint16(SP)] //| (SR & 0b_0001_0000)
		SP++
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
