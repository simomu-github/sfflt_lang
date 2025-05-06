package parser

import (
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/simomu-github/sfflt_lang/ast"
	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/lib"
	"github.com/simomu-github/sfflt_lang/token"
)

type Parser struct {
	lexer           *lexer.Lexer
	currentToken    token.Token
	peekToken       token.Token
	isFunction      bool
	nestedLoopCount int
	stackTop        int
	scopes          []map[string]*declaredVariable
	Errors          []string
	VisitedFiles    []string
}

type declaredVariable struct {
	initialized   bool
	typ           variableType
	scopeDepth    int
	argumentIndex int
	localIndex    int
}

const (
	LOCAL    = "LOCAL"
	ARGUMENT = "ARGUMENT"
)

type variableType string

func New(lexer *lexer.Lexer) *Parser {
	return newParser(lexer, []string{})
}

func newParser(lexer *lexer.Lexer, visitedFiles []string) *Parser {
	p := &Parser{
		lexer:        lexer,
		isFunction:   false,
		Errors:       []string{},
		scopes:       []map[string]*declaredVariable{},
		VisitedFiles: append(visitedFiles, lexer.Filename),
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram() []ast.Statement {
	statements := []ast.Statement{}
	for p.currentToken.Type != token.EOF {
		if p.matchToken(token.INCLUDE) {
			stmts := p.parseInclude()
			if stmts != nil {
				statements = append(statements, stmts...)
			}
		}
		stmt := p.parseDeclaration()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		p.nextToken()
	}

	return statements
}

func (p *Parser) HadErrors() bool {
	return len(p.Errors) != 0
}

func (p *Parser) parseDeclaration() ast.Statement {
	defer func() {
		if p.HadErrors() {
			p.skipStatement()
		}
	}()

	if p.matchToken(token.VAR) {
		return p.parseVarDeclaration()
	}

	if p.matchToken(token.FUNC) {
		return p.parseFunctionDeclaration()
	}

	return p.parseStatement()
}

func (p *Parser) parseVarDeclaration() ast.Statement {
	if p.currentToken.Type != token.IDENT {
		p.parseError(p.currentToken, "Expect identifier.")
		return nil
	}
	identifier := p.currentToken
	local := p.declareLocalVariable(identifier)
	p.nextToken()

	if !p.matchToken(token.ASSIGN) {
		p.parseError(p.currentToken, "Expect '=' after identifier.")
		return nil
	}

	p.pushStack()

	expr := p.parseExpression()
	p.markInitializedVariable(identifier)

	if p.currentToken.Type != token.SEMICOLON {
		p.parseError(p.currentToken, "Expect ';' after statement.")
		return nil
	}

	p.popStack()

	isLocal := false
	depth := 0
	index := 0
	if local != nil {
		isLocal = true
		depth = local.scopeDepth
		index = local.localIndex
	}

	return ast.Var{
		Identifier: identifier,
		Expression: expr,
		IsLocal:    isLocal,
		ScopeDepth: depth,
		LocalIndex: index,
	}
}

func (p *Parser) parseFunctionDeclaration() ast.Statement {
	if p.isFunction {
		p.parseError(p.currentToken, "Can not declare function inner function.")
		return nil
	}

	p.beginScope()
	p.isFunction = true

	if p.currentToken.Type != token.IDENT {
		p.parseError(p.currentToken, "Expect function name.")
		return nil
	}
	name := p.currentToken
	p.nextToken()

	if !p.matchToken(token.LPAREN) {
		p.parseError(p.currentToken, "Expect '(' after function name.")
		return nil
	}

	params := []token.Token{}
	if p.currentToken.Type != token.RPAREN {
		for {
			if p.currentToken.Type != token.IDENT {
				p.parseError(p.currentToken, "Expect argument name.")
				return nil
			}

			params = append(params, p.currentToken)
			p.declareArgumentVariable(p.currentToken, len(params))

			p.nextToken()
			if !p.matchToken(token.COMMA) {
				break
			}
		}
	}

	if !p.matchToken(token.RPAREN) {
		p.parseError(p.currentToken, "Expect ')' after parameters.")
		return nil
	}

	if !p.matchToken(token.LBRACE) {
		p.parseError(p.currentToken, "Expect '{' before function body.")
		return nil
	}

	body := p.parseBlock().(ast.Block)

	p.isFunction = false
	p.endScope()

	return ast.Function{Name: name, Params: params, Body: body.Statements}
}

func (p *Parser) parseInclude() []ast.Statement {
	if p.currentToken.Type != token.STRING {
		p.parseError(p.currentToken, "Expect include name.")
		return nil
	}

	if slices.Contains(p.VisitedFiles, p.currentToken.Literal) {
		return nil
	}

	var code string
	if lib, ok := lib.LookupBuilinLibrary(p.currentToken.Literal); ok {
		code = lib
	} else {
		bytes, err := os.ReadFile(p.currentToken.Literal)
		if err != nil {
			p.parseError(p.currentToken, fmt.Sprintf("Including file can not read. (%s)", p.currentToken.Literal))
			return nil
		}
		code = string(bytes)
	}

	lexer := lexer.New(p.currentToken.Literal, code)
	parser := newParser(lexer, p.VisitedFiles)
	statements := parser.ParseProgram()
	if parser.HadErrors() {
		p.Errors = append(p.Errors, parser.Errors...)
	}
	p.VisitedFiles = parser.VisitedFiles

	return statements
}

func (p *Parser) parseStatement() ast.Statement {
	if p.currentToken.Type == token.PUTN ||
		p.currentToken.Type == token.PUTC {
		return p.parsePutStatement()
	}

	if p.matchToken(token.IF) {
		return p.parseIf()
	}

	if p.matchToken(token.WHILE) {
		return p.parseWhile()
	}

	if p.matchToken(token.FOR) {
		return p.parseFor()
	}

	if p.matchToken(token.LBRACE) {
		return p.parseBlock()
	}

	if p.currentToken.Type == token.RETURN {
		return p.parseReturn()
	}

	if p.currentToken.Type == token.BREAK {
		return p.parseBreak()
	}

	expr := p.parseExpression()
	if expr == nil {
		return nil
	}

	if p.currentToken.Type != token.SEMICOLON {
		p.parseError(p.currentToken, "Expect ';' after statement.")
		return nil
	}

	return ast.ExpressionStatement{Expression: expr}
}

func (p *Parser) parsePutStatement() ast.Statement {
	tok := p.currentToken
	p.nextToken()
	expr := p.parseExpression()
	if p.currentToken.Type != token.SEMICOLON {
		p.parseError(p.currentToken, "Expect ';' after statement.")
		return nil
	}

	return ast.PutStatement{Token: tok, Expression: expr}
}

func (p *Parser) parseReturn() ast.Statement {
	if !p.isFunction {
		p.parseError(p.currentToken, "Can not return top-level code.")
		return nil
	}

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
		return ast.Return{Value: nil}
	}

	p.nextToken()
	expr := p.parseExpression()
	if p.currentToken.Type != token.SEMICOLON {
		p.parseError(p.currentToken, "Expect ';' after statement.")
		return nil
	}

	return ast.Return{Value: expr}
}

func (p *Parser) parseBreak() ast.Statement {
	if !p.isInLoop() {
		p.parseError(p.currentToken, "Can not use 'break' out of loop.")
		return nil
	}

	tok := p.currentToken
	if p.peekToken.Type != token.SEMICOLON {
		p.parseError(p.currentToken, "Expect ';' after statement.")
		return nil
	}
	p.nextToken()

	return ast.Break{Token: tok}
}

func (p *Parser) parseIf() ast.Statement {
	if p.currentToken.Type != token.LPAREN {
		p.parseError(p.currentToken, "Expect '(' after if.")
		return nil
	}
	p.nextToken()

	condition := p.parseExpression()

	if p.currentToken.Type != token.RPAREN {
		p.parseError(p.currentToken, "Expect ')' after if condition.")
		return nil
	}
	p.nextToken()

	thenStmt := p.parseDeclaration()
	var elseStmt ast.Statement
	if p.peekToken.Type == token.ELSE {
		p.nextToken()
		p.nextToken()
		elseStmt = p.parseDeclaration()
	}

	return ast.If{Condition: condition, Then: thenStmt, Else: elseStmt}
}

func (p *Parser) parseWhile() ast.Statement {
	p.beginLoop()

	if p.currentToken.Type != token.LPAREN {
		p.parseError(p.currentToken, "Expect '(' after while.")
		return nil
	}
	p.nextToken()

	condition := p.parseExpression()

	if p.currentToken.Type != token.RPAREN {
		p.parseError(p.currentToken, "Expect ')' after while condition.")
		return nil
	}
	p.nextToken()

	body := p.parseDeclaration()

	p.endLoop()
	return ast.While{Condition: condition, Body: body}
}

func (p *Parser) parseFor() ast.Statement {
	p.beginScope()
	p.beginLoop()

	if p.currentToken.Type != token.LPAREN {
		p.parseError(p.currentToken, "Expect '(' after while.")
		return nil
	}
	p.nextToken()

	var initializer ast.Statement
	if p.currentToken.Type == token.SEMICOLON {
		initializer = nil
	} else if p.matchToken(token.VAR) {
		initializer = p.parseVarDeclaration()
	} else {
		expr := p.parseExpression()
		if p.currentToken.Type != token.SEMICOLON {
			p.parseError(p.currentToken, "Expect ';' after for initialzier.")
			return nil
		}
		initializer = ast.ExpressionStatement{Expression: expr}
	}
	p.nextToken()

	var condition ast.Expression
	if p.currentToken.Type != token.SEMICOLON {
		condition = p.parseExpression()
		if p.currentToken.Type != token.SEMICOLON {
			p.parseError(p.currentToken, "Expect ';' after loop condition.")
			return nil
		}
	}
	p.nextToken()

	var iter ast.Expression
	if p.currentToken.Type != token.RPAREN {
		iter = p.parseExpression()
	}

	if p.currentToken.Type != token.RPAREN {
		p.parseError(p.currentToken, "Expect ')' after for clauses.")
		return nil
	}
	p.nextToken()

	body := p.parseDeclaration()

	p.endLoop()

	if iter != nil {
		body = ast.Block{
			Statements: []ast.Statement{
				body,
				ast.ExpressionStatement{Expression: iter},
			},
		}
	}
	if condition == nil {
		condition = ast.BooleanLiteral{Value: true}
	}

	body = ast.While{Condition: condition, Body: body}

	if initializer != nil {
		body = ast.Block{
			Statements: []ast.Statement{
				initializer,
				body,
			},
		}
	}

	p.endLoop()
	p.endScope()

	return body
}

func (p *Parser) parseBlock() ast.Statement {
	p.beginScope()
	stmts := []ast.Statement{}
	for p.currentToken.Type != token.RBRACE {
		if p.currentToken.Type == token.EOF {
			p.parseError(p.currentToken, "Expect '}' after statements.")
			return nil
		}
		stmts = append(stmts, p.parseDeclaration())
		p.nextToken()
	}

	p.endScope()
	return ast.Block{Statements: stmts}
}

func (p *Parser) parseExpression() ast.Expression {
	expr := p.parseAssign()
	p.popStack()
	p.nextToken()
	return expr
}

func (p *Parser) parseAssign() ast.Expression {
	expr := p.parseOr()

	switch p.peekToken.Type {
	case token.ASSIGN:
		p.nextToken()
		p.nextToken()

		p.pushStack()
		right := p.parseAssign()
		target, ok := expr.(ast.Assignable)
		if !ok {
			p.parseError(p.currentToken, "Invalid assignment target.")
			return nil
		}
		if !target.CanAssign() {
			p.parseError(p.currentToken, "Invalid assignment target.")
			return nil
		}
		p.popStack()

		p.popStack()
		return ast.Assign{Target: target, Expression: right}
	}

	return expr
}

func (p *Parser) parseOr() ast.Expression {
	expr := p.parseAnd()
	for p.matchPeekToken(token.OR) {
		p.nextToken()
		operator := p.currentToken
		p.nextToken()
		right := p.parseAnd()
		expr = ast.Binary{Left: expr, Operator: operator, Right: right}
		p.popStack()
	}

	return expr
}

func (p *Parser) parseAnd() ast.Expression {
	expr := p.parseEquality()
	for p.matchPeekToken(token.AND) {
		p.nextToken()
		operator := p.currentToken
		p.nextToken()
		right := p.parseEquality()
		expr = ast.Binary{Left: expr, Operator: operator, Right: right}
		p.popStack()
	}

	return expr
}

func (p *Parser) parseEquality() ast.Expression {
	expr := p.parseComparison()
	switch p.peekToken.Type {
	case token.EQ, token.NOT_EQ:
		p.nextToken()
		operator := p.currentToken
		p.nextToken()
		right := p.parseComparison()
		e := ast.Binary{Left: expr, Operator: operator, Right: right}
		p.popStack()
		return e
	}

	return expr
}

func (p *Parser) parseComparison() ast.Expression {
	expr := p.parseTerm()
	switch p.peekToken.Type {
	case token.LT, token.LTEQ, token.GT, token.GTEQ:
		p.nextToken()
		operator := p.currentToken
		p.nextToken()
		right := p.parseTerm()
		e := ast.Binary{Left: expr, Operator: operator, Right: right}
		p.popStack()
		return e
	}

	return expr
}

func (p *Parser) parseTerm() ast.Expression {
	expr := p.parseFactor()
	for p.matchPeekToken(token.PLUS, token.MINUS) {
		p.nextToken()
		operator := p.currentToken
		p.nextToken()
		right := p.parseFactor()
		expr = ast.Binary{Left: expr, Operator: operator, Right: right}
		p.popStack()
	}

	return expr
}

func (p *Parser) parseFactor() ast.Expression {
	expr := p.parseUnary()
	for p.matchPeekToken(token.ASTERISK, token.SLASH, token.MOD) {
		p.nextToken()
		operator := p.currentToken
		p.nextToken()
		right := p.parseUnary()
		expr = ast.Binary{Left: expr, Operator: operator, Right: right}
		p.popStack()
	}

	return expr
}

func (p *Parser) parseUnary() ast.Expression {
	switch p.currentToken.Type {
	case token.MINUS, token.BANG:
		operator := p.currentToken
		p.nextToken()
		return ast.Unary{Operator: operator, Right: p.parseCall()}
	}

	return p.parseIndex()
}

func (p *Parser) parseIndex() ast.Expression {
	expr := p.parseCall()
	if p.matchPeekToken(token.LBRACKET) {
		p.nextToken()
		p.nextToken()

		index := p.parseExpression()

		if p.currentToken.Type != token.RBRACKET {
			p.parseError(p.currentToken, "Expect ']' after index.")
			return nil
		}

		expr = ast.Index{Receiver: expr, Index: index}
	}

	return expr
}

func (p *Parser) parseCall() ast.Expression {
	if p.currentToken.Type == token.IDENT &&
		p.peekToken.Type == token.LPAREN {
		callee := p.currentToken
		p.nextToken()
		p.nextToken()

		arguments := []ast.Expression{}
		if p.currentToken.Type != token.RPAREN {
			for {
				arguments = append(arguments, p.parseExpression())

				if !p.matchToken(token.COMMA) {
					break
				}
			}
		}

		if p.currentToken.Type != token.RPAREN {
			p.parseError(p.currentToken, "Expect ')' after arguments.")
			return nil
		}
		p.pushStack()
		return ast.Call{Callee: callee, Arguments: arguments}
	}

	return p.parsePrimary()
}

func (p *Parser) parsePrimary() ast.Expression {
	switch p.currentToken.Type {
	case token.INT:
		return p.parseIntegerLiteral()
	case token.CHAR:
		p.pushStack()
		return ast.CharLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
	case token.STRING:
		p.pushStack()
		return ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
	case token.IDENT:
		return p.parseVariable()
	case token.GETC, token.GETN:
		p.pushStack()
		return ast.Get{Token: p.currentToken}
	case token.TRUE, token.FALSE:
		p.pushStack()
		return ast.BooleanLiteral{Token: p.currentToken, Value: p.currentToken.Type == token.TRUE}
	}

	if p.matchToken(token.LBRACKET) {
		return p.parseArrayLiteral()
	}

	if p.matchToken(token.LPAREN) {
		expr := p.parseExpression()
		if p.currentToken.Type != token.RPAREN {
			p.parseError(p.currentToken, "Expect ')' after expression.")
			return nil
		}
		p.pushStack()
		return expr
	}

	p.parseError(p.currentToken, "Unexpect token")
	return nil
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	p.pushStack()
	value, _ := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	return ast.IntegerLiteral{Token: p.currentToken, Value: value}
}

func (p *Parser) parseVariable() ast.Expression {
	if local := p.resolveLocal(p.currentToken); local != nil {
		top := p.stackTop
		p.pushStack()
		var typ ast.VariableType
		if local.typ == LOCAL {
			typ = ast.LOCAL
		} else {
			typ = ast.ARGUMENT
		}
		return ast.Variable{
			Identifier:    p.currentToken,
			Type:          typ,
			ScopeDepth:    local.scopeDepth,
			LocalIndex:    local.localIndex,
			ArgumentIndex: local.argumentIndex,
			RelativeIndex: top,
		}
	}
	p.pushStack()

	return ast.Variable{Identifier: p.currentToken}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	elements := []ast.Expression{}
	if p.currentToken.Type != token.RPAREN {
		for {
			elements = append(elements, p.parseExpression())

			if !p.matchToken(token.COMMA) {
				break
			}
		}
	}

	if p.currentToken.Type != token.RBRACKET {
		p.parseError(p.currentToken, "Expect ']' after array literal.")
		return nil
	}

	p.pushStack()
	return ast.ArrayLiteral{Elements: elements}
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.ScanToken()
}

func (p *Parser) matchToken(types ...token.TokenType) bool {
	for _, typ := range types {
		if p.currentToken.Type == typ {
			p.nextToken()
			return true
		}

	}

	return false
}

func (p *Parser) matchPeekToken(types ...token.TokenType) bool {
	for _, typ := range types {
		if p.peekToken.Type == typ {
			return true
		}
	}

	return false

}

func (p *Parser) beginLoop() {
	p.nestedLoopCount++
}

func (p *Parser) endLoop() {
	p.nestedLoopCount--
}

func (p *Parser) isInLoop() bool {
	return p.nestedLoopCount >= 1
}

func (p *Parser) pushStack() {
	p.stackTop++
}

func (p *Parser) popStack() {
	p.stackTop--
}

func (p *Parser) beginScope() {
	p.scopes = append(p.scopes, map[string]*declaredVariable{})
}

func (p *Parser) endScope() {
	p.scopes = p.scopes[:len(p.scopes)-1]
}

func (p *Parser) declareLocalVariable(name token.Token) *declaredVariable {
	depth := len(p.scopes)
	return p.declareVariable(name, &declaredVariable{typ: LOCAL, scopeDepth: depth})
}

func (p *Parser) declareArgumentVariable(name token.Token, argumentIndex int) {
	p.declareVariable(name, &declaredVariable{typ: ARGUMENT, argumentIndex: argumentIndex})
	p.markInitializedVariable(name)
}

func (p *Parser) declareVariable(name token.Token, variable *declaredVariable) *declaredVariable {
	if len(p.scopes) == 0 {
		return nil
	}

	scope := p.scopes[len(p.scopes)-1]
	if _, ok := scope[name.Literal]; ok {
		p.parseError(p.currentToken, "ALready avariable with this name in this scope.")
		return nil
	}

	variable.localIndex = len(scope)
	scope[name.Literal] = variable

	return variable
}

func (p *Parser) markInitializedVariable(name token.Token) {
	if len(p.scopes) == 0 {
		return
	}

	scope := p.scopes[len(p.scopes)-1]
	variable, ok := scope[name.Literal]

	if !ok {
		p.parseError(name, "Not declared variable in this scope.")
		return
	}
	variable.initialized = true
}

func (p *Parser) resolveLocal(name token.Token) *declaredVariable {
	for i := len(p.scopes) - 1; i >= 0; i-- {
		if result, ok := p.scopes[i][name.Literal]; ok && result.initialized {
			return result
		}
	}

	return nil
}

func (p *Parser) parseError(tok token.Token, message string) {
	var position string
	if tok.Type == token.EOF {
		position = "at end"
	} else {
		position = "at '" + tok.Literal + "'"
	}

	p.Errors = append(
		p.Errors,
		fmt.Sprintf("%s:%d Error %s: %s\n", p.lexer.Filename, tok.Line, position, message),
	)
}

func (p *Parser) skipStatement() {
	for p.currentToken.Type != token.EOF {
		if p.currentToken.Type == token.SEMICOLON {
			return
		}

		switch p.peekToken.Type {
		case token.VAR, token.FUNC, token.RETURN, token.BREAK, token.IF,
			token.WHILE, token.PUTN, token.PUTC, token.GETN, token.GETC:
			return
		}
		p.nextToken()
	}
}
