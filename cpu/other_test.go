package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_JMP_Absolute(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	JMP()

	RAM[0], RAM[1], RAM[2] = 0x4C, 0xFA, 0x57
	Execute()
	assert.Equal(t, uint16(0x57FA), PC)
}

func Test_CPX_Immeadiate(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	Compare()

	X = 0x80
	// Set up four different comparisons
	SetRam(0, []byte{0xE0, 0x79, 0xE0, 0x80, 0xE0, 0x81, 0xE0, 0x7f})
	Execute()
	assert.Equal(t, byte(0b0011_0001), SR)
	assert.Equal(t, uint16(0x2), PC)
	SR = SR_RESET

	Execute()
	assert.Equal(t, byte(0b0011_0011), SR)
	assert.Equal(t, uint16(0x4), PC)
	SR = SR_RESET

	Execute()
	assert.Equal(t, byte(0b1011_0000), SR)
	assert.Equal(t, uint16(0x6), PC)
	SR = SR_RESET

	X = 0xFF
	Execute()
	assert.Equal(t, byte(0b1011_0001), SR)
	assert.Equal(t, uint16(0x8), PC)
}

func Test_CPX_Zeropage(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	Compare()

	RAM[0], RAM[1] = 0xE4, 0xAB
	RAM[0xAB] = 0x81 // Make sure it doesnt mistake as 0

	X = 0x80
	Execute()
	assert.Equal(t, byte(0b1011_0000), SR)
	assert.Equal(t, uint16(0x2), PC)
}

func Test_CPX_Absolute(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	Compare()

	SetRam(0, []byte{0xEC, 0xAB, 0x56})
	RAM[0x56AB] = 0x81 // Make sure it doesnt mistake as 0

	X = 0x80
	Execute()
	assert.Equal(t, byte(0b1011_0000), SR)
	assert.Equal(t, uint16(0x3), PC)
}

// TODO: are all these bits wrong????

func Test_BIT_Zeropage(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	Other()

	SetRam(0, []byte{0x24, 0x79})
	RAM[0x79] = 0b0100_0000
	Execute()
	assert.Equal(t, byte(0b01110010), SR)
	assert.Equal(t, uint16(0x2), PC)

	// AND is true, set
	PC, RAM[0x79], SR = 0, 0b11000000, 0b0011_0000
	AC = 0b11000000
	Execute()
	assert.Equal(t, byte(0b11110000), SR)
	assert.Equal(t, uint16(0x2), PC)

}
