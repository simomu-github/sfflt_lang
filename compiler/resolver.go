package compiler

import (
	"fmt"

	"github.com/simomu-github/sfflt_lang/ast"
	"github.com/simomu-github/sfflt_lang/token"
)

type Resolver struct {
	filename          string
	statements        []ast.Statement
	declaredFunctions map[string]declaredFunction
	calledFunctions   map[string]calledFunction
	Errors            []string
}

type declaredFunction struct {
	name  token.Token
	arity int
}

type calledFunction struct {
	name  token.Token
	arity int
}

func NewResolver(filename string, statements []ast.Statement) *Resolver {
	return &Resolver{
		filename:          filename,
		statements:        statements,
		declaredFunctions: map[string]declaredFunction{},
		calledFunctions:   map[string]calledFunction{},
		Errors:            []string{},
	}
}

func (r *Resolver) Resolve() {
	for _, e := range r.statements {
		e.Visit(r)
	}

	if r.HadErrors() {
		return
	}

	for name, cf := range r.calledFunctions {
		if bf, ok := buildinFunctions[name]; ok {
			if cf.arity != bf.arity {
				r.resolveError(
					cf.name,
					fmt.Sprintf("Expected %d arguments, but got %d.", bf.arity, cf.arity),
				)
			}
			return
		}

		if df, ok := r.declaredFunctions[name]; ok {
			if cf.arity != df.arity {
				r.resolveError(
					cf.name,
					fmt.Sprintf("Expected %d arguments, but got %d.", df.arity, cf.arity),
				)
			}
		} else {
			r.resolveError(cf.name, "function is not declared.")
		}
	}
}

func (r *Resolver) VisitVar(s ast.Var) { s.Expression.Visit(r) }
func (r *Resolver) VisitFunction(s ast.Function) {
	_, ok := r.declaredFunctions[s.Name.Literal]
	if ok {
		r.resolveError(s.Name, "function is already declared.")
	}

	r.declaredFunctions[s.Name.Literal] = declaredFunction{
		name:  s.Name,
		arity: len(s.Params),
	}

	for _, stmt := range s.Body {
		stmt.Visit(r)
	}
}
func (r *Resolver) VisitPut(s ast.PutStatement) { s.Expression.Visit(r) }
func (r *Resolver) VisitReturn(s ast.Return)    { s.Value.Visit(r) }
func (r *Resolver) VisitBreak(s ast.Break)      {}
func (r *Resolver) VisitIf(s ast.If) {
	s.Condition.Visit(r)
	s.Then.Visit(r)
	if s.Else != nil {
		s.Else.Visit(r)
	}
}
func (r *Resolver) VisitWhile(s ast.While) {
	s.Condition.Visit(r)
	s.Body.Visit(r)
}
func (r *Resolver) VisitBlock(s ast.Block) {
	for _, stmt := range s.Statements {
		stmt.Visit(r)
	}
}
func (r *Resolver) VisitExpression(s ast.ExpressionStatement) { s.Expression.Visit(r) }

func (r *Resolver) VisitAssign(e ast.Assign)           { e.Expression.Visit(r) }
func (r *Resolver) VisitBinaryExpression(e ast.Binary) { e.Left.Visit(r); e.Right.Visit(r) }
func (r *Resolver) VisitUnaryExpression(e ast.Unary)   { e.Right.Visit(r) }
func (r *Resolver) VisitCall(e ast.Call) {
	for _, arg := range e.Arguments {
		arg.Visit(r)
	}

	r.calledFunctions[e.Callee.Literal] = calledFunction{
		name:  e.Callee,
		arity: len(e.Arguments),
	}
}
func (r *Resolver) VisitIntegerLiteral(e ast.IntegerLiteral) {}
func (r *Resolver) VisitCharLiteral(e ast.CharLiteral)       {}
func (r *Resolver) VisitStringLiteral(e ast.StringLiteral)   {}
func (r *Resolver) VisitBooleanLiteral(e ast.BooleanLiteral) {}
func (r *Resolver) VisitVariable(e ast.Variable)             {}
func (r *Resolver) VisitGet(e ast.Get)                       {}

func (r *Resolver) VisitArrayLiteral(e ast.ArrayLiteral) {
	for _, element := range e.Elements {
		element.Visit(r)
	}
}

func (r *Resolver) VisitIndex(e ast.Index) {
	e.Receiver.Visit(r)
	e.Index.Visit(r)
}

func (r *Resolver) HadErrors() bool {
	return len(r.Errors) != 0
}

func (r *Resolver) resolveError(tok token.Token, message string) {
	var position string
	if tok.Type == token.EOF {
		position = "at end"
	} else {
		position = "at '" + tok.Literal + "'"
	}

	r.Errors = append(
		r.Errors,
		fmt.Sprintf("%s:%d Error %s: %s\n", r.filename, tok.Line, position, message),
	)
}
