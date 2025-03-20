package parser

import (
	"testing"

	"github.com/simomu-github/sfflt_lang/ast"
	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/token"
)

func TestParsePrimary(t *testing.T) {
	input := "123; 'a'; true; false; a; getc;"
	lexer := lexer.New("script", input)
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

func TestParseArray(t *testing.T) {
	input := "[1, 2, 3];"
	lexer := lexer.New("script", input)
	parser := New(lexer)
	stmts := parser.ParseProgram()

	stmt := stmts[0].(ast.ExpressionStatement)
	arrayLiteral, ok := stmt.Expression.(ast.ArrayLiteral)

	if !ok {
		t.Fatalf("Not ArrayLiteral")
	}

	element1 := arrayLiteral.Elements[0].(ast.IntegerLiteral)
	if element1.Value != 1 {
		t.Fatalf("Array element does not match")
	}

	element2 := arrayLiteral.Elements[1].(ast.IntegerLiteral)
	if element2.Value != 2 {
		t.Fatalf("Array element does not match")
	}

	element3 := arrayLiteral.Elements[2].(ast.IntegerLiteral)
	if element3.Value != 3 {
		t.Fatalf("Array element does not match")
	}
}

func TestParseString(t *testing.T) {
	input := "\"abc\";"
	lexer := lexer.New("script", input)
	parser := New(lexer)
	stmts := parser.ParseProgram()

	stmt := stmts[0].(ast.ExpressionStatement)
	stringLiteral, ok := stmt.Expression.(ast.StringLiteral)

	if !ok {
		t.Fatalf("Not StringLiteral")
	}

	if stringLiteral.Value != "abc" {
		t.Fatalf("Value does not match")
	}
}

func TestParseArgumentVariable(t *testing.T) {
	input := "var a = 0; a = 0; f(a, 1, 2); func f(a, b, c) { 1 + b; }"
	lexer := lexer.New("script", input)
	parser := New(lexer)
	stmt := parser.ParseProgram()

	f, ok := stmt[3].(ast.Function)
	if !ok {
		t.Fatalf("Statement is not function")
	}

	expr, ok := f.Body[0].(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Body is not ExpressionStatement")
	}

	b, ok := expr.Expression.(ast.Binary)
	if !ok {
		t.Fatalf("Expression is not Binary")
	}

	v := b.Right.(ast.Variable)

	if v.Type != ast.ARGUMENT {
		t.Fatalf("Variale is not argument")
	}

	if v.ArgumentIndex != 2 {
		t.Fatalf("Argument index is not match")
	}

	if v.RelativeIndex != 1 {
		t.Fatalf("Relative index is not match")
	}
}

func TestParseLocalVariable(t *testing.T) {
	input := "{ var a = 0; { var a = 1; a; } a; }"
	lexer := lexer.New("script", input)
	parser := New(lexer)
	stmt := parser.ParseProgram()

	block, ok := stmt[0].(ast.Block)
	if !ok {
		t.Fatalf("Statement is not block")
	}

	v, ok := block.Statements[0].(ast.Var)
	if !ok {
		t.Fatalf("Body is not ExpressionStatement")
	}

	if v.Identifier.Literal != "a" {
		t.Fatalf("Declared variable name is not match")
	}

	block2, ok := block.Statements[1].(ast.Block)
	if !ok {
		t.Fatalf("Statement is not block")
	}

	v2, ok := block2.Statements[0].(ast.Var)
	if !ok {
		t.Fatalf("Body is not ExpressionStatement")
	}

	if v2.Identifier.Literal != "a" {
		t.Fatalf("Declared variable name is not match")
	}

	e1 := block2.Statements[1].(ast.ExpressionStatement)
	variable1 := e1.Expression.(ast.Variable)
	if !ok {
		t.Fatalf("Not Variable")
	}

	if variable1.Identifier.Literal != "a" {
		t.Fatalf("Identifier literal is not match")
	}

	if variable1.ScopeDepth != 2 {
		t.Fatalf("Variable scope depth is not match")
	}

	if variable1.LocalIndex != 0 {
		t.Fatalf("Variable local index is not match")
	}

	e := block.Statements[2].(ast.ExpressionStatement)
	variable := e.Expression.(ast.Variable)
	if !ok {
		t.Fatalf("Not Variable")
	}

	if variable.Identifier.Literal != "a" {
		t.Fatalf("Identifier literal is not match")
	}

	if variable.ScopeDepth != 1 {
		t.Fatalf("Variable scope depth is not match")
	}

	if variable.LocalIndex != 0 {
		t.Fatalf("Variable local index is not match")
	}
}

func TestCall(t *testing.T) {
	input := "test(1, 2)"
	lexer := lexer.New("script", input)
	parser := New(lexer)
	exp := parser.parseExpression()

	call, ok := exp.(ast.Call)
	if !ok {
		t.Fatalf("Not Call")
	}

	if call.Callee.Literal != "test" {
		t.Fatalf("Callee is not match")
	}

	if len(call.Arguments) != 2 {
		t.Fatalf("Arguments count is not match")
	}

	first, ok := call.Arguments[0].(ast.IntegerLiteral)
	if !ok {
		t.Fatalf("First arguments is not integer literal")
	}

	if first.Value != 1 {
		t.Fatalf("First arguments is not match")
	}

	second, ok := call.Arguments[1].(ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Second arguments is not integer literal")
	}

	if second.Value != 2 {
		t.Fatalf("Second arguments is not match")
	}
}

func TestParseUnary(t *testing.T) {
	input := "-123"
	lexer := lexer.New("script", input)
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
	lexer := lexer.New("script", input)
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

func TestParseIndex(t *testing.T) {
	input := "a[0]"
	lexer := lexer.New("script", input)
	parser := New(lexer)
	expr := parser.parseExpression()

	index, ok := expr.(ast.Index)

	if !ok {
		t.Fatalf("Not Index")
	}

	receiver, ok := index.Receiver.(ast.Variable)

	if !ok {
		t.Fatalf("Left expression does not Variable")
	}

	if receiver.Identifier.Literal != "a" {
		t.Fatalf("Receiver expression does not match")
	}

	i, ok := index.Index.(ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Index expression does not IntegerLiteral")
	}

	if i.Value != 0 {
		t.Fatalf("Index expression does not match")
	}
}

func TestParseIndexWithArrayLiteral(t *testing.T) {
	input := "[1][(1 + 2) * 3]"
	lexer := lexer.New("script", input)
	parser := New(lexer)
	expr := parser.parseExpression()

	index, ok := expr.(ast.Index)

	if !ok {
		t.Fatalf("Not Index")
	}

	_, ok = index.Receiver.(ast.ArrayLiteral)

	if !ok {
		t.Fatalf("Left expression does not Variable")
	}

	_, ok = index.Index.(ast.Binary)

	if !ok {
		t.Fatalf("Index expression does not Binary")
	}
}

func TestParseIndexWithCall(t *testing.T) {
	input := "call()[call()]"
	lexer := lexer.New("script", input)
	parser := New(lexer)
	expr := parser.parseExpression()

	index, ok := expr.(ast.Index)

	if !ok {
		t.Fatalf("Not Index")
	}

	_, ok = index.Receiver.(ast.Call)

	if !ok {
		t.Fatalf("Left expression does not Call")
	}

	_, ok = index.Index.(ast.Call)

	if !ok {
		t.Fatalf("Index expression does not Call")
	}
}

func TestParseTerm(t *testing.T) {
	input := "(4 - 3) * (2 + 1)"
	lexer := lexer.New("script", input)
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
	lexer := lexer.New("script", input)
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
	lexer := lexer.New("script", input)
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

func TestParseAndOr(t *testing.T) {
	input := "true || false && true"
	lexer := lexer.New("script", input)
	parser := New(lexer)
	expr := parser.parseExpression()

	binary, ok := expr.(ast.Binary)

	if !ok {
		t.Fatalf("Not Binary")
	}

	if binary.Operator.Type != token.OR {
		t.Fatalf("Operator not OR get: %s", binary.Operator.Type)
	}

	left, ok := binary.Left.(ast.BooleanLiteral)

	if !ok {
		t.Fatalf("Left expression does not BooleanLiteral")
	}

	if left.Value != true {
		t.Fatalf("Left expression value not match")
	}

	right, ok := binary.Right.(ast.Binary)

	if !ok {
		t.Fatalf("Right expression does not Binary")
	}

	if right.Operator.Type != token.AND {
		t.Fatalf("Right binary expression oeprator not match")
	}

	rl, ok := right.Left.(ast.BooleanLiteral)

	if !ok {
		t.Fatalf("Right left expression does not BooleanLiteral")
	}

	if rl.Value != false {
		t.Fatalf("Right left expression value not match")
	}

	ll, ok := right.Right.(ast.BooleanLiteral)

	if !ok {
		t.Fatalf("Right right expression does not BooleanLiteral")
	}

	if ll.Value != true {
		t.Fatalf("Right right expression value not match")
	}
}

func TestParsePut(t *testing.T) {
	input := "putn 1;"
	lexer := lexer.New("script", input)
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

func TestParseReturn(t *testing.T) {
	input := "func test() { return 1; }"
	lexer := lexer.New("script", input)
	parser := New(lexer)
	stmt := parser.ParseProgram()

	funcStmt, ok := stmt[0].(ast.Function)
	if !ok {
		t.Fatalf("Statement is not function")
	}

	r, ok := funcStmt.Body[0].(ast.Return)

	if !ok {
		t.Fatalf("Body is not Return")
	}

	intLiteral, ok := r.Value.(ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Expression is not IntegerLiteral")
	}

	if intLiteral.Value != 1 {
		t.Fatalf("IntegerLiteral value is not match")
	}
}

func TestParseBreak(t *testing.T) {
	input := "while(true) break;"
	lexer := lexer.New("script", input)
	parser := New(lexer)
	stmt := parser.ParseProgram()

	whileStmt, ok := stmt[0].(ast.While)
	if !ok {
		t.Fatalf("Statement is not while")
	}

	b, ok := whileStmt.Body.(ast.Break)

	if !ok {
		t.Fatalf("Body is not Break")
	}

	tok := b.Token

	if tok.Type != token.BREAK {
		t.Fatalf("Break token is not match")
	}
}

func TestParseAssign(t *testing.T) {
	input := "a = b = 2"
	lexer := lexer.New("script", input)
	parser := New(lexer)
	expr := parser.parseExpression()

	assign, ok := expr.(ast.Assign)
	if !ok {
		t.Fatalf("Statement is not Assign")
	}

	target, ok := assign.Target.(ast.Variable)
	if !ok {
		t.Fatalf("Target is not Variable")
	}

	if target.Identifier.Literal != "a" {
		t.Fatalf("Target literal is not match")
	}

	right, ok := assign.Expression.(ast.Assign)
	if !ok {
		t.Fatalf("assign expression it not assign")
	}

	target, ok = right.Target.(ast.Variable)
	if !ok {
		t.Fatalf("Target is not Variable")
	}

	if target.Identifier.Literal != "b" {
		t.Fatalf("Target literal is not match")
	}

	rr, ok := right.Expression.(ast.IntegerLiteral)
	if !ok {
		t.Fatalf("assign expression is not IntegerLiteral")
	}
	if rr.Value != 2 {
		t.Fatalf("assign value is not match")
	}
}

func TestParseAssignToIndex(t *testing.T) {
	input := "a[0] = 2"
	lexer := lexer.New("script", input)
	parser := New(lexer)
	expr := parser.parseExpression()

	assign, ok := expr.(ast.Assign)
	if !ok {
		t.Fatalf("Statement is not Assign")
	}

	target, ok := assign.Target.(ast.Index)
	if !ok {
		t.Fatalf("Target is not Index")
	}
	receiver, ok := target.Receiver.(ast.Variable)

	if !ok {
		t.Fatalf("Left expression does not Variable")
	}

	if receiver.Identifier.Literal != "a" {
		t.Fatalf("Receiver expression does not match")
	}

	i, ok := target.Index.(ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Index expression does not IntegerLiteral")
	}

	if i.Value != 0 {
		t.Fatalf("Index expression does not match")
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
	lexer := lexer.New("script", input)
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

func TestParseFunction(t *testing.T) {
	input := "func test(a, b) { putc 'a'; }"
	lexer := lexer.New("script", input)
	parser := New(lexer)
	stmt := parser.ParseProgram()

	funcStmt, ok := stmt[0].(ast.Function)
	if !ok {
		t.Fatalf("Statement is not function")
	}

	if funcStmt.Name.Literal != "test" {
		t.Fatalf("Function name is not match")
	}

	if len(funcStmt.Params) != 2 {
		t.Fatalf("Params count is not match")
	}

	if funcStmt.Params[0].Literal != "a" {
		t.Fatalf("First params literal is not match")
	}

	if funcStmt.Params[1].Literal != "b" {
		t.Fatalf("Second params literal is not match")
	}

	_, ok = funcStmt.Body[0].(ast.PutStatement)

	if !ok {
		t.Fatalf("Body is not PutStatement")
	}
}

func TestParseBlock(t *testing.T) {
	input := "{ a; b; }"
	lexer := lexer.New("script", input)
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
	lexer := lexer.New("script", input)
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
	lexer := lexer.New("script", input)
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

func TestParseFor(t *testing.T) {
	input := `for (;;) true;`
	lexer := lexer.New("script", input)
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

func TestParseForWithCondition(t *testing.T) {
	input := `for (;100;) true;`
	lexer := lexer.New("script", input)
	parser := New(lexer)
	stmt := parser.ParseProgram()

	whileStmt, ok := stmt[0].(ast.While)
	if !ok {
		t.Fatalf("Statement is not while")
	}

	condition, ok := whileStmt.Condition.(ast.IntegerLiteral)
	if !ok {
		t.Fatalf("condition is not IntegerLiteral")
	}

	if condition.Value != 100 {
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

func TestParseForWithInitializer(t *testing.T) {
	input := `for (var i = 0;;) true;`
	lexer := lexer.New("script", input)
	parser := New(lexer)
	stmt := parser.ParseProgram()

	block, ok := stmt[0].(ast.Block)
	if !ok {
		t.Fatalf("Statement is not block")
	}
	if len(block.Statements) != 2 {
		t.Fatalf("Block statements length is not match")
	}

	_, ok = block.Statements[0].(ast.Var)
	if !ok {
		t.Fatalf("Initializer is not Var statement")
	}

	whileStmt, ok := block.Statements[1].(ast.While)
	if !ok {
		t.Fatalf("While statement not found")
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

func TestParseForWithIterator(t *testing.T) {
	input := `for (;;1) true;`
	lexer := lexer.New("script", input)
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

	body, ok := whileStmt.Body.(ast.Block)
	if !ok {
		t.Fatalf("Body statement is not Block")
	}

	if len(body.Statements) != 2 {
		t.Fatalf("Block statements length is not match")
	}

	thenStmt, ok := body.Statements[0].(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("First statement is not ExpressionStatement")
	}

	thenExpr, ok := thenStmt.Expression.(ast.BooleanLiteral)
	if !ok {
		t.Fatalf("Body expression statement is not BooleanLiteral")
	}
	if thenExpr.Value != true {
		t.Fatalf("Body value is not match")
	}

	iteratorStmt, ok := body.Statements[1].(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Second statement is not ExpressionStatement")
	}

	iterator, ok := iteratorStmt.Expression.(ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Iterator expression statement is not IntegerLiteral")
	}
	if iterator.Value != 1 {
		t.Fatalf("Body value is not match")
	}
}

func TestParseInclude(t *testing.T) {
	input := `include "../fixtures/include.sflt";`
	lexer := lexer.New("script", input)
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

	if intLiteral.Value != 1 {
		t.Fatalf("Value does not match")
	}
}
