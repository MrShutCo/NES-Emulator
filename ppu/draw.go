package ppu

import (
	"image"

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
		op.GeoM.Translate(float64(posX), float64(posY))
		op.GeoM.Scale(2, 2)

		data := p.pattern0[tileIndex*64 : tileIndex*64+64]

		paletteID := tileAttr & 0b0000_0011
		palette := GetSpritePalette(paletteID)
		img := image.NewPaletted(image.Rect(int(posX), int(posY), int(posX)+8, int(posY)+8), palette)

		for j := 0; j < 64; j++ {
			img.SetColorIndex((j%8)+posX, (j/8)+posY, data[j])
		}

		imagio, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

		background.DrawImage(imagio, op)
	}
}

// TODO: this should slowly draw image instead of all at once
// DEPRECATED
func (p *PPU) DrawBackground(startPosX uint16) {
	for i := 0; i < 0x3c0; i++ {
		tileIndex := PPURAM[p.nametable+uint16(i)]
		tileX := i % 32
		tileY := i / 32
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tileX*8), float64(tileY*8))

		sx := (int(tileIndex) % 16) * 8
		sy := (int(tileIndex) / 16) * 8
		Image.DrawImage(p.patternTable1SpriteSheet.SubImage(image.Rect(sx, sy, sx+8, sy+8)).(*ebiten.Image), op)
		//ShowTile(tileData, x*8, y*8)
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
	// Cache any tiles for this draw cycle
	cache := map[uint16]*ebiten.Image{}
	for i := 0; i < 0x3c0; i++ {
		tileIndex := int(PPURAM[p.nametable+uint16(i)])
		palette, index := GetBackgroundPalette(i)

		// Only do update if the index AND palette have changed
		if p.cache[i].NametableIndex == byte(tileIndex) && p.cache[i].Palette == index {
			continue
		}
		p.cache[i] = struct {
			NametableIndex byte
			Palette        byte
		}{NametableIndex: byte(tileIndex), Palette: index}

		tileX := i % 32
		tileY := i / 32
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tileX*8), float64(tileY*8))

		if cache[uint16(tileIndex)] == nil {
			sx := (tileIndex % 16) * 8
			sy := (tileIndex / 16) * 8

			img := image.NewPaletted(image.Rect(int(sx), int(sy), int(sx)+8, int(sy)+8), palette)

			data := p.pattern1[tileIndex*64 : tileIndex*64+64]

			for j := 0; j < 64; j++ {
				img.SetColorIndex((j%8)+sx, (j/8)+sy, data[j])
			}

			imgio, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
			cache[uint16(tileIndex)] = imgio
		}

		Image.DrawImage(cache[uint16(tileIndex)], op)
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

func (p *PPU) GetAttributeIndex(tileX, tileY byte) uint16 {
	addr := (0x3C0 + p.nametable)
	return addr + uint16(tileY)*8 + uint16(tileX)
}

func DrawDebug(screen *ebiten.Image) {
	//t := fmt.Sprintf("PPUADDR: 0x%04X\n", _PPUADDR)
	//text.Draw(screen, t, Font, 700, 40, color.White)
}
