package type_checker

import (
	"fmt"

	"github.com/simomu-github/sfflt_lang/ast"
	"github.com/simomu-github/sfflt_lang/token"
)

type Resolver struct {
	filename          string
	statements        []ast.Statement
	DeclaredFunctions map[string]FunctionType
	Errors            []string
}

type FunctionType struct {
	Name   token.Token
	Arity  int
	Type   *DeclaredType
	Params []DeclaredType
}

func NewResolver(filename string, statements []ast.Statement) *Resolver {
	return &Resolver{
		filename:          filename,
		statements:        statements,
		DeclaredFunctions: map[string]FunctionType{},
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
}

func (r *Resolver) VisitVar(s ast.Var) {}
func (r *Resolver) VisitFunction(s ast.Function) {
	_, ok := r.DeclaredFunctions[s.Name.Literal]
	if ok {
		r.resolveError(s.Name, "function is already declared.")
	}

	params := []DeclaredType{}
	for _, param := range s.Params {
		params = append(
			params,
			DeclaredType{
				Name: param.Type.Name,
			},
		)
	}
	var retType *DeclaredType = nil
	if s.ReturnType != nil {
		retType = &DeclaredType{Name: s.ReturnType.Name}
	}
	r.DeclaredFunctions[s.Name.Literal] = FunctionType{
		Name:   s.Name,
		Arity:  len(s.Params),
		Params: params,
		Type:   retType,
	}
}
func (r *Resolver) VisitReturn(s ast.Return)                  {}
func (r *Resolver) VisitBreak(s ast.Break)                    {}
func (r *Resolver) VisitIf(s ast.If)                          {}
func (r *Resolver) VisitWhile(s ast.While)                    {}
func (r *Resolver) VisitBlock(s ast.Block)                    {}
func (r *Resolver) VisitExpression(s ast.ExpressionStatement) {}

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
