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
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		p.nextToken()
	}

	return statements
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
	return p.parseTerm()
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
