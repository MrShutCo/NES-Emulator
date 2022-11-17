package main

import (
	"6502/cpu"
	"6502/ppu"
)

type NES struct {
	MasterCycle uint64
}

func (n *NES) Step(count uint) {
	instruct := cpu.RAM[cpu.PC]
	cpu.FuncMap[instruct]()

	n.MasterCycle += uint64(cpu.Instructions[instruct].Cycles)

	ppu.CatchupToCurrent(n.MasterCycle)
}
