package main

type Breakpoint struct {
	Address      uint16
	ClearWhenHit bool

	MarkedForRemoval bool
}
