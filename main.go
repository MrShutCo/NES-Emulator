package main

import (
	"6502/cpu"
	"6502/ppu"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten"
)

type Game struct{}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {
	// Write your game's logical update.
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	//ppu.DrawNameTable0(screen)
	ppu.DrawImage(screen, ppu.PATTERN_TABLE_0, 0)
	ppu.DrawImage(screen, ppu.PATTERN_TABLE_1, 150)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	//NESGame()
	ExecuteNESTest()
}

func NESGame() {
	game := &Game{}
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Your game's title")

	cpu.Reset()
	cpu.LoadMaps()
	//cpu.Load("nes-test-roms/tutor/tutor.nes")
	cpu.Load("nes_test/donkeykong.nes")
	cpu.Start()
	cpu.PC = cpu.GetWordAt(cpu.NMI_VECTOR)
	Image, _ := ebiten.NewImage(512, 512, ebiten.FilterDefault)
	ppu.Image = Image

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func ExecuteNESTest() {
	cpu.Reset()
	cpu.LoadMaps()
	cpu.Load("nes_test/nestest.nes")
	cpu.Start()
	cpu.PC = 0xC000

	os.Truncate("access.log", 0)
	f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for i := 0; i < 5000; i++ {
		output := cpu.Execute()
		//fmt.Fscanln(os.Stdin)
		f.Write([]byte(output + "\n"))
		//fmt.Print(output)
		//fmt.Println(cpu.RAM[0x2000:0x2100])
	}
}
