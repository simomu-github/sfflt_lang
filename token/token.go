package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"
	CHAR  = "CHAR"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MOD      = "%"

	LT   = "<"
	LTEQ = "<="
	GT   = ">"
	GTEQ = ">="

	EQ     = "=="
	NOT_EQ = "!="

	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	VAR    = "VAR"
	FUNC   = "FUNC"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	ELSE   = "ELSE"
	WHILE  = "WHILE"
	RETURN = "RETURN"

	PUTN = "PUTN"
	PUTC = "PUTC"

	GETN = "GETN"
	GETC = "GETC"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

var keywords = map[string]TokenType{
	"var":    VAR,
	"func":   FUNC,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"return": RETURN,

	"putn": PUTN,
	"putc": PUTC,
	"getn": GETN,
	"getc": GETC,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
