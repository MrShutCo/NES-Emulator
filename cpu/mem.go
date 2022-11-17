package cpu

import (
	"fmt"
	"os"
)

var RAM [0xffff + 1]byte
var X byte
var Y byte
var AC byte
var SR byte
var SP byte
var PC uint16

// Useful memory pointers

const STACK uint16 = 0x0100 // Is the end of stack
const SR_RESET = 0b0011_0000
const NMI_VECTOR uint16 = 0xFFFA
const RES_VECTOR uint16 = 0xFFFD
const IRQ_VECTOR uint16 = 0xFFFE

var FuncMap map[byte]func()
var Instructions map[byte]Instruction

var output string

type Instruction struct {
	Name           string
	AddressingMode string
	Cycles         uint8
}

func (i Instruction) String() string {
	return i.Name
}

func newInst(opcode byte, name, mode string, cycles uint8) {
	Instructions[opcode] = Instruction{Name: name, AddressingMode: mode, Cycles: cycles}
}

func Execute() {
	instruct := RAM[PC]

	//output := fmt.Sprintf("%04X %02X %02X %02X", PC, instruct)

	start := fmt.Sprintf("%04X %02X %02X %02x  ", PC, RAM[PC], RAM[PC+1], RAM[PC+2])
	middle := Instructions[RAM[PC]].String()

	//fmt.Printf("Executing 0x%x: %s\tat 0x%x\n", instruct, Instructions[instruct].String(), PC)

	end := fmt.Sprintf("A:%02X X:%02X Y %02X P:%02X SP:%02X", AC, X, Y, SR, SP)

	if FuncMap[instruct] == nil {
		//printStack()
		fmt.Printf("PC: %04X\n", PC)
		fmt.Printf("Found not implemented instruction: 0x%02X\n", instruct)
	}

	FuncMap[instruct]()

	a := start + middle + " " + output
	fmt.Printf("%-40v %29v", a, end)
	output = ""
}

func Start() {
	PC = GetWordAt(0xFFFC)
	//setNegativeFlag(true)
	//setOverflowFlag(true) // So BIT can pass
	SR = 0x24
	SP = 0xFD
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

	//buffer := make([]byte, 40976)
	buffer := make([]byte, 1024*64)

	n, err := f.Read(buffer)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("Read %d bytes\n", n)
	SetRam(0xC000, buffer[0x10:]...) // Copy everything past the header into ROM
	//SetRam(0x0, buffer[:]...)
	//fmt.Println(RAM[0x8000:0x9000])
	//fmt.Println(RAM[0x8000:0xFFFF])
}

func LoadMaps() {
	FuncMap = map[byte]func(){}
	Instructions = map[byte]Instruction{}
	STA()
	STX()
	STY()
	LDA()
	LDX()
	LDY()
	JMP()
	Stack()
	Flag()
	Branch()
	Other()
	Bitwise()
}

func SetRam(start uint16, data ...byte) {
	for i := range data {
		RAM[start+uint16(i)] = data[i]
	}
}

func SetX(val byte) {
	setZeroFlag(val == 0x0)
	setNegativeFlag(val >= 0b1000_0000)
	X = val
}

func SetY(val byte) {
	setZeroFlag(val == 0x0)
	setNegativeFlag(val >= 0b1000_0000)
	Y = val
}

func SetAC(val byte) {
	setZeroFlag(val == 0x0)
	setNegativeFlag(val >= 0b1000_0000)
	AC = val
}
