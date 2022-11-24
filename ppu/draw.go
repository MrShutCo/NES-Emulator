package ppu

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

var Image *ebiten.Image
var Font font.Face

var screenBuffer []byte

// TODO: this should slowly draw image instead of all at once
func (p *PPU) DrawBackground(startPosX uint16) {
	/*s := 0
	for y := byte(0); y < 16; y++ {
		for x := byte(0); x < 16; x++ {
			image := GetImageFromPatternTable(y*16+x, p.backgroundPatternTable)
			ShowTile(image, int(x)*8+int(startPosX), int(y)*8)
			//DrawCell(val, x*8, y*8)
			//fmt.Println(val)
			s += 16
		}
	}*/
	for i := 0; i < 0x3c0; i++ {
		tileIndex := PPURAM[p.nametable+uint16(i)]
		tileData := GetImageFromPatternTable(tileIndex, p.backgroundPatternTable)
		tileX := i % 32
		tileY := i / 32
		for y := 0; y < 8; y++ {
			lower := tileData[y]
			upper := tileData[y+8]
			for x := 0; x < 8; x++ {
				colour := (lower >> (8 - x) & 1)
				colour += 2 * (upper >> (8 - x) & 1)
				// TODO: precalculate the bytes for all tiles, and then just copy them all from another array or something
				// Relative position inside the sprite + How much is added by tileX and tileY
				pos := (y*1024 + x*4 + tileX*32 + tileY*32*64*4)
				//pos := (y*256 + x + tileX*8 + i*64) * 4
				switch colour {
				case 0x0:
					screenBuffer[pos], screenBuffer[pos+1], screenBuffer[pos+2], screenBuffer[pos+3] = 50, 50, 50, 255
				case 0x1:
					screenBuffer[pos], screenBuffer[pos+1], screenBuffer[pos+2], screenBuffer[pos+3] = 100, 100, 100, 255
				case 0x2:
					screenBuffer[pos], screenBuffer[pos+1], screenBuffer[pos+2], screenBuffer[pos+3] = 150, 150, 150, 255
				case 0x3:
					screenBuffer[pos], screenBuffer[pos+1], screenBuffer[pos+2], screenBuffer[pos+3] = 255, 255, 255, 255
				}
				//pos += 4
			}
		}
		//ShowTile(tileData, x*8, y*8)
	}
	Image.ReplacePixels(screenBuffer[:])
	//screen.DrawImage(Image, nil)
}

func writebytes(array []byte, data []byte, pos int) []byte {
	for i := range data {
		array[pos+i] = data[i]
	}
	return array
}

func DrawDebug(screen *ebiten.Image) {
	t := fmt.Sprintf("PPUADDR: 0x%04X\n", _PPUADDR)
	text.Draw(screen, t, Font, 700, 40, color.White)
}

func DrawNameTable0(screen *ebiten.Image) {
	for i := 0; i < 0x3c0; i++ {
		//val := GetPatternTable(int(PPURAM[i+NAMETABLE_0]), PATTERN_TABLE_0)
		//x := i % 32
		//y := i / 32
		//ShowTile(val, x, y)
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
