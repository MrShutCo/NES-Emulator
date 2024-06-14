package ppu

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

var PPURAM [0x4000]byte
var OAM [0x100]byte

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

var ColorList []color.RGBA

func GetImageFromPatternTable(val byte, bank uint16) []byte {
	start := bank + uint16(val)*16 // Get the image in question
	return PPURAM[start : start+16]
}

func GetSpritePalette(paletteID byte) color.Palette {
	addr := 0x3F11 + 4*uint16(paletteID)
	paletteData := PPURAM[addr : addr+3]
	return color.Palette{
		color.Transparent, ColorList[paletteData[0]&0b00111111],
		ColorList[paletteData[1]&0b00111111], ColorList[paletteData[2]&0b00111111],
	}
}

func GetBackgroundPalette(paletteID byte) color.Palette {
	addr := 0x3F01 + 4*uint16(paletteID)
	paletteData := PPURAM[addr : addr+3]
	return color.Palette{
		ColorList[PPURAM[0x3F00]&0b00111111], ColorList[paletteData[0]&0b00111111],
		ColorList[paletteData[1]&0b00111111], ColorList[paletteData[2]&0b00111111],
	}
}

func GetBackgroundPaletteID(tileIndex int) byte {
	x := (tileIndex % 32) / 4
	y := (tileIndex / 32) / 4
	address := 0x23C0 + x + 8*y
	attributeByte := PPURAM[address]

	// Now we need to figure out what quadrant we are in
	cellX := (tileIndex % 32) % 4
	cellY := (tileIndex / 32) % 4

	if cellX < 2 && cellY < 2 {
		return attributeByte & 0x03
	}
	if cellX >= 2 && cellY < 2 {
		return (attributeByte >> 2) & 0x03
	}
	if cellX < 2 && cellY >= 2 {
		return (attributeByte >> 4) & 0x03
	}
	return (attributeByte >> 6) & 0x03
}

func getAttrByteFromTileIndex(tileIndex int) byte {
	// position of the 32x32 block
	tileX := tileIndex % 32
	tileY := tileIndex / 32
	attrX := tileX / 8
	attrY := (tileY / 8)
	return byte(attrX) + byte(attrY)*8
}

func get2BitsFromByte(data byte, bit byte) byte {
	return data & (0x3 << bit) >> bit
}

type TileCache struct {
	NametableIndex int
	Palette        color.Palette
	Tile           *ebiten.Image
}

type PPU struct {
	latch       byte
	NMI_enabled bool
	cycles      int
	scanlines   int
	isEvenFrame bool

	nametable              uint16
	vramIncrement          byte
	spritePatternTable     uint16
	backgroundPatternTable uint16
	is8x8Sprites           bool
	nmi_occurred           bool
	nmi_output             bool

	patternTable0SpriteSheet *ebiten.Image
	patternTable1SpriteSheet *ebiten.Image

	cache []TileCache

	// New logic to support colour
	// Contains the palette indexes of all background sprites
	pattern0 []byte
	pattern1 []byte
}

