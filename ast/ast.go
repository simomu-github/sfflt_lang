package ast

import "github.com/simomu-github/sfflt_lang/token"

type Statement interface {
	Visit(visitor StatementVisitor)
}

type StatementVisitor interface {
	VisitVar(s Var)
	VisitPut(s PutStatement)
	VisitExpression(s ExpressionStatement)
}

type Var struct {
	Identifier token.Token
	Expression Expression
}

func (v Var) Visit(visitor StatementVisitor) {
	visitor.VisitVar(v)
}

type PutStatement struct {
	Token      token.Token
	Expression Expression
}

func (pn PutStatement) Visit(visitor StatementVisitor) {
	visitor.VisitPut(pn)
}

type ExpressionStatement struct {
	Expression Expression
}

func (e ExpressionStatement) Visit(visitor StatementVisitor) {
	visitor.VisitExpression(e)
}

type Expression interface {
	Visit(visitor ExpressionVisitor)
}

type ExpressionVisitor interface {
	VisitBinaryExpression(b Binary)
	VisitUnaryExpression(e Unary)
	VisitIntegerLiteral(e IntegerLiteral)
	VisitCharLiteral(e CharLiteral)
	VisitBooleanLiteral(e BooleanLiteral)
}

type Binary struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (b Binary) Visit(visitor ExpressionVisitor) {
	visitor.VisitBinaryExpression(b)
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
