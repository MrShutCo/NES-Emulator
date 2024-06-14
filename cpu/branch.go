package cpu

func Branch() {
	// TODO: fix paging for cycles
	newInst(0x90, "BCC", "", 2)
	newInst(0xB0, "BCS", "", 2)
	newInst(0xF0, "BEQ", "", 2)
	newInst(0x30, "BMI", "", 2)
	newInst(0xD0, "BNE", "", 2)
	newInst(0x10, "BPL", "", 2)
	newInst(0x50, "BVC", "", 2)
	newInst(0x70, "BVS", "", 2)
	// BCC
	FuncMap[0x90] = func() {
		val := addsignedByteToUInt(RAM[PC+1], PC)
		////output = fmt.Sprintf("$%04X", val+2)
		if !isCarrySet() {
			PC = val
			Cycles++
		}
		PC += 2
	}
	// BCS
	FuncMap[0xB0] = func() {
		val := addsignedByteToUInt(RAM[PC+1], PC)
		////output = fmt.Sprintf("$%04X", val+2)
		if isCarrySet() {
			PC = val
			Cycles++
		}
		PC += 2
	}
	// BEQ
	FuncMap[0xF0] = func() {
		val := addsignedByteToUInt(RAM[PC+1], PC)
		//output = fmt.Sprintf("$%04X", val+2)
		if isZeroSet() {
			PC = val
			Cycles++
		}
		PC += 2
	}
	// BMI
	FuncMap[0x30] = func() {
		val := addsignedByteToUInt(RAM[PC+1], PC)
		//output = fmt.Sprintf("$%04X", val+2)
		if isNegativeSet() {
			PC = val
			Cycles++
		}
		PC += 2
	}
	// BNE
	FuncMap[0xD0] = func() {
		val := addsignedByteToUInt(RAM[PC+1], PC)
		//output = fmt.Sprintf("$%04X", val+2)
		if !isZeroSet() {
			PC = val
			Cycles++
		}
		PC += 2
	}
	// BPL
	FuncMap[0x10] = func() {
		val := addsignedByteToUInt(RAM[PC+1], PC)
		//output = fmt.Sprintf("$%04X", val+2)
		if !isNegativeSet() {
			PC = val
			Cycles++
		}
		PC += 2
	}
	// BRK
	FuncMap[0x00] = func() {
		// TODO
		push(lowByte(PC))
		push(highByte(PC))
		push(SR)
		PC = GetWordAt(IRQ_VECTOR)
		setBreakFlag(true)
		Cycles += 7
		PC++
	}
	// BVC
	FuncMap[0x50] = func() {
		val := addsignedByteToUInt(RAM[PC+1], PC)
		//output = fmt.Sprintf("$%04X", val+2)
		if !isOverflowSet() {
			PC = val
			Cycles++
		}
		PC += 2
	}
	// BVS
	FuncMap[0x70] = func() {
		val := addsignedByteToUInt(RAM[PC+1], PC)
		//output = fmt.Sprintf("$%04X", val+2)
		if isOverflowSet() {
			PC = val
			Cycles++
		}
		PC += 2
	}
}
