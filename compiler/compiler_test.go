package compiler

import (
	"testing"

	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/parser"
)

func TestCompilePrimary(t *testing.T) {
	input := "10; 'a'; true; false;"
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFFLFLFT",
		"FTT",
		"FFFLLFFFFLT",
		"FTT",
		"FFFLT",
		"FTT",
		"FFFFT",
		"FTT",
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileUnary(t *testing.T) {
	input := "-10;"
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFLLT",
		"FFFLFLFT",
		"LFFT",
		"FTT",
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileFactor(t *testing.T) {
	input := "2 * -3;"
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFFLFT",
		"FFLLT",
		"FFFLLT",
		"LFFT",
		"LFFT",
		"FTT",
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileTerm(t *testing.T) {
	input := "(4 - 3) * (2 + 1);"
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFFLFFT",
		"FFFLLT",
		"LFFL",
		"FFFLFT",
		"FFFLT",
		"LFFF",
		"LFFT",
		"FTT",
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileComparison(t *testing.T) {
	input := "1 > 2;"
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFFLT",
		"FFFLFT",
		"FTL",
		"LFFL",
		"TLLLLFFFLLLLFLLFFLFFFT",
		"FFFFT",
		"TFTLLFFFLLLLFLLFFLFLFT",
		"TFFLLFFFLLLLFLLFFLFFFT",
		"FFFLT",
		"TFFLLFFFLLLLFLLFFLFLFT",
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompilePut(t *testing.T) {
	input := "putn -1; putc 'a';"
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFLLT",
		"FFFLT",
		"LFFT",
		"LTFL",
		"FFFLLFFFFLT",
		"LTFF",
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileGlobalVariable(t *testing.T) {
	input := "var a = 1; a;"
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFFLLFFLLLLLFFFFLT",
		"FFFLT",
		"LLF",
		"FFFLLFFLLLLLFFFFLT",
		"LLL",
		"FTT",
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}
