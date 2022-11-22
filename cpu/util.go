package cpu

import (
	"fmt"
	"reflect"
	"runtime"
)

func bytesToInt16(high, low byte) uint16 {
	return uint16(high)<<8 | uint16(low)
}

func getNextWord() uint16 {
	return GetWordAt(PC + 1)
}

func GetWordAt(addr uint16) uint16 {
	return bytesToInt16(RAM[addr+1], RAM[addr])
}

func addsignedByteToUInt(val byte, larger uint16) uint16 {
	return larger + uint16(int8(val))
}

func lowByte(data uint16) byte {
	return byte(data & 0x00FF)
}

func highByte(data uint16) byte {
	return byte((data & 0xFF00) >> 8)
}

// ===== Getting status bits

func getBit(b byte, bit int8) bool {
	if bit > 7 {
		return false
	}
	return b&(1<<bit)>>bit == 1
}

func isCarrySet() bool     { return getBit(SR, 0) }
func isZeroSet() bool      { return getBit(SR, 1) }
func isInterruptSet() bool { return getBit(SR, 2) }
func isDecimalSet() bool   { return getBit(SR, 3) }
func isBreakSet() bool     { return getBit(SR, 4) }
func isBit5Set() bool      { return getBit(SR, 5) }
func isOverflowSet() bool  { return getBit(SR, 6) }
func isNegativeSet() bool  { return getBit(SR, 7) }

// ======================

func setNegativeFlag(enable bool) {
	if enable {
		SR = SR | 0b1000_0000
	} else {
		SR = SR & 0b0111_1111
	}
}

func setOverflowFlag(enable bool) {
	if enable {
		SR = SR | 0b0100_0000
	} else {
		SR = SR & 0b1011_1111
	}
}

func setEmptyFlag(enable bool) {
	if enable {
		SR = SR | 0b0010_0000
	} else {
		SR = SR & 0b1101_1111
	}
}

func setBreakFlag(enable bool) {
	if enable {
		SR = SR | 0b0001_0000
	} else {
		SR = SR & 0b1110_1111
	}
}

func setDecimalFlag(enable bool) {
	if enable {
		SR = SR | 0b0000_1000
	} else {
		SR = SR & 0b1111_0111
	}
}

func setInterruptFlag(enable bool) {
	if enable {
		SR = SR | 0b0000_0100
	} else {
		SR = SR & 0b1111_1011
	}
}

func setZeroFlag(enable bool) {
	if enable {
		SR = SR | 1<<1
	} else {
		SR = SR & 0b1111_1101
	}
}

func setCarryFlag(enable bool) {
	if enable {
		SR = SR | 1
	} else {
		SR = SR & 0b1111_1110
	}
}

// Addressing mode
func apply(functions []foo) {
	for i := range functions {
		FuncMap[functions[i].opcode] = functions[i].f
	}
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func ac() byte {
	output = "A"
	PC++
	return AC
}

func immed() byte {
	val := GetRAM(PC + 1)
	output = fmt.Sprintf("#$%02X", RAM[PC+1])
	PC += 2
	return val
}

func zeropage() byte {
	val := GetRAM(zeropageAddr())

	PC += 2
	return val
}

func zeropageAddr() uint16 {
	addr := uint16(GetRAM(PC + 1))
	output = fmt.Sprintf("$%02X = %02X", addr, RAM[addr])
	return addr
}

func zeropageX() byte {
	val := GetRAM(zeropageXAddr())
	PC += 2
	return val
}

func zeropageXAddr() uint16 {
	addr := byte(GetRAM(PC+1)) + X
	output = fmt.Sprintf("$%02X,X @ %02X = %02X", GetRAM(PC+1), addr, GetRAM(uint16(addr)))
	return uint16(addr)
}

func zeropageY() byte {
	val := GetRAM(zeropageYAddr() % 256)
	PC += 2
	return val
}

func zeropageYAddr() uint16 {
	addr := byte(GetRAM(PC+1)) + Y
	output = fmt.Sprintf("$%02X,Y @ %02X = %02X", GetRAM(PC+1), addr, GetRAM(uint16(addr)))
	return uint16(addr)
}

func absolute() byte {
	addr := absoluteAddr()
	val := GetRAM(addr)
	PC += 3
	return val
}

func absoluteAddr() uint16 {
	addr := getNextWord()
	output = fmt.Sprintf("$%04X = %02X", addr, RAM[addr])
	return addr
}

func absoluteX() byte {
	addrX := absoluteXAddr()
	addr := getNextWord()
	if highByte(addr) != highByte(addrX) {
		Cycles++
	}
	PC += 3
	//output = fmt.Sprintf("$%04X = %02X", addrX, RAM[addrX])
	val := GetRAM(addrX)
	return val
}

func absoluteXAddr() uint16 {
	addrX := bytesToInt16(GetRAM(PC+2), GetRAM(PC+1)) + uint16(X)
	output = fmt.Sprintf("$%04X,X @ %04X = %02X", getNextWord(), addrX, RAM[addrX])
	return addrX
}

func absoluteY() byte {
	addrY := absoluteYAddr()
	addr := getNextWord()
	if highByte(addr) != highByte(addrY) {
		Cycles++
	}
	val := GetRAM(addrY)
	PC += 3
	return val
}

func absoluteYAddr() uint16 {
	addrY := bytesToInt16(GetRAM(PC+2), GetRAM(PC+1)) + uint16(Y)
	output = fmt.Sprintf("$%04X,Y @ %04X = %02X", getNextWord(), addrY, RAM[addrY])
	return addrY
}

func indirectX16() uint16 {
	zeropageAddr := uint16(GetRAM(PC+1)) + uint16(X)
	low := GetRAM(zeropageAddr % 256)
	hi := GetRAM((zeropageAddr + 1) % 256) // Wraparound
	addr := bytesToInt16(hi, low)
	output = fmt.Sprintf("($%02X,X) @ %02X = %04X = %02X", RAM[PC+1], RAM[PC+1]+X, addr, RAM[addr])
	return addr
}

func indirectX() byte {
	addr := indirectX16()
	PC += 2
	return GetRAM(addr)
}

func indirectYAddr() uint16 {
	low := RAM[RAM[PC+1]]
	hi := RAM[RAM[PC+1]+1]
	addr := bytesToInt16(hi, low)
	addrY := int(addr) + int(Y)
	// If its on a different page, add a cycle
	if PC == 0xDB65 {
		i := 0
		fmt.Println(i)
	}
	if addrY > 0xFFFF {
		Cycles++
	} else if highByte(addr) != highByte(uint16(addrY)) {
		Cycles++
	}

	output = fmt.Sprintf("($%02X),Y = %04X @ %04X = %02X", RAM[PC+1], addr, addrY, RAM[uint16(addrY)])
	return uint16(addrY)
}

func indirectY() byte {
	addr := indirectYAddr()
	PC += 2
	return GetRAM(addr)
}
