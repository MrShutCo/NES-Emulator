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

func (ls LabelSet) GetLabelForRam(addr uint16) (Label, bool) {
	for _, l := range ls.Labels {
		if l.CPUAddress == addr && (l.Type == "NesInternalRam" || l.Type == "NesMemory") {
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
		lineData := strings.Split(line, ":")
		if len(lineData) < 2 {
			continue
		}
		cpuAddr, _ := strconv.ParseUint(lineData[1], 16, 32)
		l := Label{
			Name: lineData[2],
		}

		if strings.HasPrefix(lineData[0], "NesPrgRom") {
			l.RomAddr = uint16(cpuAddr)
			l.CPUAddress = uint16(cpuAddr) + 0xC000
			l.Type = "NesPrgRom"
		}
		if strings.HasPrefix(lineData[0], "NesMemory") {
			l.CPUAddress = uint16(cpuAddr)
			l.Type = "NesMemory"
		}
		if strings.HasPrefix(lineData[0], "NesInternalRam") {
			l.CPUAddress = uint16(cpuAddr)
			l.RomAddr = uint16(cpuAddr)
			l.Type = "NesInternalRam"
		}

		labels.Labels = append(labels.Labels, l)
	}
	return labels
}
