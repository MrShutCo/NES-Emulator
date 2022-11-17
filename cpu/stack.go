package cpu

import "fmt"

func Stack() {
	// PHA
	FuncMap[0x48] = func() {
		push(AC)
	}
	// PHP
	FuncMap[0x68] = func() {
		push(SR)
		SR = SR | 0b0001_0000 // TODO: is this right?
	}
	// PLA
	FuncMap[0x68] = func() {
		AC = pull()
	}
	// PLP
	// TODO: determine what break flag does???
	FuncMap[0x28] = func() {
		// Ignore bit 5
		SR = RAM[STACK+uint16(SP)] | (SP & 0b_0001_0000)
		SP++
	}
}

func pull() byte {
	data := RAM[STACK+uint16(SP)]
	SP++
	return data
}

func push(data byte) {
	RAM[STACK+uint16(SP)] = data
	SP--
}

func printStack() {
	for i := 0x0100; i < 0x01FF; i++ {
		fmt.Printf("%02X,", RAM[i])
	}
}
