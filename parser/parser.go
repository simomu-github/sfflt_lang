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
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer}
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

func (p *Parser) parseDeclaration() ast.Statement {
	if p.currentToken.Type == token.VAR {
		return p.parseVarDeclaration()
	}

	return p.parseStatement()
}

func (p *Parser) parseVarDeclaration() ast.Statement {
	p.nextToken()
	if p.currentToken.Type != token.IDENT {
		panic("Parser error")
	}
	identifier := p.currentToken
	p.nextToken()

	if p.currentToken.Type != token.ASSIGN {
		panic("Parser error")
	}
	p.nextToken()

	expr := p.parseExpression()

	if p.peekToken.Type != token.SEMICOLON {
		panic("Parser error")
	}
	p.nextToken()

	return ast.Var{Identifier: identifier, Expression: expr}
}

func (p *Parser) parseStatement() ast.Statement {
	if p.currentToken.Type == token.PUTN ||
		p.currentToken.Type == token.PUTC {
		return p.parsePutStatement()
	}

	expr := p.parseExpression()
	if p.peekToken.Type != token.SEMICOLON {
		panic("Parser error")
	}
	p.nextToken()

	return ast.ExpressionStatement{Expression: expr}
}

func (p *Parser) parsePutStatement() ast.Statement {
	tok := p.currentToken
	p.nextToken()
	expr := p.parseExpression()
	if p.peekToken.Type != token.SEMICOLON {
		panic("Parser error")
	}
	p.nextToken()

	return ast.PutStatement{Token: tok, Expression: expr}
}

func (p *Parser) parseExpression() ast.Expression {
	return p.parseEquality()
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
	case token.ASTERISK, token.SLASH:
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
	case token.TRUE, token.FALSE:
		return ast.BooleanLiteral{Token: p.currentToken, Value: p.currentToken.Type == token.TRUE}
	case token.LPAREN:
		p.nextToken()
		expr := p.parseExpression()
		p.nextToken()
		if p.currentToken.Type != token.RPAREN {
			// TODO: parser error
			panic(fmt.Sprintf("Parser Error at %s", p.currentToken.Type))
		}
		return expr
	}

	// TODO: parser error
	panic("Parser Error")
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, _ := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	return ast.IntegerLiteral{Token: p.currentToken, Value: value}
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.ScanToken()
}
