package compiler

const (
	PUSH    = "FF"
	DUP     = "FTF"
	SWAP    = "FTL"
	DISCARD = "FTT"
	COPY    = "FLF"
	SLIDE   = "FLT"

	ADD = "LFFF"
	SUB = "LFFL"
	MUL = "LFFT"
	DIV = "LFLF"
	MOD = "LFLL"

	STORE    = "LLF"
	RETRIEVE = "LLL"

	GETC = "LTLF"
	GETN = "LTLL"
	PUTC = "LTFF"
	PUTN = "LTFL"

	LABEL          = "TFF"
	JUMP           = "TFT"
	JUMP_WHEN_ZERO = "TLF"
	JUMP_WHEN_NEGA = "TLL"

	CALLSUB = "TFL"
	ENDSUB  = "TLT"

	END = "TTT"
)

const (
	POSI = "F"
	NEGA = "L"

	ZERO      = "FF"
	ONE       = "FL"
	MINUS_ONE = "LL"
)

type InstructionType string
