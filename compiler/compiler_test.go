package compiler

import (
	"testing"

	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/parser"
)

func TestCompile(t *testing.T) {
	input := "10 'a' true false"
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFFLFLFT",
		"FFFLLFFFFLT",
		"FFFLT",
		"FFFFT",
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction  wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}