func NewPPU() *PPU {

	/*preset := []byte{
		0x80, 0x80, 0x80, 0x00, 0x3D, 0xA6, 0x00, 0x12, 0xB0, 0x44, 0x00, 0x96, 0xA1, 0x00, 0x5E,
		0xC7, 0x00, 0x28, 0xBA, 0x06, 0x00, 0x8C, 0x17, 0x00, 0x5C, 0x2F, 0x00, 0x10, 0x45, 0x00,
		0x05, 0x4A, 0x00, 0x00, 0x47, 0x2E, 0x00, 0x41, 0x66, 0x00, 0x00, 0x00, 0x05, 0x05, 0x05,
		0x05, 0x05, 0x05, 0xC7, 0xC7, 0xC7, 0x00, 0x77, 0xFF, 0x21, 0x55, 0xFF, 0x82, 0x37, 0xFA,
		0xEB, 0x2F, 0xB5, 0xFF, 0x29, 0x50, 0xFF, 0x22, 0x00, 0xD6, 0x32, 0x00, 0xC4, 0x62, 0x00,
		0x35, 0x80, 0x00, 0x05, 0x8F, 0x00, 0x00, 0x8A, 0x55, 0x00, 0x99, 0xCC, 0x21, 0x21, 0x21,
		0x09, 0x09, 0x09, 0x09, 0x09, 0x09, 0xFF, 0xFF, 0xFF, 0x0F, 0xD7, 0xFF, 0x69, 0xA2, 0xFF,
		0xD4, 0x80, 0xFF, 0xFF, 0x45, 0xF3, 0xFF, 0x61, 0x8B, 0xFF, 0x88, 0x33, 0xFF, 0x9C, 0x12,
		0xFA, 0xBC, 0x20, 0x9F, 0xE3, 0x0E, 0x2B, 0xF0, 0x35, 0x0C, 0xF0, 0xA4, 0x05, 0xFB, 0xFF,
		0x5E, 0x5E, 0x5E, 0x0D, 0x0D, 0x0D, 0x0D, 0x0D, 0x0D, 0xFF, 0xFF, 0xFF, 0xA6, 0xFC, 0xFF,
		0xB3, 0xEC, 0xFF, 0xDA, 0xAB, 0xEB, 0xFF, 0xA8, 0xF9, 0xFF, 0xAB, 0xB3, 0xFF, 0xD2, 0xB0,
		0xFF, 0xEF, 0xA6, 0xFF, 0xF7, 0x9C, 0xD7, 0xE8, 0x95, 0xA6, 0xED, 0xAF, 0xA2, 0xF2, 0xDA,
		0x99, 0xFF, 0xFC, 0xDD, 0xDD, 0xDD, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
	}*/
	preset := []byte{
		0x7C, 0x7C, 0x7C,
		0x00, 0x00, 0xFC,
		0x00, 0x00, 0xBC,
		0x44, 0x28, 0xBC,
		0x94, 0x00, 0x84,
		0xA8, 0x00, 0x20,
		0xA8, 0x10, 0x00,
		0x88, 0x14, 0x00,
		0x50, 0x30, 0x00,
		0x00, 0x78, 0x00,
		0x00, 0x68, 0x00,
		0x00, 0x58, 0x00,
		0x00, 0x40, 0x58,
		0x00, 0x00, 0x00,
		0x00, 0x00, 0x00,
		0x00, 0x00, 0x00,
		0xBC, 0xBC, 0xBC,
		0x00, 0x78, 0xF8,
		0x00, 0x58, 0xF8,
		0x68, 0x44, 0xFC,
		0xD8, 0x00, 0xCC,
		0xE4, 0x00, 0x58,
		0xF8, 0x38, 0x00,
		0xE4, 0x5C, 0x10,
		0xAC, 0x7C, 0x00,
		0x00, 0xB8, 0x00,
		0x00, 0xA8, 0x00,
		0x00, 0xA8, 0x44,
		0x00, 0x88, 0x88,
		0x00, 0x00, 0x00,
		0x00, 0x00, 0x00,
		0x00, 0x00, 0x00,
		0xF8, 0xF8, 0xF8,
		0x3C, 0xBC, 0xFC,
		0x68, 0x88, 0xFC,
		0x98, 0x78, 0xF8,
		0xF8, 0x78, 0xF8,
		0xF8, 0x58, 0x98,
		0xF8, 0x78, 0x58,
		0xFC, 0xA0, 0x44,
		0xF8, 0xB8, 0x00,
		0xB8, 0xF8, 0x18,
		0x58, 0xD8, 0x54,
		0x58, 0xF8, 0x98,
		0x00, 0xE8, 0xD8,
		0x78, 0x78, 0x78,
		0x00, 0x00, 0x00,
		0x00, 0x00, 0x00,
		0xFC, 0xFC, 0xFC,
		0xA4, 0xE4, 0xFC,
		0xB8, 0xB8, 0xF8,
		0xD8, 0xB8, 0xF8,
		0xF8, 0xB8, 0xF8,
		0xF8, 0xA4, 0xC0,
		0xF0, 0xD0, 0xB0,
		0xFC, 0xE0, 0xA8,
		0xF8, 0xD8, 0x78,
		0xD8, 0xF8, 0x78,
		0xB8, 0xF8, 0xB8,
		0xB8, 0xF8, 0xD8,
		0x00, 0xFC, 0xFC,
		0xF8, 0xD8, 0xF8,
		0x00, 0x00, 0x00,
		0x00, 0x00, 0x00,
	}
	ColorList = make([]color.RGBA, len(preset)/3)

	for i := byte(0); i <= 0x3F; i++ {
		pi := i * 3
		ColorList[i] = color.RGBA{preset[pi], preset[pi+1], preset[pi+2], 255}
	}
	//cache = make(map[uint16]*ebiten.Image)

	return &PPU{
		vramIncrement: 1,
		is8x8Sprites:  true,
		isEvenFrame:   true,
		nmi_output:    true,
		nmi_occurred:  false,
		nametable:     NAMETABLE_0,
		cache:         make([]TileCache, 0x3C0),
	}
}

func (p *PPU) ShouldTriggerNMI() bool {
	return p.nmi_occurred && p.nmi_output
}

