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
}

func (n *NES) Simulate() {
	if n.PPU.NMI_enabled && !n.hasInterrupted {
		n.hasInterrupted = true
		n.interrupt()
	}
	oldCycles := cpu.Cycles
	cpu.Execute()
	//n.stdout += o + "\n"
	newCycles := cpu.Cycles
	doneDrawing1 := n.PPU.StepPPU(byte(newCycles - oldCycles))
	doneDrawing2 := n.PPU.StepPPU(byte(newCycles - oldCycles))
	doneDrawing3 := n.PPU.StepPPU(byte(newCycles - oldCycles))
	if doneDrawing1 || doneDrawing2 || doneDrawing3 {
		n.hasInterrupted = false
	}
}

func (n *NES) interrupt() {
	cpu.NMI_Interrupt()
}
