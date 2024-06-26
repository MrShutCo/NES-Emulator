package ppu

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

var Image *ebiten.Image
var Font font.Face
var Background *ebiten.Image

func (p *PPU) DrawSprites2(background *ebiten.Image) {
	for i := 0x100 - 4; i >= 0; i -= 4 {
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
		p.CheckSprite0Hit(i, data, posX, posY)

		paletteID := tileAttr & 0x03
		c1, c2, c3, c4 := GetSpritePalette(paletteID)

		pixels := make([]byte, 64*4)
		for j := 0; j < 256; j += 4 {
			col := chooseColor(c1, c2, c3, c4, data[j/4])
			pixels[j] = col.R
			pixels[j+1] = col.G
			pixels[j+2] = col.B
			pixels[j+3] = col.A
		}
		p.sprites[i/4].ReplacePixels(pixels)

		background.DrawImage(p.sprites[i/4], op)
	}
}

func (p *PPU) CheckSprite0Hit(i int, data []byte, posX, posY int) {
	// Sprite 0 hit. See https://www.nesdev.org/wiki/PPU_OAM#Sprite_zero_hits
	// TODO: need some extra cases for when it does and doesnt happen
	if i == 0 && !p.hasSprite0ThisFrame {
		for x := 0; x < 8; x++ {
			for y := 0; y < 8; y++ {
				pX, pY := posX+x, posY+y
				//fmt.Printf("%d,%d\n", pX, pY)
				if data[x+y*8] == 0x00 && Background.At(pX, pY) == ColorList[PPURAM[0x3F00]&0b00111111] { // Transparent pixel
					if pX == 255 { // Extra conditions
						continue
					}
					// Set sprite 0 hit
					_PPUSTATUS = SetBit(_PPUSTATUS, 6)
					p.hasSprite0ThisFrame = true
				}
			}
		}
	}
}

func DrawPalettes(background *ebiten.Image, startX, startY float64) {
	for i := 0; i <= 3; i++ {
		palette := GetSpritePalettePacked(byte(i))
		for x := range palette {
			drawX := startX + float64(x)*32
			drawY := startY + float64(i)*32
			DrawSolidColour(background, palette[x], 32, drawX, drawY)
			//text.Draw(background, fmt.Sprintf("%x", PPURAM[0x3F11+4*i+x]), Font, int(drawX)+16, int(drawY)+16, color.White)
		}
		//palette = GetBackgroundPalette(byte(i))
		//for x := range palette {
		//	DrawSolidColour(background, palette[x], 32, startX+float64(x*32)+128, startY+float64(i)*32)
		//}
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

func chooseColor(c1, c2, c3, c4 color.RGBA, i byte) color.RGBA {
	if i == 0 {
		return c1
	} else if i == 1 {
		return c2
	} else if i == 2 {
		return c3
	} else if i == 3 {
		return c4
	}
	return c1
}

func (p *PPU) DrawDebug() {
	for x := 0; x < 32; x++ {
		for y := 0; y < 30; y++ {
			text.Draw(Image, fmt.Sprintf("%x", GetBackgroundPaletteID(y*32+x)), Font, x*8, y*8, color.White)
		}
	}
	/*for x := 0; x <= 0x0F; x++ {
			for y := 0; y <= 0x03; y++ {
	<<<<<<< HEAD
				c := ppu.ColorMap[byte(y*0x10+x)]
	=======
				c := ppu.ColorList[byte(y*0x10+x)]
	>>>>>>> lots-of-changes
				ppu.DrawSolidColour(screen, c, 32, float64(x)*32, 600+float64(y)*32)
			}

		}*/
}

func (p *PPU) DrawBackgroundRow(startY int) {
	for i := startY * 32; i < startY*32+32; i++ {
		tileIndex := int(PPURAM[p.nametable+uint16(i)])
		index := GetBackgroundPaletteID(i)
		c1, c2, c3, c4 := GetBackgroundPalette(index)

		// Only do update if the index AND palette have changed
		if p.cache[i].NametableIndex == tileIndex && c1 == p.cache[i].Color1 && c2 == p.cache[i].Color2 && c3 == p.cache[i].Color3 && c4 == p.cache[i].Color4 {
			continue
		}

		tileX := i % 32
		tileY := i / 32
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tileX*8), float64(tileY*8))

		var data []byte
		if p.backgroundUsesTable0 {
			data = p.pattern0[tileIndex*64 : tileIndex*64+64]
		} else {
			data = p.pattern1[tileIndex*64 : tileIndex*64+64]
		}

		colourBytes := make([]byte, 256)
		for j := 0; j < 256; j += 4 {
			col := chooseColor(c1, c2, c3, c4, data[j/4])
			colourBytes[j], colourBytes[j+1], colourBytes[j+2], colourBytes[j+3] = col.R, col.G, col.B, col.A
		}

		p.cache[i] = TileCache{
			NametableIndex: tileIndex, Color1: c1, Color2: c2, Color3: c3, Color4: c4,
		}
		p.tiles[i].ReplacePixels(colourBytes)

		Background.DrawImage(p.tiles[i], op)
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

/*
// TODO: this should slowly draw image instead of all at once
func (p *PPU) DrawBackground2(startPosX uint16) {
	for i := 0; i < 0x3c0; i++ {
		tileIndex := int(PPURAM[p.nametable+uint16(i)])
		index := GetBackgroundPaletteID(i)
		c1,c2,c3,c4 := GetBackgroundPalette(index)

		// Only do update if the index AND palette have changed
		if p.cache[i].NametableIndex == tileIndex && equalPalette(p.cache[i].Palette, palette) {
			continue
		}
		p.cache[i] = TileCache{
			NametableIndex: tileIndex,
			Palette:        index,
		}

		tileX := i % 32
		tileY := i / 32
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tileX*8), float64(tileY*8))

		//sx := (tileIndex % 32) * 8
		//sy := (tileIndex / 32) * 8

		//img := image.NewPaletted(image.Rect(int(sx), int(sy), int(sx)+8, int(sy)+8), palette)

		data := p.pattern1[tileIndex*64 : tileIndex*64+64]

		//for j := 0; j < 64; j++ {
		//	img.SetColorIndex((j%8)+sx, (j/8)+sy, data[j])
		//}

		//imgio, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
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
//}

func (p *PPU) GetAttributeIndex(tileX, tileY byte) uint16 {
	addr := (0x3C0 + p.nametable)
	return addr + uint16(tileY)*8 + uint16(tileX)
}
