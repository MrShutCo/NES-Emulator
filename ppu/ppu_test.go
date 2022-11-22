package ppu

import (
	"fmt"
	"testing"

	"github.com/hajimehoshi/ebiten"
	"github.com/stretchr/testify/assert"
)

func Test_ShowTile(t *testing.T) {
	Image, _ = ebiten.NewImage(255, 255, ebiten.FilterDefault)
	fmt.Println(Image.Size())

	result := ShowTile([]byte{
		0x0, 0xFF, 0x00, 0xFF, 0b00001111, 0, 0, 0,
		0x0, 0x00, 0xFF, 0xFF, 0b00110011, 0, 0, 0,
	}, 0, 0)

	assert.Equal(t, result[0:8], []byte{0, 0, 0, 0, 0, 0, 0, 0})
	assert.Equal(t, result[8:16], []byte{1, 1, 1, 1, 1, 1, 1, 1})
	assert.Equal(t, result[16:24], []byte{2, 2, 2, 2, 2, 2, 2, 2})
	assert.Equal(t, result[24:32], []byte{3, 3, 3, 3, 3, 3, 3, 3})
	assert.Equal(t, result[32:40], []byte{3, 3, 1, 1, 2, 2, 0, 0})
}
