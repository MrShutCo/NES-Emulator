package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BCC(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	Branch()

	// Branch not taken
	setCarryFlag(true)
	SetRam(0, []byte{0x90, 0x04, 0x90, 0x05})
	SetRam(0x7, []byte{0x90, 0xFE}) // Move back 2
	Execute()
	assert.Equal(t, uint16(0x2), PC)

	// Branch taken, move ahead 5 spaces from 0x2
	setCarryFlag(false)
	Execute()
	assert.Equal(t, uint16(0x7), PC)

	// Move backwards 2 spaces to 0x5
	Execute()
	assert.Equal(t, uint16(0x5), PC)
}

func Test_BCS(t *testing.T) {
	Reset()
	FuncMap = map[byte]func(){}
	Branch()

	// Branch not taken
	setCarryFlag(false)
	SetRam(0, []byte{0xB0, 0x04, 0xB0, 0x05})
	SetRam(0x7, []byte{0xB0, 0xFE}) // Move back 2
	Execute()
	assert.Equal(t, uint16(0x2), PC)

	// Branch taken, move ahead 5 spaces from 0x2
	setCarryFlag(true)
	Execute()
	assert.Equal(t, uint16(0x7), PC)

	// Move backwards 2 spaces to 0x5
	Execute()
	assert.Equal(t, uint16(0x5), PC)
}
