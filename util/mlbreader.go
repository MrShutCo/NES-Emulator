package util

import (
	"os"
	"strconv"
	"strings"
)

type LabelSet struct {
	Labels []Label
}

type Label struct {
	RomAddr    uint16
	CPUAddress uint16
	Name       string
	Comment    string
	Type       string
}

func (ls LabelSet) GetLabelAt(addr uint16) (Label, bool) {
	for _, l := range ls.Labels {
		if l.CPUAddress == addr {
			return l, true
		}
	}
	return Label{}, false
}

func ReadlabelFile(file string) LabelSet {
	data, _ := os.ReadFile(file)
	lines := strings.Split(string(data), "\n")

	labels := LabelSet{Labels: make([]Label, 0)}
	for _, line := range lines {
		if strings.HasPrefix(line, "NesPrgRom") {
			lineData := strings.Split(line, ":")
			cpuAddr, _ := strconv.ParseUint(lineData[1], 16, 32)
			labels.Labels = append(labels.Labels, Label{
				RomAddr:    uint16(cpuAddr),
				CPUAddress: uint16(cpuAddr) + 0xC000,
				Name:       lineData[2],
				Comment:    "",
				Type:       "NesPrgRom",
			})
		}
		if strings.HasPrefix(line, "NesMemory") {

		}
	}
	return labels
}
