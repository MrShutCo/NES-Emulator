package cpu

import (
	"6502/ppu"
	"fmt"
	"log"
	"os"
)

var RAM [0xffff + 1]byte
var X byte
var Y byte
var AC byte
var SR byte
var SP byte
var PC uint16
var Cycles uint64

type CPU6502 struct {
	ram    []byte
	x      byte
	y      byte
	ac     byte
	sr     byte
	sp     byte
	pc     byte
	cycles uint16

	instructions map[byte]Instruction
}

// Useful memory pointers

const STACK uint16 = 0x0100 // Is the end of stack
const SR_RESET = 0b0011_0000
const NMI_VECTOR uint16 = 0xFFFA
const RES_VECTOR uint16 = 0xFFFC
const IRQ_VECTOR uint16 = 0xFFFE

var FuncMap map[byte]func()
var Instructions map[byte]Instruction

var output string

type Instruction struct {
	Name           string
	AddressingMode string
	Cycles         uint8
	Execute        func()
}

func (i Instruction) String() string {
	return i.Name
}

func newInst(opcode byte, name, mode string, cycles uint8) {
	Instructions[opcode] = Instruction{Name: name, AddressingMode: mode, Cycles: cycles}
}

func NewInst(opcode byte, name, mode string, cycles uint8, execute func()) {
	Instructions[opcode] = Instruction{Name: name, AddressingMode: mode, Cycles: cycles, Execute: execute}
}

func Execute() string {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Last output: %s", output)
			log.Printf("Address: %04X\n", PC)
			log.Panicf("Instruction: %02X", RAM[PC])
			log.Println(err)
			panic(err)
		}
	}()

	instruct := RAM[PC]

	start := fmt.Sprintf("%04X  %02X %02X %02x  ", PC, RAM[PC], RAM[PC+1], RAM[PC+2])
	middle := Instructions[RAM[PC]].String()
	regData := fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X SP:%02X", AC, X, Y, SR, SP)

	cycleData := fmt.Sprintf("CYC:%d", Cycles)
	if FuncMap[instruct] != nil {
		FuncMap[instruct]()
	} else {
		Instructions[instruct].Execute()
	}

	a := start + middle + " " + output
	output = ""
	Cycles += uint64(Instructions[instruct].Cycles)
	return fmt.Sprintf("%-47v %v %v", a, regData, cycleData)
}

func Start() {
	PC = GetWordAt(RES_VECTOR)
	RAM[0x2002] = 0b1100_0000
	SR = 0x24
	SP = 0xFD
	Cycles = 7
}

func SetRAM(addr uint16, data byte) {
	RAM[addr] = data
	//fmt.Printf("ADDR: %04X\n", addr)
	switch addr {
	case 0x2000:
		fallthrough
	case 0x2006:
		fallthrough
	case 0x2007:
		ppu.DataStruct.WriteBus(addr, data)
	case 0x4014:
		page := uint16(data) << 8
		var arr [0x100]byte
		copy(arr[:], RAM[page:page+0x100])
		ppu.DataStruct.OAMDMA(arr)
	}
}

func GetRAM(addr uint16) byte {
	switch addr {
	case 0x2002:
		return ppu.DataStruct.ReadBus(addr)
	}
	return RAM[addr]
}

func Reset() {
	Instructions = map[byte]Instruction{}
	X, Y, PC, AC, SP, SR = 0, 0, 0, 0, 0xFF, 0b00110000
	for x := range RAM {
		RAM[x] = 0x00
	}
}

func Load(file string) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("%s", err)
	}
	defer f.Close()

	buffer := make([]byte, 40976)
	//buffer := make([]byte, 1024*64)

	n, err := f.Read(buffer)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("Read %d bytes\n", n)

	// Parse file
	header := buffer[0x0:0x10]

	if string(header[0x0:0x3]) != "NES" {
		fmt.Println("ERROR! File format is not NES")
		return
	}

	PRGROMSize := header[0x4] // In 16KB units
	CHRROMSize := header[0x5] // In 8KB units
	//Flags1 := buffer[0x6]

	// Copy PRGROM
	ROMSize := 16 * 1024 * int(PRGROMSize)
	SetRam(0x8000, buffer[0x10:0x10+ROMSize])
	if n != 40976 {
		SetRam(0xC000, buffer[0x10:0x10+ROMSize])
		fmt.Println("Copying first half of ROM into second half")
	}

	// Copy CHRROM
	startOfCHRROM := 0x10 + ROMSize
	CHRSize := 8 * 1024 * int(CHRROMSize)
	ppu.SetMemory(ppu.PATTERN_TABLE_0, buffer[startOfCHRROM:startOfCHRROM+CHRSize])

	fmt.Printf("PRGROM Size: %04X\nCHRROM Size: %04X\n", ROMSize, CHRSize)
	fmt.Printf("PRGROM copied from [0x%04X,0x%04X] to [0x%04X, 0x%04X] in CPU\n", 0x10, 0x10+ROMSize, 0xC000, 0xC000+ROMSize)
	fmt.Printf("CHRROM copied from [0x%04X,0x%04X] to [0x%04X, 0x%04X] in PPU\n", startOfCHRROM, startOfCHRROM+CHRSize, ppu.PATTERN_TABLE_0, ppu.PATTERN_TABLE_0+CHRSize)

	// HUH??? SetRAM is broken TODO
	for x := 0; x < 0x7FFF; x++ {
		RAM[x] = 0x00
	}
}

func LoadMaps() {
	FuncMap = map[byte]func(){}
	Instructions = map[byte]Instruction{}
	STA()
	STX()
	STY()
	LoadInstructions()
	JMP()
	Stack()
	Flag()
	Branch()
	Other()
	Bitwise()
	Math()
	Compare()
}

func SetRam(start uint16, data []byte) {
	for i := range data {
		RAM[start+uint16(i)] = data[i]
	}
}

func SetX(val byte) {
	setZeroFlag(val == 0x0)
	setNegativeFlag(val >= 0x80)
	X = val
}

func SetY(val byte) {
	setZeroFlag(val == 0x0)
	setNegativeFlag(val >= 0x80)
	Y = val
}

func SetAC(val byte) {
	setZeroFlag(val == 0x0)
	setNegativeFlag(val >= 0x80)
	AC = val
}
