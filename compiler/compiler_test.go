package compiler

import (
	"testing"

	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/parser"
)

func TestCompilePrimary(t *testing.T) {
	input := "0; 10; 'a'; true; false; getn;"
	instructions := compile(input, t)
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

	assertInstructions(instructions, expects, t)
}

func TestCompileBang(t *testing.T) {
	input := "!true;"
	instructions := compile(input, t)
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

	assertInstructions(instructions, expects, t)
}

func TestCompileCall(t *testing.T) {
	input := "a(1, 2, 3);"
	instructions := compile(input, t)
	expects := []string{
		"FFFLT",                     // push 1
		"FFFLFT",                    // push 2
		"FFFLLT",                    // push 3
		"TFLLLFFLLLLLFFLLFLLFFFFLT", // call sub
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileUnary(t *testing.T) {
	input := "-10;"
	instructions := compile(input, t)
	expects := []string{
		"FFLLT",
		"FFFLFLFT",
		"LFFT",
		"FTT",
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileFactor(t *testing.T) {
	input := "1 / 2 * -3 % 4;"
	instructions := compile(input, t)
	expects := []string{
		"FFFLT",   // push 1
		"FFFLFT",  // push 2
		"LFLF",    // div
		"FFLLT",   // push -1
		"FFFLLT",  // push 3
		"LFFT",    // mul
		"LFFT",    // mul
		"FFFLFFT", // push 4
		"LFLL",    // mod 4
		"FTT",
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileTerm(t *testing.T) {
	input := "(4 - 3) * (2 + 1);"
	instructions := compile(input, t)
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

	assertInstructions(instructions, expects, t)
}

func TestCompileComparison(t *testing.T) {
	input := "1 > 2;"
	instructions := compile(input, t)
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

	assertInstructions(instructions, expects, t)
}

func TestCompileComparisonWithEqual(t *testing.T) {
	input := "1 >= 2;"
	instructions := compile(input, t)
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

	assertInstructions(instructions, expects, t)
}

func TestCompileEquality(t *testing.T) {
	input := "1 != 2;"
	instructions := compile(input, t)
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

	assertInstructions(instructions, expects, t)
}

func TestCompileAnd(t *testing.T) {
	input := "true && false;"
	instructions := compile(input, t)
	expects := []string{
		"FFFLT",        // push lhs
		"TLFLLFLLFFFT", // jump label when zero
		"FFFFT",        // push rhs
		"TLFLLFLLFFFT", // jump label when zero
		"FFFLT",        // push 1
		"TFTLLFLLFFLT", // jump label to end
		"TFFLLFLLFFFT", // mark label when zero
		"FFFFT",        // push 0
		"TFFLLFLLFFLT", // mark label end
		"FTT",          // discard
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileOr(t *testing.T) {
	input := "true || false;"
	instructions := compile(input, t)
	expects := []string{
		"FFFLT",         // push lhs
		"TLFLLFLLFFFT",  // jump label when zero
		"FFFLT",         // push 1
		"TFTLLFLLFFLFT", // jump label to end

		"TFFLLFLLFFFT",  // mark label when zero
		"FFFFT",         // push rhs
		"TLFLLFLLFFLT",  // jump label when zero
		"FFFLT",         // push 1
		"TFTLLFLLFFLFT", // jump label to end

		"TFFLLFLLFFLT",  // mark label when zero
		"FFFFT",         // push 0
		"TFFLLFLLFFLFT", // mark label end
		"FTT",           // discard
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileAssign(t *testing.T) {
	input := "var a = 1; a = 2;"
	instructions := compile(input, t)
	expects := []string{
		"FFFLLFFLLLLLLFLLFLLFFFFLT", // push "a" address
		"FFFLT",                     // push 1
		"LLF",                       // store

		"FFFLLFFLLLLLLFLLFLLFFFFLT", // push "a" address
		"LLL",                       // retrieve
		"FTT",                       // discard
		"FFFLFT",                    // push 2
		"FTF",                       // dup
		"FFFLLFFLLLLLLFLLFLLFFFFLT", // push "a" address
		"FTL",                       // swap
		"LLF",                       // store
		"FTT",                       // discard
	}

	assertInstructions(instructions, expects, t)
}

func TestCompilePut(t *testing.T) {
	input := "putn -1; putc 'a';"
	instructions := compile(input, t)
	expects := []string{
		"FFLLT",
		"FFFLT",
		"LFFT",
		"LTFL",
		"FFFLLFFFFLT",
		"LTFF",
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileIf(t *testing.T) {
	input := "if (true) { 1; } else { 2;}"
	instructions := compile(input, t)
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

	assertInstructions(instructions, expects, t)
}

func TestCompileWhile(t *testing.T) {
	input := "while (true) true;"
	instructions := compile(input, t)
	expects := []string{
		"TFFLLFLLFFFT", // mark label loop
		"FFFLT",        // condition
		"TLFLLFLLFFLT", // jump label when zero
		"FFFLT",        // body statement
		"FTT",          // body statement
		"TFTLLFLLFFFT", // jump label to loop
		"TFFLLFLLFFLT", // mark label zero
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileGlobalVariable(t *testing.T) {
	input := "var a = 1; a;"
	instructions := compile(input, t)
	expects := []string{
		"FFFLLFFLLLLLLFLLFLLFFFFLT",
		"FFFLT",
		"LLF",
		"FFFLLFFLLLLLLFLLFLLFFFFLT",
		"LLL",
		"FTT",
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileArgumentVariable(t *testing.T) {
	input := "func f(a, b, c) { a + c; }"
	instructions := compile(input, t)
	expects := []string{
		"TTT",
		"TFFLLFFLLLLLFFLLFLLFFLLFT",
		"FLFFLFT", // copy 2
		"FLFFLT",  // copy 1
		"LFFF",    // add
		"FTT",     // discard
		"FFFFT",   // push 0
		"FLTFLLT", // slide 3
		"TLT",     // end sub
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileFunction(t *testing.T) {
	input := "func a() { 1;} a();"
	instructions := compile(input, t)
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

	assertInstructions(instructions, expects, t)
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

	assertInstructions(instructions, expects, t)
}

func TestCompileBreak(t *testing.T) {
	input := "while(true) { while(true) { break; } break; }"
	instructions := compile(input, t)
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

	assertInstructions(instructions, expects, t)
}

func compile(input string, t *testing.T) []string {
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	if parser.HadErrors() {
		t.Fatalf("Parse error occurred.")
	}
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	return compiler.Compile()
}

func assertInstructions(actuals []string, expects []string, t *testing.T) {
	for i, expect := range expects {
		if len(actuals) <= i {
			t.Fatalf("tests[%d] - expected instruction does not exists. expected=%q", i, expect)
		}
		if actuals[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, actuals[i])
		}
	}
}
