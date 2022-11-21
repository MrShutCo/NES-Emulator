package ppu

import (
	"6502/util"
	"fmt"
)

var PPURAM [0x4000]byte
var ppuCycles uint64

var _PPUCONTROL byte
var _PPUMASK byte
var _PPUSTATUS byte
var _OAMADDR byte
var _OAMDATA byte
var _PPUSCROLL byte
var _PPUADDR uint16
var _PPUDATA byte

const PATTERN_TABLE_0 = 0x0
const PATTERN_TABLE_1 = 0x1000
const NAMETABLE_0 = 0x2000
const NAMETABLE_1 = 0x2400
const NAMETABLE_2 = 0x2800
const NAMETABLE_3 = 0x2C00
const MIRRODED_REGISTERS = 0x3000
const PALETTE_RAM_INDEXES = 0x4000
const MIRRORS_OF_PALETTE = 0x3F20

var DataStruct *PPU

func GetCell(x byte, y byte) byte {
	var start = int(NAMETABLE_0) + int(30*y) + int(x)
	return PPURAM[start]
}

type PPU struct {
	latch byte
}

func SetMemory(start uint16, data []byte) {
	for i := range data {
		PPURAM[start+uint16(i)] = data[i]
	}
	fmt.Printf("Length of data: %04X\n", len(data))
	fmt.Printf("%04X\n", start)
	util.PrintPage(data[:], 0x00)
}

func (b *PPU) ReadBus(cpuAddr uint16) byte {
	switch cpuAddr {
	case 0x2007:
		// HUH
	}
	return b.latch
}

func (b *PPU) WriteBus(cpuAddr uint16, data byte) {
	b.latch = data
	switch cpuAddr {
	case 0x2006:
		_PPUADDR = _PPUADDR << 8           // Shift lo -> high
		_PPUADDR = _PPUADDR & 0xFF00       // Set lo = 0
		_PPUADDR = _PPUADDR | uint16(data) // Set lo bits
	case 0x2007:
		PPURAM[_PPUADDR] = data
		_PPUADDR += 1 // TODO: determined by bit 2 of $2000
	}

}

func CatchupToCurrent(masterCycle uint64) {

}
