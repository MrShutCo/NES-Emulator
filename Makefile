asmbuild:
	ca65 nes_test/test.asm
	ld65 nes_test/test.o -t nes -o nes_test/test.nes

testbuild:
	ca65 6502_test/first.asm
	ld65 6502_test/first.o -C 6502_test/test.cfg -o 6502_test/first.nes

build:
	go build .

run: build
	./6502