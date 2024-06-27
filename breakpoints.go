package main

type Breakpoint struct {
	Address      uint16
	ClearWhenHit bool

	MarkedForRemoval bool
}

type Disassembly struct {
	Instructions []DisInstruction
}

type DisInstruction struct {
	Command string
	Args    string
	Address uint16
}
