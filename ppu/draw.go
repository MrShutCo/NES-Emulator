package ppu

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

var Image *ebiten.Image

// TODO: this should slowly draw image instead of all at once
func DrawImage(screen *ebiten.Image) {
	for y := byte(0); y < 30; y++ {
		for x := byte(0); x < 32; x++ {
			val := GetCell(x, y)
			DrawCell(val, x*8, y*8)
			//fmt.Println(val)
		}
	}
	screen.DrawImage(Image, nil)
}

func DrawCell(cell, startX, startY byte) {
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			Image.Set(int(startX)+x, int(startY)+y, color.RGBA{R: cell, G: 0, B: 0, A: 255})
			//fmt.Printf("%d,", cell)
		}
	}
}
