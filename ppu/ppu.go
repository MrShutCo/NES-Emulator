package ppu

import (
	"6502/util"
	"fmt"
	"image/color"
)

var PPURAM [0x4000]byte

//var ppuCycles uint64

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

func GetImageFromPatternTable(val byte, bank uint16) []byte {
	start := bank + uint16(val)*16 // Get the image in question
	return PPURAM[start : start+16]
}

type PPU struct {
	latch       byte
	NMI_enabled bool
	cycles      int
	scanlines   int

	nametable              uint16
	vramIncrement          byte
	spritePatternTable     uint16
	backgroundPatternTable uint16
	is8x8Sprites           bool
	generateNMIatVBI       bool
}

func NewPPU() *PPU {
	screenBuffer = make([]byte, 256*256*4)
	return &PPU{
		vramIncrement:    1,
		is8x8Sprites:     true,
		generateNMIatVBI: true,
		nametable:        NAMETABLE_0,
	}
}

func (p *PPU) StepPPU(cycles byte) bool {
	p.cycles += int(cycles)
	if p.cycles > 340 {
		p.cycles -= 341
		p.scanlines++
	}

	// Before update we weren't in VB, but after we are
	if p.cycles-int(cycles) < 241 && p.cycles >= 241 {
		// Generate NMI interrupt if enabled
		if p.generateNMIatVBI {
			p.NMI_enabled = true
		}
		p.DrawBackground(0)
	}

	if p.scanlines >= 262 {
		p.scanlines = 0
		return true
	}
	return false
}

func ShowTile(tile []byte, startX, startY int) []byte {
	b := []byte{}
	for y := 0; y < 8; y++ {
		lower := tile[y]
		upper := tile[y+8]
		for x := 0; x < 8; x++ {
			colour := (lower >> (8 - x) & 1)
			colour += 2 * (upper >> (8 - x) & 1)
			posX := int(startX) + (x)
			posY := int(startY) + (y)
			switch colour {
			case 0x0:
				Image.Set(posX, posY, color.RGBA{50, 50, 50, 255})
			case 0x1:
				Image.Set(posX, posY, color.RGBA{100, 100, 100, 255})
			case 0x2:
				Image.Set(posX, posY, color.RGBA{150, 150, 150, 255})
			case 0x3:
				Image.Set(posX, posY, color.RGBA{255, 255, 255, 255})
			}
			b = append(b, colour)

		}
	}
	return b
}

func SetMemory(start uint16, data []byte) {
	for i := range data {
		PPURAM[start+uint16(i)] = data[i]
	}
	fmt.Printf("Length of data: %04X\n", len(data))
	fmt.Printf("%04X\n", start)
	util.PrintPage(data[:], 0x00)
}

func (p *PPU) ppuctrl(data byte) {
	// Which nametable to access
	nametable := data & 0x3
	p.nametable = 0x2000 + 0x400*uint16(nametable)

	// 0 if adding 1 for a PPUDATA write, 1 for adding 32
	incrPPUData := (data & 0b0100) >> 2
	if incrPPUData == 0 {
		p.vramIncrement = 1
	} else {
		p.vramIncrement = 32
	}

	// 0: $0000, 1: $1000, ignored in 8x16 mode
	spritePatternTable := (data & 0b1000) >> 3
	p.spritePatternTable = 0x1000 * uint16(spritePatternTable)

	// 0: $0000, 1: $1000, ignored in 8x16 mode
	backgroundPatternTable := (data & 0x10) >> 4
	p.backgroundPatternTable = 0x1000 * uint16(backgroundPatternTable)

	// 0: 8x8 pixels, 1: 8x16 pixels
	spriteSize := (data & 0x20) >> 5
	p.is8x8Sprites = spriteSize == 0

	// Unsure of next bit
	// 0: off, 1:
	generateNMIatVBI := (data & 0x80) >> 7
	p.generateNMIatVBI = generateNMIatVBI == 1
	//fmt.Println("Set PPUCTRL")
	//fmt.Printf("Sprite Table:     0x%04X\n", p.spritePatternTable)
	//fmt.Printf("Background Table: 0x%04X\n", p.backgroundPatternTable)
}

func (b *PPU) WriteBus(cpuAddr uint16, data byte) {
	b.latch = data
	switch cpuAddr {
	case 0x2000:
		b.ppuctrl(data)
	case 0x2006:
		_PPUADDR = _PPUADDR << 8           // Shift lo -> high
		_PPUADDR = _PPUADDR & 0xFF00       // Set lo = 0
		_PPUADDR = _PPUADDR | uint16(data) // Set lo bits
	case 0x2007:
		PPURAM[_PPUADDR] = data
		_PPUADDR += uint16(b.vramIncrement)
	}

}

func (b *PPU) ReadBus(cpuAddr uint16) byte {
	switch cpuAddr {
	case 0x2002:
		return b.ppustatus()
	}
	return b.latch
}

func (p *PPU) ppustatus() byte {
	_PPUADDR = 0x0
	if p.NMI_enabled {
		return 0xC0
	}
	return 0xC0
}
