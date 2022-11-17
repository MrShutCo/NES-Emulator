package cpu

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