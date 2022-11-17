.segment "HEADER"
  LDX #$FF
  STX $2001
  LDY #$45
  STY $2002
forever:
  JMP forever