package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SEI(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	Flag()

	SR = 0b0011_0000
	RAM[0] = 0x78
	Execute()
	assert.Equal(t, byte(0b0011_0100), SR)
	assert.Equal(t, uint16(0x1), PC)
}

func Test_SEC(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	Flag()

	SR = 0b0011_0010
	RAM[0] = 0x38
	Execute()
	assert.Equal(t, byte(0b0011_0011), SR)
	assert.Equal(t, uint16(0x1), PC)
}

func Test_SED(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	Flag()

	SR = 0b0011_0001
	RAM[0] = 0xF8
	Execute()
	assert.Equal(t, byte(0b0011_1001), SR)
	assert.Equal(t, uint16(0x1), PC)
}
