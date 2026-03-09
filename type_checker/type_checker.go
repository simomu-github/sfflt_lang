package type_checker

import (
	"fmt"

	"github.com/simomu-github/sfflt_lang/ast"
	"github.com/simomu-github/sfflt_lang/token"
)

const (
	GLOBAL                 = "GLOBAL"
	WITHIN_FUNCTION        = "WITHIN_FUNCTION"
	WITHIN_FUNCTION_BRANCH = "WITHIN_FUNCTION_BRANCH"
)

type TypeCheckState string

type TypeChecker struct {
	filename                       string
	statements                     []ast.Statement
	declaredFunctions              map[string]FunctionType
	typeEnvironment                *TypeEnvironment
	currentFunctionType            Type
	inFunctionBranchlessReturnType Type
	typeCheckState                 TypeCheckState
	Errors                         []string
}

func NewTypeChecker(filename string, statements []ast.Statement, declaredFunctions map[string]FunctionType) *TypeChecker {
	return &TypeChecker{
		filename:                       filename,
		statements:                     statements,
		declaredFunctions:              declaredFunctions,
		typeEnvironment:                NewTypeEnvironment(nil),
		currentFunctionType:            Type{Tag: "invalid"},
		inFunctionBranchlessReturnType: Type{Tag: "invalid"},
		typeCheckState:                 GLOBAL,
		Errors:                         []string{},
	}
}

func (t *TypeChecker) TypeCheck() {
	for _, e := range t.statements {
		e.Visit(t)
	}
}

func (t *TypeChecker) VisitVar(s ast.Var) {
	typ := t.typeCheck(s.Expression)
	t.typeEnvironment.AddVariableType(s.Identifier.Literal, typ)
}
func (t *TypeChecker) VisitFunction(s ast.Function) {
	t.beginFunction()
	if s.ReturnType != nil {
		t.currentFunctionType = Type{Tag: s.ReturnType.Name.Literal}
	} else {
		t.currentFunctionType = Type{Tag: "void"}
	}

	for _, param := range s.Params {
		t.typeEnvironment.AddVariableType(
			param.Name.Literal,
			Type{Tag: param.Type.Name.Literal},
		)
	}

	for _, stmt := range s.Body {
		stmt.Visit(t)
	}

	returnTypeMatch :=
		typeEq(t.currentFunctionType, t.inFunctionBranchlessReturnType) ||
			(t.currentFunctionType.Tag == "void" && t.inFunctionBranchlessReturnType.Tag == "void")
	if !returnTypeMatch {
		t.typeCheckError(s.Name, "return type does not match.")
		return
	}
	t.currentFunctionType = Type{Tag: "invalid"}
	t.endFunction()
}
func (t *TypeChecker) VisitReturn(s ast.Return) {
	if s.Value != nil {
		returnType := t.typeCheck(s.Value)
		if !typeEq(returnType, t.currentFunctionType) {
			t.typeCheckError(s.Token, "return type does not match.")
			return
		}
		t.registerReturnType(returnType)
	} else {
		returnType := Type{Tag: "void"}
		if !typeEq(t.currentFunctionType, returnType) {
			t.typeCheckError(s.Token, "return type does not match.")
			return
		}
		t.registerReturnType(returnType)
	}
}
func (t *TypeChecker) VisitBreak(s ast.Break) {}
func (t *TypeChecker) VisitIf(s ast.If) {
	conditionType := t.typeCheck(s.Condition)
	if !typeEq(conditionType, Type{Tag: "boolean"}) {
		t.typeCheckError(s.Token, "condition expression does not bool.")
		return
	}

	t.beginBranch()
	s.Then.Visit(t)
	if s.Else != nil {
		s.Else.Visit(t)
	}
	t.endBranch()
}
func (t *TypeChecker) VisitWhile(s ast.While) {
	conditionType := t.typeCheck(s.Condition)
	if !typeEq(conditionType, Type{Tag: "boolean"}) {
		t.typeCheckError(s.Token, "condition expression does not bool.")
		return
	}
	t.beginBranch()
	s.Body.Visit(t)
	t.endBranch()
}
func (t *TypeChecker) VisitBlock(s ast.Block) {
	t.pushScope()
	for _, stmt := range s.Statements {
		stmt.Visit(t)
	}
	t.popScope()
}

func (t *TypeChecker) VisitExpression(s ast.ExpressionStatement) {
	t.typeCheck(s.Expression)
}

func (t *TypeChecker) HadErrors() bool {
	return len(t.Errors) != 0
}

func (t *TypeChecker) typeCheckError(tok token.Token, message string) {
	var position string
	if tok.Type == token.EOF {
		position = "at end"
	} else {
		position = "at '" + tok.Literal + "'"
	}

	t.Errors = append(
		t.Errors,
		fmt.Sprintf("%s:%d Error %s: %s\n", t.filename, tok.Line, position, message),
	)
}

func (t *TypeChecker) pushScope() {
	prev := t.typeEnvironment
	t.typeEnvironment = NewTypeEnvironment(prev)
}

func (t *TypeChecker) popScope() {
	t.typeEnvironment = t.typeEnvironment.parent
}

