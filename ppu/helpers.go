package ppu

import (
	"fmt"
	"image/color"
)

func GetImageFromPatternTable(val byte, bank uint16) []byte {
	start := bank + uint16(val)*16 // Get the image in question
	return PPURAM[start : start+16]
}

func GetSpritePalette(paletteID byte) (color.RGBA, color.RGBA, color.RGBA, color.RGBA) {
	addr := 0x3F11 + 4*uint16(paletteID)
	paletteData := PPURAM[addr : addr+3]
	return color.RGBA{0, 0, 0, 0}, ColorList[paletteData[0]&0b00111111],
		ColorList[paletteData[1]&0b00111111], ColorList[paletteData[2]&0b00111111]
}

func GetSpritePalettePacked(paletteID byte) color.Palette {
	addr := 0x3F11 + 4*uint16(paletteID)
	paletteData := PPURAM[addr : addr+3]
	return color.Palette{color.RGBA{0, 0, 0, 0}, ColorList[paletteData[0]&0b00111111],
		ColorList[paletteData[1]&0b00111111], ColorList[paletteData[2]&0b00111111]}
}

func GetBackgroundPalette(paletteID byte) (color.RGBA, color.RGBA, color.RGBA, color.RGBA) {
	addr := 0x3F01 + 4*uint16(paletteID)
	paletteData := PPURAM[addr : addr+3]
	return ColorList[PPURAM[0x3F00]&0b00111111], ColorList[paletteData[0]&0b00111111],
		ColorList[paletteData[1]&0b00111111], ColorList[paletteData[2]&0b00111111]
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

func SetMemory(start uint16, data []byte) {
	for i := range data {
		PPURAM[start+uint16(i)] = data[i]
	}
	fmt.Printf("Length of data: %04X\n", len(data))
	fmt.Printf("%04X\n", start)
}

func SetBit(b byte, index int) byte {
	return b | (1 << index)
}

func ClearBit(b byte, index int) byte {
	return b & ^(1 << index)
}

func ToggleBit(b byte, index int) byte {
	return b ^ (1 << index)
}
