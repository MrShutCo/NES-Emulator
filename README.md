# NES Emulator
This is a side project to implement the Nintendo Entertainment System from scratch using Golang and the ebiten game engine. 

# What's done so far
* 6502 Processor official opcodes are implemented using https://www.masswerk.at/6502/6502_instruction_set.html as reference, logs verified by https://github.com/christopherpow/nes-test-roms/blob/master/other/nestest.log 
* PPU can render single screen games such Donkey Kong, displaying sprites in full colour ([reference here](https://www.nesdev.org/wiki/PPU_OAM)) and background [reference here](https://www.nesdev.org/wiki/PPU_pattern_tables)

# What's left to do
* More automation testing
* Visual bugs
* Various other flags to implement other features
* Sprite 0 hits
* Screen scrolling
* Line-by-line rendering to mimic the real PPU instead of once per frame
* Audio

