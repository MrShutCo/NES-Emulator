package util

import "fmt"

func PrintPage(data []byte, page byte) {
	start := uint16(page) << 8
	for y := uint16(0); y < 0x10; y++ {
		fmt.Printf("%04X:  ", start+y*0x10)
		for x := uint16(0); x < 0x10; x++ {
			fmt.Printf("%02X,", data[start+y*0x10+x])
		}
		fmt.Printf("\n")
	}
}
