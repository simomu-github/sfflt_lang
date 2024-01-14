package lexer

import (
	"testing"

	"github.com/simomu-github/sfflt_lang/token"
)

func TestScanToken(t *testing.T) {
	input := `(){};
+-*/%=!
==
!=
<><=>=
'a'123
var func if else while true false putn putc getn getc hoge_fuga0
`

	expects := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LPAREN, "(", 1, 1},
		{token.RPAREN, ")", 1, 2},
		{token.LBRACE, "{", 1, 3},
		{token.RBRACE, "}", 1, 4},
		{token.SEMICOLON, ";", 1, 5},

		{token.PLUS, "+", 2, 1},
		{token.MINUS, "-", 2, 2},
		{token.ASTERISK, "*", 2, 3},
		{token.SLASH, "/", 2, 4},
		{token.MOD, "%", 2, 5},
		{token.ASSIGN, "=", 2, 6},
		{token.BANG, "!", 2, 7},

		{token.EQ, "==", 3, 2},
		{token.NOT_EQ, "!=", 4, 2},

		{token.LT, "<", 5, 1},
		{token.GT, ">", 5, 2},
		{token.LTEQ, "<=", 5, 4},
		{token.GTEQ, ">=", 5, 6},

		{token.CHAR, "a", 6, 3},
		{token.INT, "123", 6, 6},

		{token.VAR, "var", 7, 3},
		{token.FUNC, "func", 7, 8},
		{token.IF, "if", 7, 11},
		{token.ELSE, "else", 7, 16},
		{token.WHILE, "while", 7, 22},
		{token.TRUE, "true", 7, 27},
		{token.FALSE, "false", 7, 33},
		{token.PUTN, "putn", 7, 38},
		{token.PUTC, "putc", 7, 43},
		{token.GETN, "getn", 7, 48},
		{token.GETC, "getc", 7, 53},
		{token.IDENT, "hoge_fuga0", 7, 64},

		{token.EOF, string(byte(0)), 8, 0},
	}

	lexer := New("script", input)

	for i, expect := range expects {
		token := lexer.ScanToken()
		if token.Type != expect.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, expect.expectedType, token.Type)
		}

		if token.Literal != expect.expectedLiteral {
			t.Fatalf("tests[%d] - ligteral wrong. expected=%q, got=%q", i, expect.expectedLiteral, token.Literal)
		}

		if token.Line != expect.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d", i, expect.expectedLine, token.Line)
		}

		if token.Column != expect.expectedColumn {
			t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d", i, expect.expectedColumn, token.Column)
		}

	}
}
