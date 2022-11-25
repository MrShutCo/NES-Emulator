package ppu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_get2BitsFromByte(t *testing.T) {
	assert.Equal(t, byte(0), get2BitsFromByte(0b11100100, 0))
	assert.Equal(t, byte(1), get2BitsFromByte(0b11100100, 2))
	assert.Equal(t, byte(2), get2BitsFromByte(0b11100100, 4))
	assert.Equal(t, byte(3), get2BitsFromByte(0b11100100, 6))
}

func TestUnit_getAttrByteFromTileIndex(t *testing.T) {
	assert.Equal(t, byte(0), getAttrByteFromTileIndex(0))
	assert.Equal(t, byte(0), getAttrByteFromTileIndex(3))
	assert.Equal(t, byte(0), getAttrByteFromTileIndex(4))
	assert.Equal(t, byte(0), getAttrByteFromTileIndex(7))

	assert.Equal(t, byte(1), getAttrByteFromTileIndex(8))
	assert.Equal(t, byte(1), getAttrByteFromTileIndex(9))
	assert.Equal(t, byte(1), getAttrByteFromTileIndex(14))
	assert.Equal(t, byte(1), getAttrByteFromTileIndex(15))

	assert.Equal(t, byte(10), getAttrByteFromTileIndex(8))
	assert.Equal(t, byte(10), getAttrByteFromTileIndex(9))
	assert.Equal(t, byte(10), getAttrByteFromTileIndex(14))
	assert.Equal(t, byte(10), getAttrByteFromTileIndex(15))
}
