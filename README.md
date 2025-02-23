# NES Emulator
This is a side project to implement the Nintendo Entertainment System from scratch using Golang and the ebiten game engine. 

![image](https://github.com/user-attachments/assets/f9927846-d23b-4b25-9cd6-34ef220bff8c)

Donkey Kong game is fully playable at 60FPS. The biggest overhead is issues with how NES graphics works vs ebiten. Writing individual pixels to a texture is slow, so caching was required to make
this run smoothly. If a tiles palette and index hasnt changed then we dont need to waste time redrawing the frame

# What's done so far
* 6502 Processor official opcodes are implemented using https://www.masswerk.at/6502/6502_instruction_set.html as reference, logs verified by https://github.com/christopherpow/nes-test-roms/blob/master/other/nestest.log 
* PPU can render single screen games such Donkey Kong, displaying sprites in full colour ([reference here](https://www.nesdev.org/wiki/PPU_OAM)) and background [reference here](https://www.nesdev.org/wiki/PPU_pattern_tables). G



# What's left to do
* More automation testing
* Various other flags to implement other features
* Sprite 0 hits
* Screen scrolling
* Audio
* Line-by-line rendering to mimic the real PPU instead of once per frame
