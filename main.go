package main

import (
	"6502/cpu"
	"6502/ppu"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
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

type Game struct{}

var keymap map[ebiten.Key]byte

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {

	for key, val := range keymap {
		if inpututil.IsKeyJustPressed(key) {
			cpu.ButtonStatus |= val
		} else if inpututil.IsKeyJustReleased(key) {
			cpu.ButtonStatus &= ^byte(val)
		}
	}

	//t := time.Now()
	for cpu.Cycles < 29780 {
		nes.Simulate()
	}

	//fmt.Println(time.Since(t).Milliseconds())
	cpu.Cycles -= 29780

	//s := time.Now()
	//
	//a := time.Now()
	//fmt.Printf("Time difference: %s\n", a.Sub(s).String())

	return nil
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

	// Useful for debugging
	//nes.PPU.DrawDebug()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 768
}

func main() {
	NESGame()
	//ExecuteNESTest()
}

func NESGame() {
	f, _ := os.Create(".cpuprofile.pprof")
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	game := &Game{}
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("NES Emulator")
	ebiten.SetMaxTPS(60)
	cpu.Reset()
	cpu.LoadMaps()
	//cpu.Load("nes-te/st-roms/tutor/tutor.nes")
	cpu.Load("nes-test-roms/sprite_hit_tests_2005.10.05/01.basics.nes")
	//cpu.Load("nes_test/smb.nes")
	//cpu.Load("nes_test/donkeykong.nes")
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