func (t *TypeChecker) beginFunction() {
	t.pushScope()
	t.typeCheckState = WITHIN_FUNCTION
	t.inFunctionBranchlessReturnType = Type{Tag: "void"}
}

func (t *TypeChecker) endFunction() {
	t.popScope()
	t.typeCheckState = GLOBAL
	t.inFunctionBranchlessReturnType = Type{Tag: "invalid"}
}

func (t *TypeChecker) registerReturnType(typ Type) {
	if t.typeCheckState == WITHIN_FUNCTION {
		t.inFunctionBranchlessReturnType = typ
	}
}

func (t *TypeChecker) beginBranch() {
	if t.typeCheckState == WITHIN_FUNCTION {
		t.typeCheckState = WITHIN_FUNCTION_BRANCH
	}
}

func (t *TypeChecker) endBranch() {
	if t.typeCheckState == WITHIN_FUNCTION_BRANCH {
		t.typeCheckState = WITHIN_FUNCTION
	}
}

func (t *TypeChecker) typeCheck(e ast.Expression) Type {
	switch expr := e.(type) {
	case ast.CharLiteral:
		return Type{Tag: "char"}
	case ast.IntegerLiteral:
		return Type{Tag: "int"}
	case ast.BooleanLiteral:
		return Type{Tag: "boolean"}
	case ast.ArrayLiteral:
		return Type{Tag: "array"}
	case ast.StringLiteral:
		return Type{Tag: "string"}
	case ast.Variable:
		variableType := t.typeEnvironment.FindVariableType(expr.Identifier.Literal)
		if variableType == nil {
			t.typeCheckError(expr.Identifier, "Does not declared.")
			return Type{Tag: "invalid"}
		}
		return *variableType
	case ast.Call:
		{
			callee, ok := t.declaredFunctions[expr.Callee.Literal]
			if !ok {
				t.typeCheckError(expr.Callee, "Does not declared.")
				return Type{Tag: "invalid"}
			}
			return convertoToType(callee.Type)
		}
	case ast.Unary:
		{
			rht := t.typeCheck(expr.Right)
			if rht.Tag == "int" || rht.Tag == "char" {
				return Type{Tag: rht.Tag}
			} else {
				t.typeCheckError(expr.Operator, "It operator can not use this type.")
				return Type{Tag: "invalid"}
			}
		}
	case ast.Binary:
		{
			lht := t.typeCheck(expr.Left)
			rht := t.typeCheck(expr.Right)
			switch expr.Operator.Type {
			case token.AND, token.OR:
				if !typeEq(lht, Type{Tag: "boolean"}) {
					t.typeCheckError(expr.Operator, "It operator can not use this type.")
					return Type{Tag: "invalid"}
				}
				if !typeEq(rht, Type{Tag: "boolean"}) {
					t.typeCheckError(expr.Operator, "It operator can not use this type.")
					return Type{Tag: "invalid"}
				}
				return Type{Tag: "boolean"}
			case token.EQ, token.NOT_EQ, token.GT, token.GTEQ, token.LT, token.LTEQ:
				if !typeEq(lht, rht) {
					t.typeCheckError(expr.Operator, "It operator can not use this type")
					return Type{Tag: "invalid"}
				}
				return Type{Tag: "boolean"}
			case token.PLUS, token.MINUS, token.ASTERISK, token.SLASH, token.MOD:
				if !typeEq(lht, rht) {
					t.typeCheckError(expr.Operator, "It operator can not use this type")
					return Type{Tag: "invalid"}
				}
				return Type{Tag: "int"}
			default:
				panic("Not implemented")
			}

		}
	}
	return Type{Tag: "invalid"}
}

func typeEq(type1 Type, type2 Type) bool {
	switch type1.Tag {
	case "char":
		return type2.Tag == "int" || type2.Tag == "char"
	case "int":
		return type2.Tag == "int" || type2.Tag == "char"
	case "array":
		return type2.Tag == "array"
	case "string":
		return type2.Tag == "string"
	case "boolean":
		return type2.Tag == "boolean"
	case "void":
		return false
	}

	return false
}

func convertoToType(declaredType *DeclaredType) Type {
	if declaredType == nil {
		return Type{Tag: "void"}
	}

	return Type{Tag: declaredType.Name.Literal}
}

func isLogicalExpression(tok token.Token) bool {
	switch tok.Type {
	case token.AND, token.OR:
		return true
	default:
		return false
	}
}

func isCompareExpression(tok token.Token) bool {
	switch tok.Type {
	case token.EQ, token.NOT_EQ, token.GT, token.GTEQ, token.LT, token.LTEQ:
		return true
	default:
		return false
	}
}

type TypeEnvironment struct {
	parent    *TypeEnvironment
	variables map[string]Type
}

func NewTypeEnvironment(parent *TypeEnvironment) *TypeEnvironment {
	return &TypeEnvironment{
		parent:    parent,
		variables: map[string]Type{},
	}
}

func (te *TypeEnvironment) FindVariableType(name string) *Type {
	if v, ok := te.variables[name]; ok {
		return &v
	}

	if te.parent == nil {
		return nil
	} else {
		return te.parent.FindVariableType(name)
	}
}

func (te *TypeEnvironment) AddVariableType(name string, typ Type) {
	te.variables[name] = typ
}
