package parser

import (
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

func (p *Parser) ParseProgram() []ast.Expression {
	expressions := []ast.Expression{}
	for p.currentToken.Type != token.EOF {
		expr := p.parseFactor()
		if expr != nil {
			expressions = append(expressions, expr)
		}
		p.nextToken()
	}

	return expressions
}

func (p *Parser) parseFactor() ast.Expression {
	expr := p.parseUnary()
	p.nextToken()
	switch p.currentToken.Type {
	case token.ASTERISK, token.SLASH:
		operator := p.currentToken
		p.nextToken()
		return ast.Factor{Left: expr, Operator: operator, Right: p.parseUnary()}
	}

	return expr
}

func (p *Parser) parseUnary() ast.Expression {
	switch p.currentToken.Type {
	case token.MINUS:
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
	case token.TRUE, token.FALSE:
		return ast.BooleanLiteral{Token: p.currentToken, Value: p.currentToken.Type == token.TRUE}
	}

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
