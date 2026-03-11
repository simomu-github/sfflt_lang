package type_checker

import (
	"testing"

	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/parser"
)

func TestResolveDeclaredFunctionType(t *testing.T) {
	input := "func test1(a: int, b: char) int { return 1; } func test2() { 2; }"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	stmts := parser.ParseProgram()
	resolver := NewResolver("script", stmts)
	resolver.Resolve()

	if resolver.HadErrors() {
		t.Fatalf("Resolver error occurs.")
	}

	testFunc1, ok := resolver.DeclaredFunctions["test1"]
	if !ok {
		t.Fatalf("'test1' function is not found")
	}

	if testFunc1.Arity != 2 {
		t.Fatalf("'test1' function arity is not 2")
	}

	if testFunc1.Params[0].Name.Literal != "int" {
		t.Fatalf("'test1' param[0] type is not match")
	}

	if testFunc1.Params[1].Name.Literal != "char" {
		t.Fatalf("'test1' param[0] type is not match")
	}

	if testFunc1.Type.Name.Literal != "int" {
		t.Fatalf("'test1' return type is not match ")
	}

	testFunc2, ok := resolver.DeclaredFunctions["test2"]
	if !ok {
		t.Fatalf("'test2' function is not found")
	}

	if testFunc2.Arity != 0 {
		t.Fatalf("'test2' function arity is not 0")
	}

	if testFunc2.Type != nil {
		t.Fatalf("'test2' return type is not void")
	}
}
