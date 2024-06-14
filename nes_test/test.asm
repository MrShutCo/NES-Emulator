.segment "HEADER"
.byte $4e, $45, $53, $1a, $01, $01, $00, $00

.segment "CODE"
.proc irq_handler
  RTI
.endproc

.proc nmi_handler
  LDA #$45
  STA $2007,X
  INX
  RTI
.endproc

.proc reset_handler
  SEI         ; 0x78
  CLD         ; 0x D8
  LDX #$00
  STX $2000
  STX $2001
vblankwait:
  BIT $2002
  BPL vblankwait
  JMP main
.endproc

.proc main
  LDX $2002   ; 0xAE 0x02 0x20
  LDX #$3f
  STX $2006
  LDX #$00
  STX $2006
  LDA #$29
  STA $2007
  LDA #%00011110
  STA $2001
forever:
  JMP forever
.endproc

.segment "VECTORS"
.addr nmi_handler, reset_handler, irq_handler

.segment "CHARS"
.res 8192
.segment "STARTUP"