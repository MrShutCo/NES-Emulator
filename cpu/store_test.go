package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_STX_ZeroPage(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	STX()

	X = 0x11
	RAM[0], RAM[1] = 0x86, 0xFA
	Execute()
	assert.Equal(t, uint16(0x2), PC)
	assert.Equal(t, uint8(0x11), RAM[0x00FA])
}

func Test_STX_ZeroPageY(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	STX()

	X = 0x11
	Y = 0x5
	RAM[0], RAM[1] = 0x96, 0xFA
	Execute()
	assert.Equal(t, uint16(0x2), PC)
	assert.Equal(t, uint8(0x11), RAM[0x00FF])
}

func Test_STX_Absolute(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	STX()

	X = 0x22
	RAM[0], RAM[1], RAM[2] = 0x8E, 0x80, 0x55
	Execute()
	assert.Equal(t, uint16(0x3), PC)
	assert.Equal(t, uint8(0x22), RAM[0x5580])
}
