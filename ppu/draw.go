package ppu

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

var Image *ebiten.Image

// TODO: this should slowly draw image instead of all at once
func DrawImage(screen *ebiten.Image, bank uint16, startPosX uint16) {
	s := 0
	for y := byte(0); y < 16; y++ {
		for x := byte(0); x < 16; x++ {
			val := GetPatternTable(int(s), bank)
			ShowTile(val, int(x)*8+int(startPosX), int(y)*8)
			//DrawCell(val, x*8, y*8)
			//fmt.Println(val)
			s += 16
		}
	}
	screen.DrawImage(Image, nil)
}

func DrawNameTable0(screen *ebiten.Image) {
	for i := 0; i < 0x3c0; i++ {
		val := GetPatternTable(int(PPURAM[i+NAMETABLE_0]), PATTERN_TABLE_0)
		x := i % 32
		y := i / 32
		ShowTile(val, x, y)
	}
}

func DrawCell(cell byte, startX, startY int) {
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			Image.Set(int(startX)+x, int(startY)+y, color.RGBA{R: cell, G: 0, B: 0, A: 255})
			//fmt.Printf("%d,", cell)
		}
	}
}
