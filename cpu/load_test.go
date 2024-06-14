package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LDA_Immeadiate(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	LDA()

	RAM[0], RAM[1] = 0xA9, 0xFA
	Execute()
	assert.Equal(t, uint16(0x2), PC)
	assert.Equal(t, uint8(0xFA), AC)
}

func Test_LDA_ZeroPage(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	LDA()

	RAM[0x00E9] = 0xAA
	RAM[0x1000], RAM[0x1001] = 0xA5, 0xE9
	PC = 0x1000
	Execute()
	assert.Equal(t, uint16(0x1002), PC)
	assert.Equal(t, uint8(0xAA), AC)
}

func Test_LDA_ZeroPageX(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	LDA()

	RAM[0x00E9] = 0xAA
	RAM[0x00E9+4] = 0x66
	RAM[0x1000], RAM[0x1001], RAM[0x1002] = 0xB5, 0xE9, 0x04
	PC = 0x1000
	Execute()
	assert.Equal(t, uint16(0x1002), PC)
	assert.Equal(t, uint8(0x66), AC)
}
