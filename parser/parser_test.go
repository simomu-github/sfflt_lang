package parser

import (
	"testing"

	"github.com/simomu-github/sfflt_lang/ast"
	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/token"
)

func TestParsePrimary(t *testing.T) {
	input := "123; 'a'; true; false; a; getc;"
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

	stmt = stmts[4].(ast.ExpressionStatement)
	variable, ok := stmt.Expression.(ast.Variable)

	if !ok {
		t.Fatalf("Not Variable")
	}

	if variable.Identifier.Literal != "a" {
		t.Fatalf("variable literal is not match")
	}

	stmt = stmts[5].(ast.ExpressionStatement)
	get, ok := stmt.Expression.(ast.Get)

	if !ok {
		t.Fatalf("Not get")
	}

	if get.Token.Type != token.GETC {
		t.Fatalf("token type is not match")
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

func TestParseComparison(t *testing.T) {
	input := "(a + b) < c"
	lexer := lexer.New(input)
	parser := New(lexer)
	expr := parser.parseExpression()

	binary, ok := expr.(ast.Binary)

	if !ok {
		t.Fatalf("Not Binary")
	}

	if binary.Operator.Type != token.LT {
		t.Fatalf("Operator not LT get: %s", binary.Operator.Type)
	}

	left, ok := binary.Left.(ast.Binary)

	if !ok {
		t.Fatalf("Left expression does not Binary")
	}

	if left.Operator.Type != token.PLUS {
		t.Fatalf("Left expression operator type does not match")
	}

	ll, ok := left.Left.(ast.Variable)
	if !ok {
		t.Fatalf("Left left expression does not Variable")
	}

	if ll.Identifier.Literal != "a" {
		t.Fatalf("Left left variable does not match")
	}

	lr, ok := left.Right.(ast.Variable)
	if !ok {
		t.Fatalf("Left right expression does not Variable")
	}

	if lr.Identifier.Literal != "b" {
		t.Fatalf("Left right variable identifier does not match")
	}

	right, ok := binary.Right.(ast.Variable)

	if !ok {
		t.Fatalf("Right expression does not Variable")
	}

	if right.Identifier.Literal != "c" {
		t.Fatalf("Right variable identifier does not match")
	}
}

func TestParseEquality(t *testing.T) {
	input := "true != false"
	lexer := lexer.New(input)
	parser := New(lexer)
	expr := parser.parseExpression()

	binary, ok := expr.(ast.Binary)

	if !ok {
		t.Fatalf("Not Binary")
	}

	if binary.Operator.Type != token.NOT_EQ {
		t.Fatalf("Operator not NOT_EQ get: %s", binary.Operator.Type)
	}

	left, ok := binary.Left.(ast.BooleanLiteral)

	if !ok {
		t.Fatalf("Left expression does not BooleanLiteral")
	}

	if left.Value != true {
		t.Fatalf("Left expression value not match")
	}

	right, ok := binary.Right.(ast.BooleanLiteral)

	if !ok {
		t.Fatalf("Right expression does not BooleanLiteral")
	}

	if right.Value != false {
		t.Fatalf("Right expression value not match")
	}
}

func TestParsePut(t *testing.T) {
	input := "putn 1;"
	lexer := lexer.New(input)
	parser := New(lexer)
	stmt := parser.ParseProgram()

	put, ok := stmt[0].(ast.PutStatement)
	if !ok {
		t.Fatalf("Statement is not put")
	}

	if put.Token.Type != token.PUTN {
		t.Fatalf("Token is not PUTN")
	}

	intLiteral, ok := put.Expression.(ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Expression is not IntegerLiteral")
	}

	if intLiteral.Value != 1 {
		t.Fatalf("IntegerLiteral value is not match")
	}
}

func TestParseAssign(t *testing.T) {
	input := "a = 2"
	lexer := lexer.New(input)
	parser := New(lexer)
	expr := parser.parseExpression()

	assign, ok := expr.(ast.Assign)
	if !ok {
		t.Fatalf("Statement is not Assign")
	}

	target := assign.Target
	if target.Literal != "a" {
		t.Fatalf("Target literal is not match")
	}

	right, ok := assign.Expression.(ast.IntegerLiteral)
	if !ok {
		t.Fatalf("assign expression is not IntegerLiteral")
	}
	if right.Value != 2 {
		t.Fatalf("assign value is not match")
	}
}

func TestParseVar(t *testing.T) {
	input := "var a = 1;"
	lexer := lexer.New(input)
	parser := New(lexer)
	stmt := parser.ParseProgram()

	varStmt, ok := stmt[0].(ast.Var)
	if !ok {
		t.Fatalf("Statement is not var")
	}

	if varStmt.Identifier.Type != token.IDENT {
		t.Fatalf("Var statement does not have Identifier")
	}

	intLiteral, ok := varStmt.Expression.(ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Expression is not IntegerLiteral")
	}

	if intLiteral.Value != 1 {
		t.Fatalf("IntegerLiteral value is not match")
	}
}

func TestParseBlock(t *testing.T) {
	input := "{ a; b; }"
	lexer := lexer.New(input)
	parser := New(lexer)
	stmts := parser.ParseProgram()

	group := stmts[0].(ast.Block)

	stmt := group.Statements[0].(ast.ExpressionStatement)
	variable, ok := stmt.Expression.(ast.Variable)

	if !ok {
		t.Fatalf("Not Variable")
	}

	if variable.Identifier.Literal != "a" {
		t.Fatalf("Identifier literal is not match")
	}

	stmt = group.Statements[1].(ast.ExpressionStatement)
	variable, ok = stmt.Expression.(ast.Variable)

	if !ok {
		t.Fatalf("Not Variable")
	}

	if variable.Identifier.Literal != "b" {
		t.Fatalf("Identifier literal is not match")
	}
}

func TestParseIf(t *testing.T) {
	input := `
if (true)
    true;
else
    false;
`
	lexer := lexer.New(input)
	parser := New(lexer)
	stmt := parser.ParseProgram()

	ifStmt, ok := stmt[0].(ast.If)
	if !ok {
		t.Fatalf("Statement is not if")
	}

	condition, ok := ifStmt.Condition.(ast.BooleanLiteral)
	if !ok {
		t.Fatalf("condition is not BooleanLiteral")
	}

	if condition.Value != true {
		t.Fatalf("condition value is not match")
	}

	thenStmt, ok := ifStmt.Then.(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Then statement is not ExpressionStatement")
	}
	thenExpr, ok := thenStmt.Expression.(ast.BooleanLiteral)
	if !ok {
		t.Fatalf("Then expression statement is not BooleanLiteral")
	}
	if thenExpr.Value != true {
		t.Fatalf("Then value is not match")
	}

	elseStmt, ok := ifStmt.Else.(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Else statement is not ExpressionStatement")
	}
	elseExpr, ok := elseStmt.Expression.(ast.BooleanLiteral)
	if !ok {
		t.Fatalf("Else expression statement is not BooleanLiteral")
	}
	if elseExpr.Value != false {
		t.Fatalf("Else value is not match")
	}
}

func TestParseWhile(t *testing.T) {
	input := `while (true) true;`
	lexer := lexer.New(input)
	parser := New(lexer)
	stmt := parser.ParseProgram()

	whileStmt, ok := stmt[0].(ast.While)
	if !ok {
		t.Fatalf("Statement is not while")
	}

	condition, ok := whileStmt.Condition.(ast.BooleanLiteral)
	if !ok {
		t.Fatalf("condition is not BooleanLiteral")
	}

	if condition.Value != true {
		t.Fatalf("condition value is not match")
	}

	thenStmt, ok := whileStmt.Body.(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Body statement is not ExpressionStatement")
	}
	thenExpr, ok := thenStmt.Expression.(ast.BooleanLiteral)
	if !ok {
		t.Fatalf("Body expression statement is not BooleanLiteral")
	}
	if thenExpr.Value != true {
		t.Fatalf("Body value is not match")
	}
}
