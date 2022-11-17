asmbuild:
	ca65 nes_test/valid.asm
	ld65 nes_test/valid.o -o nes_test/valid.nes

testbuild:
	ca65 6502_test/first.asm
	ld65 6502_test/first.o -C 6502_test/test.cfg -o 6502_test/first.nes

build:
	go build .

run: build
	./6502