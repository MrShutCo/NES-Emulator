package ppu

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/text"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

var Image *ebiten.Image
var Font font.Face

func (p *PPU) DrawSprites2(background *ebiten.Image) {
	for i := 0; i < 0x100; i += 4 {
		posY := int(OAM[i] - 1)
		posX := int(OAM[i+3])
		tileIndex := int(OAM[i+1])
		tileAttr := OAM[i+2]
		op := &ebiten.DrawImageOptions{}

		flipVertical := tileAttr&0x80 == 0x80
		flipHorizontal := tileAttr&0x40 == 0x40

		if flipHorizontal {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(float64(8), 0)
		}
		if flipVertical {
			op.GeoM.Scale(1, -1)
		}

		op.GeoM.Translate(float64(posX), float64(posY))
		op.GeoM.Scale(2, 2)

		data := p.pattern0[tileIndex*64 : tileIndex*64+64]

		paletteID := tileAttr & 0x03
		palette := GetSpritePalette(paletteID)
		img := image.NewPaletted(image.Rect(int(posX), int(posY), int(posX)+8, int(posY)+8), palette)

		for j := 0; j < 64; j++ {
			img.SetColorIndex((j%8)+posX, (j/8)+posY, data[j])
		}

		imagio, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

		background.DrawImage(imagio, op)
	}
}

func DrawPalettes(background *ebiten.Image, startX, startY float64) {
	for i := 0; i <= 3; i++ {
		palette := GetSpritePalette(byte(i))
		for x := range palette {
			drawX := startX + float64(x)*32
			drawY := startY + float64(i)*32
			DrawSolidColour(background, palette[x], 32, drawX, drawY)
			//text.Draw(background, fmt.Sprintf("%x", PPURAM[0x3F11+4*i+x]), Font, int(drawX)+16, int(drawY)+16, color.White)
		}
		palette = GetBackgroundPalette(byte(i))
		for x := range palette {
			DrawSolidColour(background, palette[x], 32, startX+float64(x*32)+128, startY+float64(i)*32)
		}
	}
}

func DrawAttributeTable(background *ebiten.Image, startX, startY float64) {

}

func DrawSolidColour(background *ebiten.Image, color color.Color, size int, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	img, _ := ebiten.NewImage(size, size, ebiten.FilterDefault)
	data := make([]byte, size*size*4)
	for i := 0; i < size*size; i++ {
		r, g, b, a := color.RGBA()
		data[i*4], data[i*4+1], data[i*4+2], data[i*4+3] = byte(r), byte(g), byte(b), byte(a)
	}
	img.ReplacePixels(data)
	background.DrawImage(img, op)
}

func equalPalette(p1, p2 color.Palette) bool {
	if (p1 == nil) != (p2 == nil) {
		return false
	}
	for i := 0; i < 3; i++ {
		if p1[i] != p2[i] {
			return false
		}
	}
	return true
}

func (p *PPU) DrawDebug() {
	for x := 0; x < 32; x++ {
		for y := 0; y < 30; y++ {
			text.Draw(Image, fmt.Sprintf("%x", GetBackgroundPaletteID(y*32+x)), Font, x*8, y*8, color.White)
		}
	}
	/*for x := 0; x <= 0x0F; x++ {
		for y := 0; y <= 0x03; y++ {
			c := ppu.ColorList[byte(y*0x10+x)]
			ppu.DrawSolidColour(screen, c, 32, float64(x)*32, 600+float64(y)*32)
		}

	}*/
}

func (p *PPU) DrawBackgroundRow(startY int) {
	for i := startY * 32; i < startY*32+32; i++ {
		tileIndex := int(PPURAM[p.nametable+uint16(i)])
		index := GetBackgroundPaletteID(i)
		palette := GetBackgroundPalette(index)

		// Only do update if the index AND palette have changed
		if p.cache[i].NametableIndex == tileIndex && equalPalette(p.cache[i].Palette, palette) {
			continue
		}

		tileX := i % 32
		tileY := i / 32
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tileX*8), float64(tileY*8))

		sx := (tileIndex % 32) * 8
		sy := (tileIndex / 32) * 8

		img := image.NewPaletted(image.Rect(int(sx), int(sy), int(sx)+8, int(sy)+8), palette)

		data := p.pattern1[tileIndex*64 : tileIndex*64+64]

		for j := 0; j < 64; j++ {
			img.SetColorIndex((j%8)+sx, (j/8)+sy, data[j])
		}

		imgio, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		p.cache[i] = TileCache{
			NametableIndex: tileIndex, Palette: palette, Tile: imgio,
		}

		Image.DrawImage(imgio, op)
	}
	// PALETTE_0
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(256+64, 0)
	Image.DrawImage(p.patternTable0SpriteSheet, op)

	// PALETTE_1
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(256+64), float64(128+32))
	Image.DrawImage(p.patternTable1SpriteSheet, op)
}

// TODO: this should slowly draw image instead of all at once
func (p *PPU) DrawBackground2(startPosX uint16) {
	for i := 0; i < 0x3c0; i++ {
		tileIndex := int(PPURAM[p.nametable+uint16(i)])
		index := GetBackgroundPaletteID(i)
		palette := GetBackgroundPalette(index)

		// Only do update if the index AND palette have changed
		if p.cache[i].NametableIndex == tileIndex && equalPalette(p.cache[i].Palette, palette) {
			continue
		}
		/*p.cache[i] = TileCache{
			NametableIndex: tileIndex,
			Palette:        index,
		}*/

		tileX := i % 32
		tileY := i / 32
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tileX*8), float64(tileY*8))

		sx := (tileIndex % 32) * 8
		sy := (tileIndex / 32) * 8

		img := image.NewPaletted(image.Rect(int(sx), int(sy), int(sx)+8, int(sy)+8), palette)

		data := p.pattern1[tileIndex*64 : tileIndex*64+64]

		for j := 0; j < 64; j++ {
			img.SetColorIndex((j%8)+sx, (j/8)+sy, data[j])
		}

		imgio, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		p.cache[i] = TileCache{
			NametableIndex: tileIndex, Palette: palette, Tile: imgio,
		}

		Image.DrawImage(imgio, op)
	}
	// PALETTE_0
	/*op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(256+64, 0)
	Image.DrawImage(p.patternTable0SpriteSheet, op)

	// PALETTE_1
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(256+64), float64(128+32))
	Image.DrawImage(p.patternTable1SpriteSheet, op)*/
}

func (p *PPU) GetAttributeIndex(tileX, tileY byte) uint16 {
	addr := (0x3C0 + p.nametable)
	return addr + uint16(tileY)*8 + uint16(tileX)
}
