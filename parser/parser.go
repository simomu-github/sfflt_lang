package parser

import (
	"fmt"
	"strconv"

	"github.com/simomu-github/sfflt_lang/ast"
	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	hasError     bool
	Errors       []string
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, Errors: []string{}}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram() []ast.Statement {
	statements := []ast.Statement{}
	for p.currentToken.Type != token.EOF {
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

	if p.currentToken.Type == token.VAR {
		return p.parseVarDeclaration()
	}

	return p.parseStatement()
}

func (p *Parser) parseVarDeclaration() ast.Statement {
	p.nextToken()
	if p.currentToken.Type != token.IDENT {
		p.parseError(p.currentToken, "Expect identifier.")
		return nil
	}
	identifier := p.currentToken
	p.nextToken()

	if p.currentToken.Type != token.ASSIGN {
		p.parseError(p.currentToken, "Expect '=' after identifier.")
		return nil
	}
	p.nextToken()

	expr := p.parseExpression()

	if p.peekToken.Type != token.SEMICOLON {
		p.parseError(p.currentToken, "Expect ';' after statement.")
		return nil
	}
	p.nextToken()

	return ast.Var{Identifier: identifier, Expression: expr}
}

func (p *Parser) parseStatement() ast.Statement {
	if p.currentToken.Type == token.PUTN ||
		p.currentToken.Type == token.PUTC {
		return p.parsePutStatement()
	}

	if p.currentToken.Type == token.IF {
		return p.parseIf()
	}

	if p.currentToken.Type == token.WHILE {
		return p.parseWhile()
	}

	if p.currentToken.Type == token.LBRACE {
		return p.parseBlock()
	}

	expr := p.parseExpression()
	if p.peekToken.Type != token.SEMICOLON {
		p.parseError(p.currentToken, "Expect ';' after statement.")
		return nil
	}
	p.nextToken()

	return ast.ExpressionStatement{Expression: expr}
}

func (p *Parser) parsePutStatement() ast.Statement {
	tok := p.currentToken
	p.nextToken()
	expr := p.parseExpression()
	if p.peekToken.Type != token.SEMICOLON {
		p.parseError(p.currentToken, "Expect ';' after statement.")
		return nil
	}
	p.nextToken()

	return ast.PutStatement{Token: tok, Expression: expr}
}

func (p *Parser) parseIf() ast.Statement {
	p.nextToken()
	if p.currentToken.Type != token.LPAREN {
		p.parseError(p.currentToken, "Expect '(' after if.")
		return nil
	}
	p.nextToken()

	condition := p.parseExpression()
	p.nextToken()

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
	p.nextToken()
	if p.currentToken.Type != token.LPAREN {
		p.parseError(p.currentToken, "Expect '(' after while.")
		return nil
	}
	p.nextToken()

	condition := p.parseExpression()
	p.nextToken()

	if p.currentToken.Type != token.RPAREN {
		p.parseError(p.currentToken, "Expect ')' after while condition.")
		return nil
	}
	p.nextToken()

	body := p.parseDeclaration()

	return ast.While{Condition: condition, Body: body}
}

func (p *Parser) parseBlock() ast.Statement {
	p.nextToken()
	stmts := []ast.Statement{}
	for p.currentToken.Type != token.RBRACE {
		if p.currentToken.Type == token.EOF {
			p.parseError(p.currentToken, "Expect '}' after statements.")
			return nil
		}
		stmts = append(stmts, p.parseDeclaration())
		p.nextToken()
	}

	return ast.Block{Statements: stmts}
}

func (p *Parser) parseExpression() ast.Expression {
	return p.parseAssign()
}

func (p *Parser) parseAssign() ast.Expression {
	expr := p.parseEquality()

	switch p.peekToken.Type {
	case token.ASSIGN:
		p.nextToken()
		p.nextToken()
		right := p.parseComparison()
		variable, ok := expr.(ast.Variable)
		if !ok {
			p.parseError(p.currentToken, "Invalid assignment target.")
			return nil
		}
		return ast.Assign{Target: variable.Identifier, Expression: right}
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
		return ast.Binary{Left: expr, Operator: operator, Right: right}
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
		return ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) parseTerm() ast.Expression {
	expr := p.parseFactor()
	switch p.peekToken.Type {
	case token.PLUS, token.MINUS:
		p.nextToken()
		operator := p.currentToken
		p.nextToken()
		right := p.parseFactor()
		return ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) parseFactor() ast.Expression {
	expr := p.parseUnary()
	switch p.peekToken.Type {
	case token.ASTERISK, token.SLASH, token.MOD:
		p.nextToken()
		operator := p.currentToken
		p.nextToken()
		right := p.parseUnary()
		return ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) parseUnary() ast.Expression {
	switch p.currentToken.Type {
	case token.MINUS, token.BANG:
		operator := p.currentToken
		p.nextToken()
		return ast.Unary{Operator: operator, Right: p.parsePrimary()}
	}

	return p.parseCall()
}

func (p *Parser) parseCall() ast.Expression {
	if p.currentToken.Type == token.IDENT &&
		p.peekToken.Type == token.LPAREN {
		callee := p.currentToken
		p.nextToken()
		p.nextToken()
		// TODO: parse arguments
		if p.currentToken.Type != token.RPAREN {
			p.parseError(p.currentToken, "Expect ')' after arguments.")
			return nil
		}
		return ast.Call{Callee: callee}
	}

	return p.parsePrimary()
}

func (p *Parser) parsePrimary() ast.Expression {
	switch p.currentToken.Type {
	case token.INT:
		return p.parseIntegerLiteral()
	case token.CHAR:
		return ast.CharLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
	case token.IDENT:
		return ast.Variable{Identifier: p.currentToken}
	case token.GETC, token.GETN:
		return ast.Get{Token: p.currentToken}
	case token.TRUE, token.FALSE:
		return ast.BooleanLiteral{Token: p.currentToken, Value: p.currentToken.Type == token.TRUE}
	case token.LPAREN:
		p.nextToken()
		expr := p.parseExpression()
		p.nextToken()
		if p.currentToken.Type != token.RPAREN {
			p.parseError(p.currentToken, "Expect ')' after expression.")
			return nil
		}
		return expr
	}

	p.parseError(p.currentToken, "Unexpect token")
	return nil
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, _ := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	return ast.IntegerLiteral{Token: p.currentToken, Value: value}
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.ScanToken()
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
	for p.currentToken.Type == token.EOF {
		if p.currentToken.Type == token.SEMICOLON {
			return
		}

		switch p.peekToken.Type {
		case token.VAR, token.IF, token.WHILE, token.PUTN, token.PUTC, token.GETN, token.GETC:
			return
		}
		p.nextToken()
	}
}
