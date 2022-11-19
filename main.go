package main

import (
	"6502/cpu"
	"fmt"
	"log"
	"os"
)

func main() {
	DoProgram()
}

func DoProgram() {
	//Load("test.nes")
	cpu.Reset()
	cpu.LoadMaps()
	//cpu.Load("nes-test-roms/tutor/tutor.nes")
	cpu.Load("nes_test/nestest.nes")
	fmt.Printf("0x%d\n", cpu.RAM[cpu.RES_VECTOR])
	cpu.Start()
	cpu.PC = 0xC000

	os.Truncate("access.log", 0)
	f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for {
		output := cpu.Execute()
		//fmt.Fscanln(os.Stdin)
		f.Write([]byte(output + "\n"))
		//fmt.Print(output + "\n")
		//fmt.Println(cpu.RAM[0x2000:0x2100])
	}
}
