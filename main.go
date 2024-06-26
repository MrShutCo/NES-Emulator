package main

import (
	"6502/cpu"
	"6502/ppu"
	"6502/util"
	"fmt"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var nes NES

const RIGHT = 0b10000000
const LEFT = 0b01000000
const DOWN = 0b00100000
const UP = 0b00010000
const START = 0b00001000
const SELECT = 0b00000100
const BUTTON_B = 0b00000010
const BUTTON_A = 0b00000001

var IsPaused = false

type Game struct {
	Labels                    util.LabelSet
	BreakPointMemoryAddresses []*Breakpoint
}

var keymap map[ebiten.Key]byte

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {

	for key, val := range keymap {
		if inpututil.IsKeyJustPressed(key) {
			cpu.ButtonStatus |= val
		} else if inpututil.IsKeyJustReleased(key) {
			cpu.ButtonStatus &= ^val
		}
	}

	// Pause
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		IsPaused = !IsPaused
	}

	// Step one instruction
	if IsPaused && inpututil.IsKeyJustPressed(ebiten.KeyPeriod) {
		nes.Simulate()
		if cpu.Cycles >= 29780 {
			cpu.Cycles -= 29780
		}
	}

	// Step over instruction, insert breakpoint at the next instruction and ignore the rest
	if IsPaused && inpututil.IsKeyJustPressed(ebiten.KeyComma) {
		g.BreakPointMemoryAddresses = append(g.BreakPointMemoryAddresses, &Breakpoint{
			Address:      cpu.PC + uint16(cpu.GetInstructionAt(cpu.PC).ByteCount),
			ClearWhenHit: true,
		})
		IsPaused = false
	}

	if IsPaused {
		return nil
	}

	for cpu.Cycles < 29780 {
		hasHitBreakPoint := g.CheckBreakpoints()
		if hasHitBreakPoint { // break out so we can continue on debugging
			return nil
		}
		nes.Simulate()
	}

	cpu.Cycles -= 29780

	return nil
}

func (g *Game) CheckBreakpoints() bool {
	// Break at any important points
	i := 0
	hasHitBreakPoint := false
	for _, breakpoint := range g.BreakPointMemoryAddresses {
		if breakpoint != nil && cpu.PC == breakpoint.Address {
			IsPaused = true
			hasHitBreakPoint = true
			if !breakpoint.ClearWhenHit {
				g.BreakPointMemoryAddresses[i] = breakpoint
				i++
			}
		} else {
			g.BreakPointMemoryAddresses[i] = breakpoint
			i++
		}
	}
	for j := i; j < len(g.BreakPointMemoryAddresses); j++ {
		g.BreakPointMemoryAddresses[j] = nil
	}
	g.BreakPointMemoryAddresses = g.BreakPointMemoryAddresses[:i]
	return hasHitBreakPoint
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {

	screen.DrawImage(ppu.Image, &ebiten.DrawImageOptions{
		GeoM: ebiten.ScaleGeo(2, 2),
	})

	screen.DrawImage(ppu.Background, &ebiten.DrawImageOptions{
		GeoM: ebiten.ScaleGeo(2, 2),
	})

	nes.PPU.DrawSprites2(screen)
	ppu.DrawPalettes(screen, 32, 600)

	g.DebugInfo(screen)

}

func (g *Game) DebugInfo(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("IsPaused: %t", IsPaused), 520, 80)
	if IsPaused {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PC: %04X", cpu.PC), 520, 100)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("A: %02X", cpu.AC), 520, 110)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X: %02X", cpu.X), 520, 120)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Y: %02X", cpu.Y), 520, 130)

		pc := cpu.PC
		y := 150
		for i := 0; i < 15; i++ {
			currentInstruction := cpu.Instructions[cpu.RAM[pc]]
			s := FormatInstructionData(currentInstruction, pc)
			if label, ok := g.Labels.GetLabelAt(pc); ok {
				ebitenutil.DebugPrintAt(screen, label.Name, 520, y)
				y += 13
			}
			if currentInstruction.Name == "JSR" || currentInstruction.Name == "JMP" {
				if label, ok := g.Labels.GetLabelAt(cpu.GetWordAt(pc + 1)); ok {
					s = label.Name
				}
			}

			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  %04X %s %s", pc, currentInstruction.String(), s), 520, y)
			if currentInstruction.Name == "RTS" {
				y += 13
			}
			pc += uint16(currentInstruction.ByteCount)
			y += 13
		}

	}
}

