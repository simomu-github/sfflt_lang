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

func TestParseUnary(t *testing.T) {
	input := "-123"
	lexer := lexer.New(input)
	parser := New(lexer)
	exps := parser.ParseProgram()

	unary, ok := exps[0].(ast.Unary)

	if !ok {
		t.Fatalf("Not Unary")
	}

	if unary.Operator.Type != token.MINUS {
		t.Fatalf("TokenType not INT")
	}

	intLiteral, ok := unary.Right.(ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Right expression does not IntegerLiteral")
	}

	if intLiteral.Value != 123 {
		t.Fatalf("Right expression does not match")
	}
}

func TestParseFactor(t *testing.T) {
	input := "2 * -3"
	lexer := lexer.New(input)
	parser := New(lexer)
	exps := parser.ParseProgram()

	binary, ok := exps[0].(ast.Binary)

	if !ok {
		t.Fatalf("Not Binary")
	}

	if binary.Operator.Type != token.ASTERISK {
		t.Fatalf("Operator not ASTERISK")
	}

	left, ok := binary.Left.(ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Left expression does not IntegerLiteral")
	}

	if left.Value != 2 {
		t.Fatalf("Left expression does not match")
	}

	right, ok := binary.Right.(ast.Unary)

	if !ok {
		t.Fatalf("Right expression does not Unary")
	}

	if right.Operator.Type != token.MINUS {
		t.Fatalf("TokenType not MINUS")
	}

	intLiteral, ok := right.Right.(ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Right expression does not IntegerLiteral")
	}

	if intLiteral.Value != 3 {
		t.Fatalf("Right expression does not match")
	}
}

func TestParseTerm(t *testing.T) {
	input := "2 + -3"
	lexer := lexer.New(input)
	parser := New(lexer)
	exps := parser.ParseProgram()

	binary, ok := exps[0].(ast.Binary)

	if !ok {
		t.Fatalf("Not Binary")
	}

	if binary.Operator.Type != token.PLUS {
		t.Fatalf("Operator not PLUS")
	}

	left, ok := binary.Left.(ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Left expression does not IntegerLiteral")
	}

	if left.Value != 2 {
		t.Fatalf("Left expression does not match")
	}

	right, ok := binary.Right.(ast.Unary)

	if !ok {
		t.Fatalf("Right expression does not Unary %q", binary.Right)
	}

	if right.Operator.Type != token.MINUS {
		t.Fatalf("TokenType not MINUS")
	}

	intLiteral, ok := right.Right.(ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Right expression does not IntegerLiteral")
	}

	if intLiteral.Value != 3 {
		t.Fatalf("Right expression does not match")
	}
}
