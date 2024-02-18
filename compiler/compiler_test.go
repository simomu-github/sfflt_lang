package compiler

import (
	"testing"

	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/parser"
)

func TestCompilePrimary(t *testing.T) {
	input := "0; 10; 'a'; true; false; getn;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFFFT",
		"FTT",
		"FFFLFLFT",
		"FTT",
		"FFFLLFFFFLT",
		"FTT",
		"FFFLT",
		"FTT",
		"FFFFT",
		"FTT",
		"FFFLLLFLFFLLLFLFFLLFLLFLLLLFFFFT",
		"LTLL",
		"FFFLLLFLFFLLLFLFFLLFLLFLLLLFFFFT",
		"LLL",
		"FTT",
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileBang(t *testing.T) {
	input := "!true;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFFLT",        // push 1
		"TLFLLFLLFFFT", // jump label when zero
		"FFFFT",        // push 0
		"TFTLLFLLFFLT", // jump
		"TFFLLFLLFFFT", // mark label
		"FFFLT",        // push 1
		"TFFLLFLLFFLT", // mark label
		"FTT",          // discard
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileCall(t *testing.T) {
	input := "a();"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"TFLLLFFLLLLLFFLLFLLFFFFLT",
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileUnary(t *testing.T) {
	input := "-10;"
	lexer := lexer.New("script", input)
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
	lexer := lexer.New("script", input)
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
	lexer := lexer.New("script", input)
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
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFFLT",        // push 1
		"FFFLFT",       // push 2
		"FTL",          // swap
		"LFFL",         // sub
		"TLLLLFLLFFFT", // jump label when negative
		"FFFFT",        // push 0
		"TFTLLFLLFFLT", // jump label
		"TFFLLFLLFFFT", // mark label
		"FFFLT",        // push 1
		"TFFLLFLLFFLT", // mark label
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileComparisonWithEqual(t *testing.T) {
	input := "1 >= 2;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFFLT",         // push lhs
		"FFFLFT",        // push rhs
		"FTL",           // swap
		"LFFL",          // sub
		"FTF",           // dup
		"TLFLLFLLFFFT",  // jump label when zero
		"TLLLLFLLFFLT",  // jump label when negative
		"FFFFT",         // push 0
		"TFTLLFLLFFLFT", // jump label to end
		"TFFLLFLLFFFT",  // mark label zero
		"FTT",           // discard
		"TFFLLFLLFFLT",  // mark label negative
		"FFFLT",         // push 1
		"TFFLLFLLFFLFT", // mark label end
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileEquality(t *testing.T) {
	input := "1 != 2;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFFLT",        // push lhs
		"FFFLFT",       // push rhs
		"LFFL",         // sub
		"TLFLLFLLFFFT", // jump label when zero
		"FFFLT",        // push 1
		"TFTLLFLLFFLT", // jump label to end
		"TFFLLFLLFFFT", // mark label when zero
		"FFFFT",        // push 0
		"TFFLLFLLFFLT", // mark label end
		"FTT",          // discard
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileAssign(t *testing.T) {
	input := "var a = 1; a = 2;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFFLLFFLLLLLLFLLFLLFFFFLT",
		"FFFLT",
		"LLF",
		"FFFLLFFLLLLLLFLLFLLFFFFLT",
		"LLL",
		"FTT",
		"FFFLLFFLLLLLLFLLFLLFFFFLT",
		"FFFLFT",
		"LLF",
		"FTT",
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompilePut(t *testing.T) {
	input := "putn -1; putc 'a';"
	lexer := lexer.New("script", input)
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

func TestCompileIf(t *testing.T) {
	input := "if (true) { 1; } else { 2;}"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFFLT",        // condition
		"TLFLLFLLFFFT", // jump label when zero
		"FFFLT",        // then statement
		"FTT",          // then statement
		"TFTLLFLLFFLT", // jump label to end
		"TFFLLFLLFFFT", // mark label zero
		"FFFLFT",       // else statement
		"FTT",          // else statement
		"TFFLLFLLFFLT", //mark label end
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileWhile(t *testing.T) {
	input := "while (true) true;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"TFFLLFLLFFFT", // mark label loop
		"FFFLT",        // condition
		"TLFLLFLLFFLT", // jump label when zero
		"FFFLT",        // body statement
		"FTT",          // body statement
		"TFTLLFLLFFFT", // jump label to loop
		"TFFLLFLLFFLT", // mark label zero
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileGlobalVariable(t *testing.T) {
	input := "var a = 1; a;"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"FFFLLFFLLLLLLFLLFLLFFFFLT",
		"FFFLT",
		"LLF",
		"FFFLLFFLLLLLLFLLFLLFFFFLT",
		"LLL",
		"FTT",
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileFunction(t *testing.T) {
	input := "func a() { 1;} a();"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"TFLLLFFLLLLLFFLLFLLFFFFLT",
		"FTT",
		"TTT",
		"TFFLLFFLLLLLFFLLFLLFFFFLT",
		"FFFLT",
		"FTT",
		"FFFFT",
		"TLT",
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileReturn(t *testing.T) {
	input := "func a() { return 1; } a();"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"TFLLLFFLLLLLFFLLFLLFFFFLT",
		"FTT",
		"TTT",
		"TFFLLFFLLLLLFFLLFLLFFFFLT",
		"FFFLT",
		"TLT",
		"FFFFT",
		"TLT",
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}

func TestCompileBreak(t *testing.T) {
	input := "while(true) { while(true) { break; } break; }"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		// outer while
		"TFFLLFLLFFFT",  // mark label loop
		"FFFLT",         // condition
		"TLFLLFLLFFLLT", // jump label when zero

		// inner while
		"TFFLLFLLFFLT",  // mark label loop
		"FFFLT",         // condition
		"TLFLLFLLFFLFT", // jump label when zero
		"TFTLLFLLFFLFT", // break, jump label to end
		"TFTLLFLLFFLT",  // jump label to loop
		"TFFLLFLLFFLFT", // mark label zero

		// outer while
		"TFTLLFLLFFLLT", // break, jump label to end
		"TFTLLFLLFFFT",  // jump label to loop
		"TFFLLFLLFFLLT", // mark label zero
	}

	for i, expect := range expects {
		if instructions[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, instructions[i])
		}
	}
}
