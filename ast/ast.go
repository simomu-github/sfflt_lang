package ast

import "github.com/simomu-github/sfflt_lang/token"

type ExpressionVisitor interface {
	VisitIntegerLiteral(e IntegerLiteral)
	VisitCharLiteral(e CharLiteral)
	VisitBooleanLiteral(e BooleanLiteral)
}

type Statement interface {
	statement()
}

type Expression interface {
	Visit(visitor ExpressionVisitor)
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (l IntegerLiteral) Visit(visitor ExpressionVisitor) {
	visitor.VisitIntegerLiteral(l)
}

type CharLiteral struct {
	Token token.Token
	Value string
}

func (l CharLiteral) Visit(visitor ExpressionVisitor) {
	visitor.VisitCharLiteral(l)
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (l BooleanLiteral) Visit(visitor ExpressionVisitor) {
	visitor.VisitBooleanLiteral(l)
}
