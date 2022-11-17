package main

import (
	"6502/cpu"
	"fmt"
	"log"
	"os"
)

func main() {
	//Load("test.nes")
	cpu.Reset()
	cpu.LoadMaps()
	//cpu.Load("nes-test-roms/tutor/tutor.nes")
	cpu.Load("nes_test/nestest.nes")
	fmt.Printf("0x%d\n", cpu.RAM[cpu.RES_VECTOR])
	cpu.Start()
	cpu.PC = 0xC000

	f, _ := os.Create("output.txt")
	defer f.Close()

	f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for {
		cpu.Execute()
		fmt.Fscanln(os.Stdin)
		//start := fmt.Sprintf("%04X %02X %02X %02x  ", cpu.PC, cpu.RAM[cpu.PC], cpu.RAM[cpu.PC+1], cpu.RAM[cpu.PC+2])
		//middle := fmt.Sprintf("%s", cpu.Instructions[cpu.RAM[cpu.PC]])
		//end := fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X", cpu.AC, cpu.X, cpu.Y, cpu.SP)
		//s := ""
		//f.Write([]byte(start + middle + "\n"))
		//fmt.Println(cpu.RAM[0x2000:0x2100])
	}
}
