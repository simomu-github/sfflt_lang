package ast

import "github.com/simomu-github/sfflt_lang/token"

type Statement interface {
	Visit(visitor StatementVisitor)
}

type StatementVisitor interface {
	VisitVar(s Var)
	VisitFunction(f Function)
	VisitPut(s PutStatement)
	VisitReturn(s Return)
	VisitBreak(s Break)
	VisitIf(s If)
	VisitWhile(s While)
	VisitBlock(s Block)
	VisitExpression(s ExpressionStatement)
}

type Var struct {
	Identifier token.Token
	IsLocal    bool
	ScopeDepth int
	LocalIndex int
	Expression Expression
}

func (v Var) Visit(visitor StatementVisitor) {
	visitor.VisitVar(v)
}

type Function struct {
	Name   token.Token
	Params []token.Token
	Body   []Statement
}

func (f Function) Visit(visitor StatementVisitor) {
	visitor.VisitFunction(f)
}

type PutStatement struct {
	Token      token.Token
	Expression Expression
}

func (pn PutStatement) Visit(visitor StatementVisitor) {
	visitor.VisitPut(pn)
}

type Return struct {
	Value Expression
}

func (r Return) Visit(visitor StatementVisitor) {
	visitor.VisitReturn(r)
}

type Break struct {
	Token token.Token
}

func (b Break) Visit(visitor StatementVisitor) {
	visitor.VisitBreak(b)
}

type If struct {
	Condition Expression
	Then      Statement
	Else      Statement
}

func (i If) Visit(visitor StatementVisitor) {
	visitor.VisitIf(i)
}

type While struct {
	Condition Expression
	Body      Statement
}

func (w While) Visit(visitor StatementVisitor) {
	visitor.VisitWhile(w)
}

type Block struct {
	Statements []Statement
}

func (b Block) Visit(visitor StatementVisitor) {
	visitor.VisitBlock(b)
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
	VisitAssign(a Assign)
	VisitBinaryExpression(b Binary)
	VisitUnaryExpression(e Unary)
	VisitCall(e Call)
	VisitIntegerLiteral(e IntegerLiteral)
	VisitCharLiteral(e CharLiteral)
	VisitBooleanLiteral(e BooleanLiteral)
	VisitVariable(e Variable)
	VisitGet(e Get)
}

type Assign struct {
	Target     Variable
	Expression Expression
}

func (a Assign) Visit(visitor ExpressionVisitor) {
	visitor.VisitAssign(a)
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

type Call struct {
	Callee    token.Token
	Arguments []Expression
}

func (c Call) Visit(visitor ExpressionVisitor) {
	visitor.VisitCall(c)
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

type Variable struct {
	Identifier    token.Token
	Type          VariableType
	ScopeDepth    int
	LocalIndex    int
	ArgumentIndex int
	RelativeIndex int
}

const (
	LOCAL    = "LOCAL"
	ARGUMENT = "ARGUMENT"
)

type VariableType string

func (v Variable) Visit(visitor ExpressionVisitor) {
	visitor.VisitVariable(v)
}

type Get struct {
	Token token.Token
}

func (g Get) Visit(visitor ExpressionVisitor) {
	visitor.VisitGet(g)
}
