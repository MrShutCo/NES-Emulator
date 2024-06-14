package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareOutputs(t *testing.T) {
	compare(t)
}

func compare(t *testing.T) {
	ourOutput, _ := os.OpenFile("access.log", os.O_RDONLY, 0644)
	goodOutput, _ := os.OpenFile("nes_test/nestest.log", os.O_RDONLY, 0644)
	defer ourOutput.Close()
	defer goodOutput.Close()

	ourScanner := bufio.NewScanner(ourOutput)
	goodScanner := bufio.NewScanner(goodOutput)

	lastO := ""
	lastG := ""

	for ourScanner.Scan() {
		goodScanner.Scan()
		o := ourScanner.Text()
		g := goodScanner.Text()
		if !(checkAddr(t, o, g) && checkParsedOpcode(t, o, g) && checkRegisters(t, o, g) && checkCycles(t, o, g)) {
			fmt.Printf("expected:\n%s\n%s\nactual:\n%s\n%s\n", lastG, g, lastO, o)
			return
		}
		lastO = o
		lastG = g
	}
}

func checkAddr(t *testing.T, ours, good string) bool {
	return assert.Equal(t, good[0:4], ours[0:4])

}

func checkParsedOpcode(t *testing.T, ours, good string) bool {
	return assert.Equal(t, good[16:35], ours[16:35])
}

func checkRegisters(t *testing.T, ours, good string) bool {
	return assert.Equal(t, good[48:73], ours[48:73])
}

func checkCycles(t *testing.T, ours, good string) bool {
	return assert.Equal(t, good[86:], ours[74:])
}
