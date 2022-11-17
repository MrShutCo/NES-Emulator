package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetStatus(t *testing.T) {
	Reset()
	SR = 0b1000_1001
	assert.Equal(t, getBit(SR, 0), true)
	assert.Equal(t, getBit(SR, 1), false)
	assert.Equal(t, getBit(SR, 2), false)
	assert.Equal(t, getBit(SR, 3), true)
	assert.Equal(t, getBit(SR, 4), false)
	assert.Equal(t, getBit(SR, 5), false)
	assert.Equal(t, getBit(SR, 6), false)
	assert.Equal(t, getBit(SR, 7), true)
}

func Test_LowByte(t *testing.T) {
	assert.Equal(t, byte(0xCD), lowByte(0xABCD))
}

func Test_HighByte(t *testing.T) {
	assert.Equal(t, byte(0xAB), highByte(0xABCD))
}
