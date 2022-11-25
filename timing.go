package main

import (
	"6502/cpu"
	"6502/ppu"
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
		n.hasInterrupted = true // TODO: this may not be fully correcy
		n.interrupt()
	}
	oldCycles := cpu.Cycles
	cpu.Execute()
	cycleDiff := cpu.Cycles - oldCycles
	doneDrawing1 := n.PPU.StepPPU(byte(cycleDiff * 3))
	if doneDrawing1 {
		n.hasInterrupted = false
	}
}

func (n *NES) interrupt() {
	cpu.NMI_Interrupt()
}