func FormatInstructionData(instruct cpu.Instruction, pc uint16) string {
	switch instruct.AddressingMode {
	case "immediate":
		return fmt.Sprintf("#$%02X", cpu.RAM[pc+1])
	case "zeropage":
		return fmt.Sprintf("$%02X", cpu.RAM[pc+1])
	case "zeropage,X":
		return fmt.Sprintf("$%02X,X", cpu.RAM[pc+1])
	case "absolute":
		return fmt.Sprintf("$%04X", cpu.GetWordAt(pc+1))
	case "absolute,X":
		return fmt.Sprintf("$%04X,X", cpu.GetWordAt(pc+1))
	case "absolute,Y":
		return fmt.Sprintf("$%04X,Y", cpu.GetWordAt(pc+1))
	case "indirect":
		return fmt.Sprintf("($%04X)", cpu.GetWordAt(pc+1))
	case "(indirect,X)":
		return fmt.Sprintf("($%02X,X)", cpu.RAM[pc+1])
	case "(indirect),Y":
		return fmt.Sprintf("($%02X),Y", cpu.RAM[pc+1])
	}
	return ""
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 768
}

func main() {
	NESGame()
}

func NESGame() {
	f, _ := os.Create(".cpuprofile.pprof")
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	//labels := util.ReadlabelFile("01.basics.mlb")
	//fmt.Println(labels)

	game := &Game{
		//Labels:                    labels,
		//BreakPointMemoryAddresses: make([]*Breakpoint, 0),
	}
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("NES Emulator")
	ebiten.SetMaxTPS(60)
	cpu.Reset()
	cpu.LoadMaps()
	//cpu.Load("nes-te/st-roms/tutor/tutor.nes")
	//cpu.Load("nes-test-roms/sprite_hit_tests_2005.10.05/01.basics.nes")

	//cpu.Load("nes_test/smb.nes")
	cpu.Load("nes_test/donkeykong.nes")
	//cpu.Load("../nes-test-roms/cpu_dummy_reads/vbl_nmi_timing/7.nmi_timing.nes")
	cpu.Start()
	ppu.DataStruct = ppu.NewPPU()
	ppu.DataStruct.PreloadPalleteTable(ppu.PATTERN_TABLE_0)
	ppu.DataStruct.PreloadPalleteTable(ppu.PATTERN_TABLE_1)
	ppu.DataStruct.LoadPaletteV2(ppu.PATTERN_TABLE_0)
	ppu.DataStruct.LoadPaletteV2(ppu.PATTERN_TABLE_1)
	cpu.PC = cpu.GetWordAt(cpu.RES_VECTOR)
	Image, _ := ebiten.NewImage(512, 512, ebiten.FilterDefault)
	Background, _ := ebiten.NewImage(256, 240, ebiten.FilterDefault)
	ppu.Image = Image
	ppu.Background = Background
	nes = NES{
		PPU: ppu.DataStruct,
	}

	var err error
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		panic(err)
	}
	ppu.Font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    10,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	// Wait for user to exit, then dump logs
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		os.Truncate("output.log", 0)
		f, err := os.OpenFile("output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		f.Write([]byte(nes.stdout))
	}()

	// Setup input
	keymap = map[ebiten.Key]byte{
		ebiten.KeyA:         LEFT,
		ebiten.KeyD:         RIGHT,
		ebiten.KeyS:         DOWN,
		ebiten.KeyW:         UP,
		ebiten.KeySpace:     BUTTON_A,
		ebiten.KeyE:         BUTTON_B,
		ebiten.KeyEnter:     START,
		ebiten.KeyBackspace: SELECT,
	}

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func ExecuteNESTest() {
	cpu.Reset()
	cpu.LoadMaps()
	//cpu.Load("nes_test/colour_test.nes")

	cpu.Load("nes_test/donkeykong.nes")
	cpu.Start()

	ppu.DataStruct = ppu.NewPPU()
	cpu.PC = cpu.GetWordAt(cpu.RES_VECTOR)

	os.Truncate("cpu_dummy_reads.log", 0)
	f, err := os.OpenFile("cpu_dummy_reads.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for i := 0; i < 10000; i++ {
		start := time.Now()
		for cpu.Cycles < 29780 {
			output := cpu.Execute()
			f.Write([]byte(output + "\n"))
			fmt.Print(output)
		}
		end := time.Now()
		fmt.Printf("Time difference: %s\n", end.Sub(start).String())
		cpu.Cycles -= 29780

		//fmt.Println(cpu.RAM[0x2000:0x2100])
	}
}
