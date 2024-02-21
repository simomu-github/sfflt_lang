package compiler

import (
	"strings"
	"testing"

	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/parser"
)

func TestResolveFunctionDeclaration(t *testing.T) {
	input := "func a() { 1; } func a() { 2; }"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()

	if !resolver.HadErrors() {
		t.Fatalf("No error occurs.")
	}

	if !strings.Contains(resolver.Errors[0], "already declared") {
		t.Fatalf("Does not includes already declared error.")
	}
}

func TestResolveFunctionArity(t *testing.T) {
	input := "func f(a) { 1; } f(1,2,3);"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()

	if !resolver.HadErrors() {
		t.Fatalf("No error occurs.")
	}

	if !strings.Contains(resolver.Errors[0], "Expected 1 arguments, but got 3.") {
		t.Fatalf("Does not includes function arity error.")
	}
}
