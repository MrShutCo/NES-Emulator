package main

import (
	"6502/cpu"
	"6502/ppu"
	"fmt"
)

type NES struct {
	MasterCycle    uint64
	PPU            *ppu.PPU
	IsRunning      bool
	hasInterrupted bool
	stdout         string

	controllerInput byte
}

func (n *NES) Simulate() {
	if n.PPU.ShouldTriggerNMI() && !n.hasInterrupted {
		n.hasInterrupted = true // TODO: this may not be fully correct
		n.interrupt()
	}
	oldCycles := cpu.Cycles
	output := cpu.Execute()
	if cpu.OutputCommands {
		fmt.Println(output)
	}
	cycleDiff := cpu.Cycles - oldCycles
	for i := 0; i < int(cycleDiff)*3; i++ {
		doneDrawing := n.PPU.StepPPU()
		if doneDrawing {
			n.hasInterrupted = false
		}
	}

}

func (n *NES) interrupt() {
	cpu.NMI_Interrupt()
}
