package type_checker

import (
	"testing"

	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/parser"
)

func TestTypeCheckSimpleValidExpression(t *testing.T) {
	input := "1 + 2;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if type_checker.HadErrors() {
		t.Fatalf("Type check error occurs.")
	}
}

func TestTypeCheckSimpleInvalidExpression(t *testing.T) {
	input := "1 + [1];"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if !type_checker.HadErrors() {
		t.Fatalf("Type check error does not occur.")
	}
}

func TestTypeCheckMultiValidExpression(t *testing.T) {
	input := "2 * (3 + 4) == (5 + 6) * (7 + 8);"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if type_checker.HadErrors() {
		t.Fatalf("Type check error occurs.")
	}
}

func TestTypeCheckMultiInvalidExpression(t *testing.T) {
	input := "2 * (3 + 4) == false;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if !type_checker.HadErrors() {
		t.Fatalf("Type check error does not occur.")
	}
}

func TestTypeCheckWithCall(t *testing.T) {
	input := "func a() int { return 0; } a() + 1;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if type_checker.HadErrors() {
		t.Fatalf("Type check error occurs.")
	}
}

func TestTypeCheckWithInvalidCall(t *testing.T) {
	input := "func a() { 0; } a + 1;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if !type_checker.HadErrors() {
		t.Fatalf("Type check error does not occur.")
	}
}

func TestTypeCheckValidIf(t *testing.T) {
	input := "if ( true && false ) true;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if type_checker.HadErrors() {
		t.Fatalf("Type check error occurs.")
	}
}

func TestTypeCheckInvalidIf(t *testing.T) {
	input := "if ( true && 1 ) true;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if !type_checker.HadErrors() {
		t.Fatalf("Type check error does not occur.")
	}
}

func TestTypeCheckValidWhile(t *testing.T) {
	input := "while ( true && false ) true;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if type_checker.HadErrors() {
		t.Fatalf("Type check error occurs.")
	}
}

func TestTypeCheckInvalidWhile(t *testing.T) {
	input := "while ( true && 1 ) true;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if !type_checker.HadErrors() {
		t.Fatalf("Type check error does not occur.")
	}
}

func TestTypeCheckValidFor(t *testing.T) {
	input := "for (; 1 <= 0; ) true;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if type_checker.HadErrors() {
		t.Fatalf("Type check error occurs.")
	}
}

func TestTypeCheckValidReturn(t *testing.T) {
	input := "func f() int { return 0; }"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if type_checker.HadErrors() {
		t.Fatalf("Type check error occurs.")
	}
}

func TestTypeCheckInvalidReturn(t *testing.T) {
	input := "func f() { return 0; }"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if !type_checker.HadErrors() {
		t.Fatalf("Type check error does not occur.")
	}
}

func TestTypeCheckInvalidVoidReturn(t *testing.T) {
	input := "func f() int { if(false) { return 0; } }"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if !type_checker.HadErrors() {
		t.Fatalf("Type check error does not occur.")
	}
}

func TestTypeCheckValidVaiable(t *testing.T) {
	input := "var a = 0; { var a = 1; a + 1; } a + 1;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if type_checker.HadErrors() {
		t.Fatalf("Type check error occurs.")
	}
}

func TestTypeCheckInvalidVaiable(t *testing.T) {
	input := "var a = 0; { var a = 1; a + 1; } a + true;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if !type_checker.HadErrors() {
		t.Fatalf("Type check error does not occur.")
	}
}

func TestTypeCheckValidArgument(t *testing.T) {
	input := "func f(a: int) { a + 1; }"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if type_checker.HadErrors() {
		t.Fatalf("Type check error occurs.")
	}
}

func TestTypeCheckInvalidArgument(t *testing.T) {
	input := "func f(a: int) { a + true; }"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()
	type_checker := NewTypeChecker("script", stmts, resolver.DeclaredFunctions)

	type_checker.TypeCheck()

	if !type_checker.HadErrors() {
		t.Fatalf("Type check error does not occur.")
	}
}
