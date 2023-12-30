package parser

import (
	"testing"

	"github.com/simomu-github/sfflt_lang/ast"
	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/token"
)

func TestParsePrimary(t *testing.T) {
	input := "123; 'a'; true; false;"
	lexer := lexer.New(input)
	parser := New(lexer)
	stmts := parser.ParseProgram()

	stmt := stmts[0].(ast.ExpressionStatement)
	intLiteral, ok := stmt.Expression.(ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Not IntegerLiteral")
	}

	if intLiteral.Token.Type != token.INT {
		t.Fatalf("TokenType not INT")
	}

	if intLiteral.Value != 123 {
		t.Fatalf("Value does not match")
	}

	stmt = stmts[1].(ast.ExpressionStatement)
	charLiteral, ok := stmt.Expression.(ast.CharLiteral)

	if !ok {
		t.Fatalf("Not CharLiteral")
	}

	if charLiteral.Token.Type != token.CHAR {
		t.Fatalf("TokenType not CHAR")
	}

	if charLiteral.Value != "a" {
		t.Fatalf("Value does not match")
	}

	stmt = stmts[2].(ast.ExpressionStatement)
	boolLiteral, ok := stmt.Expression.(ast.BooleanLiteral)

	if !ok {
		t.Fatalf("Not BooleanLiteral")
	}

	if boolLiteral.Token.Type != token.TRUE {
		t.Fatalf("TokenType not TRUE")
	}

	if boolLiteral.Value != true {
		t.Fatalf("Value does not match")
	}

	stmt = stmts[3].(ast.ExpressionStatement)
	boolLiteral, ok = stmt.Expression.(ast.BooleanLiteral)

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
	exp := parser.parseExpression()

	unary, ok := exp.(ast.Unary)

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
	expr := parser.parseExpression()

	binary, ok := expr.(ast.Binary)

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
	input := "(4 - 3) * (2 + 1)"
	lexer := lexer.New(input)
	parser := New(lexer)
	expr := parser.parseExpression()

	binary, ok := expr.(ast.Binary)

	if !ok {
		t.Fatalf("Not Binary")
	}

	if binary.Operator.Type != token.ASTERISK {
		t.Fatalf("Operator not ASTERISK get: %s", binary.Operator.Type)
	}

	left, ok := binary.Left.(ast.Binary)

	if !ok {
		t.Fatalf("Left expression does not Binary")
	}

	if left.Operator.Type != token.MINUS {
		t.Fatalf("Left expression operator type does not match")
	}

	ll, ok := left.Left.(ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Left left expression does not IntegerLiteral")
	}

	if ll.Value != 4 {
		t.Fatalf("Left left expression value does not match")
	}

	lr, ok := left.Right.(ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Left left expression does not IntegerLiteral")
	}

	if lr.Value != 3 {
		t.Fatalf("Left left expression value does not match")
	}

	right, ok := binary.Right.(ast.Binary)

	if !ok {
		t.Fatalf("Right expression does not Binary")
	}

	if right.Operator.Type != token.PLUS {
		t.Fatalf("Right expression operator type does not match")
	}

	rl, ok := right.Left.(ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Right left expression does not IntegerLiteral")
	}

	if rl.Value != 2 {
		t.Fatalf("Right left expression value does not match")
	}

	rr, ok := right.Right.(ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Right right expression does not IntegerLiteral")
	}

	if rr.Value != 1 {
		t.Fatalf("Right right expression value does not match")
	}
}

func TestParsePutn(t *testing.T) {
	input := "putn 1;"
	lexer := lexer.New(input)
	parser := New(lexer)
	stmt := parser.ParseProgram()

	putn, ok := stmt[0].(ast.PutnStatement)
	if !ok {
		t.Fatalf("Statement is not Putn")
	}

	if putn.Token.Type != token.PUTN {
		t.Fatalf("Token is not PUTN")
	}

	intLiteral, ok := putn.Expression.(ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Expression is not IntegerLiteral")
	}

	if intLiteral.Value != 1 {
		t.Fatalf("IntegerLiteral value is not match")
	}
}
