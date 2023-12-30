package ast

import "github.com/simomu-github/sfflt_lang/token"

type ExpressionVisitor interface {
	VisitFactorExpression(f Factor)
	VisitUnaryExpression(e Unary)
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

type Factor struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (f Factor) Visit(visitor ExpressionVisitor) {
	visitor.VisitFactorExpression(f)
}

type Unary struct {
	Operator token.Token
	Right    Expression
}

func (u Unary) Visit(visitor ExpressionVisitor) {
	visitor.VisitUnaryExpression(u)
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