func (p *PPU) StepPPU(cycles byte) bool {
	p.cycles += int(cycles)
	if p.cycles > 340 {
		p.cycles -= 341
		p.scanlines++
	}

	if p.scanlines >= 8 && p.scanlines < 240 && p.scanlines%8 == 0 {
		p.DrawBackgroundRow(p.scanlines/8 - 1)
	}

	// Before update we weren't in VB, but after we are
	if p.scanlines == 241 {
		// Generate NMI interrupt if enabled
		//if p.nmi_occurred && p.nmi_output {
		//	p.NMI_enabled = true
		//}

		p.nmi_occurred = true

	}

	p.MirrorMemory()

	if p.scanlines >= 262 {
		p.scanlines = 0
		p.nmi_occurred = false
		// Even/odd frame counting
		p.isEvenFrame = !p.isEvenFrame
		return true
	}
	return false
}

func (p *PPU) MirrorMemory() {
	// Palette Mirroring
	PPURAM[0x3F10] = PPURAM[0x3F00]
	PPURAM[0x3F14] = PPURAM[0x3F04]
	PPURAM[0x3F18] = PPURAM[0x3F08]
	PPURAM[0x3F1C] = PPURAM[0x3F0C]
}

func SetMemory(start uint16, data []byte) {
	for i := range data {
		PPURAM[start+uint16(i)] = data[i]
	}
	fmt.Printf("Length of data: %04X\n", len(data))
	fmt.Printf("%04X\n", start)
}

// LoadPaletteV2 - starts scanning the PPURAM at the given address 16 bytes at a time
//
//	to load a single 8x8 tile. Two bits are read that are 8 bytes apart at a time and used
//	to index into the palette. p.pattern0 and p.pattern1 are byte arrays of the palette colour per pixel.
//	All the writes are done sequentually in memory, tile by tile so fetching 64 bytes sequentially in pattern
//	will correspond to a tile.
func (p *PPU) LoadPaletteV2(table uint16) {
	pos := 0
	palette := make([]byte, 256*256)
	for i := 0; i < 0x100; i++ {
		tileData := GetImageFromPatternTable(byte(i), table)
		for y := 0; y < 8; y++ {
			lower := tileData[y]
			upper := tileData[y+8]
			for x := 0; x < 8; x++ {
				colour := (lower >> (7 - x) & 1)
				colour += 2 * (upper >> (7 - x) & 1)
				palette[pos] = colour
				pos++
			}
		}
	}
	if table == PATTERN_TABLE_0 {
		p.pattern0 = palette
	} else if table == PATTERN_TABLE_1 {
		p.pattern1 = palette
	}
}

func (p *PPU) PreloadPalleteTable(table uint16) {
	spriteSheet, _ := ebiten.NewImage(256, 256, ebiten.FilterDefault)
	for i := 0; i <= 0xFF; i++ {
		tileData := GetImageFromPatternTable(byte(i), table)
		data := make([]byte, 64*4)
		for y := 0; y < 8; y++ {
			lower := tileData[y]
			upper := tileData[y+8]
			for x := 0; x < 8; x++ {
				colour := (lower >> (7 - x) & 1)
				colour += 2 * (upper >> (7 - x) & 1)
				pos := y*32 + x*4
				switch colour {
				case 0x0:
					data[pos], data[pos+1], data[pos+2], data[pos+3] = 50, 50, 50, 255
				case 0x1:
					data[pos], data[pos+1], data[pos+2], data[pos+3] = 100, 100, 100, 255
				case 0x2:
					data[pos], data[pos+1], data[pos+2], data[pos+3] = 150, 150, 150, 255
				case 0x3:
					data[pos], data[pos+1], data[pos+2], data[pos+3] = 255, 255, 255, 255
				}
			}
		}
		img, _ := ebiten.NewImage(8, 8, ebiten.FilterDefault)
		img.ReplacePixels(data)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i%16)*8, float64(i/16)*8)
		spriteSheet.DrawImage(img, op)

	}
	if table == PATTERN_TABLE_0 {
		p.patternTable0SpriteSheet = spriteSheet
	} else if table == PATTERN_TABLE_1 {
		p.patternTable1SpriteSheet = spriteSheet
	}
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
	p.nmi_output = generateNMIatVBI == 1
	//fmt.Println("Set PPUCTRL")
	//fmt.Printf("Sprite Table:     0x%04X\n", p.spritePatternTable)
	//fmt.Printf("Background Table: 0x%04X\n", p.backgroundPatternTable)
}

func (p *PPU) OAMDMA(data [0x100]byte) {
	OAM = data
}

func (b *PPU) WriteBus(cpuAddr uint16, data byte) {
	b.latch = data
	switch cpuAddr {
	case 0x2000:
		b.ppuctrl(data)
	case 0x2005:

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
	if p.nmi_occurred {
		p.nmi_occurred = false
		return 0xC0 //TODO: race condition
	}
	return 0b11000000
}
