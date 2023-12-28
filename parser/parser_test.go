package parser

import (
	"testing"

	"github.com/simomu-github/sfflt_lang/ast"
	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/token"
)

func TestParsePrimary(t *testing.T) {
	input := "123 'a' true false"
	lexer := lexer.New(input)
	parser := New(lexer)
	exps := parser.ParseProgram()

	intLiteral, ok := exps[0].(ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Not IntegerLiteral")
	}

	if intLiteral.Token.Type != token.INT {
		t.Fatalf("TokenType not INT")
	}

	if intLiteral.Value != 123 {
		t.Fatalf("Value does not match")
	}

	charLiteral, ok := exps[1].(ast.CharLiteral)

	if !ok {
		t.Fatalf("Not CharLiteral")
	}

	if charLiteral.Token.Type != token.CHAR {
		t.Fatalf("TokenType not CHAR")
	}

	if charLiteral.Value != "a" {
		t.Fatalf("Value does not match")
	}

	boolLiteral, ok := exps[2].(ast.BooleanLiteral)

	if !ok {
		t.Fatalf("Not BooleanLiteral")
	}

	if boolLiteral.Token.Type != token.TRUE {
		t.Fatalf("TokenType not TRUE")
	}

	if boolLiteral.Value != true {
		t.Fatalf("Value does not match")
	}

	boolLiteral, ok = exps[3].(ast.BooleanLiteral)

	if !ok {
		t.Fatalf("Not BooleanLiteral")
	}

	if boolLiteral.Token.Type != token.FALSE {
		t.Fatalf("TokenType not FALSE")
	}

	if boolLiteral.Value != false {
		t.Fatalf("Value does not match")
	}
}
