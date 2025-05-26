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
		"FFFFT",
		"LTLL",
		"FFFFT",
		"LLL",
		"FTT",
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileArrayLiteral(t *testing.T) {
	input := "[1, 2];"
	instructions := compile(input, t)
	expects := []string{
		// allocate 6
		"FFFLFFFFFFFFFFFFFFFFFT", // push last heap allocate address
		"LLL",                    // retrieve
		"FTF",                    // dup
		"FFFLLFT",                // push 6 ( length * 2 + 2 )
		"LFFF",                   // add
		"FFFLFFFFFFFFFFFFFFFFFT", // push last heap allocate address
		"FTL",                    // swap
		"LLF",                    // store

		// setup array
		"FTF",     // dup
		"FFFLFT",  // push 2
		"LLF",     // store
		"FTF",     // dup
		"FFFLT",   // push 1
		"LFFF",    // add
		"FFFLFFT", // push 4
		"LLF",     // store

		// array[0]
		"FTF",    // dup
		"FFFLFT", // push 2
		"LFFF",   // add
		"FFFLT",  // push 1
		"LLF",    // store

		// array[1]
		"FTF",    // dup
		"FFFLLT", // push 3
		"LFFF",   // add
		"FFFLFT", // push 1
		"LLF",    // store
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileStringLiteral(t *testing.T) {
	input := "\"abc\";"
	instructions := compile(input, t)
	expects := []string{
		// allocate 8
		"FFFLFFFFFFFFFFFFFFFFFT", // push last heap allocate address
		"LLL",                    // retrieve
		"FTF",                    // dup
		"FFFLFFFT",               // push 8 ( length * 2 + 2 )
		"LFFF",                   // add
		"FFFLFFFFFFFFFFFFFFFFFT", // push last heap allocate address
		"FTL",                    // swap
		"LLF",                    // store

		// setup string
		"FTF",     // dup
		"FFFLLT",  // push 3
		"LLF",     // store
		"FTF",     // dup
		"FFFLT",   // push 1
		"LFFF",    // add
		"FFFLLFT", // push 6
		"LLF",     // store

		// str[0] = 'a'
		"FTF",         // dup
		"FFFLFT",      // push 2
		"LFFF",        // add
		"FFFLLFFFFLT", // push 'a'
		"LLF",         // store

		// str[1] = 'b'
		"FTF",         // dup
		"FFFLLT",      // push 3
		"LFFF",        // add
		"FFFLLFFFLFT", // push 'b'
		"LLF",         // store

		// str[1] = 'c'
		"FTF",         // dup
		"FFFLFFT",     // push 3
		"LFFF",        // add
		"FFFLLFFFLLT", // push 'c'
		"LLF",         // store
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileBang(t *testing.T) {
	input := "!true;"
	instructions := compile(input, t)
	expects := []string{
		"FFFLT", // push 1
		"TLFFT", // jump label when zero
		"FFFFT", // push 0
		"TFTLT", // jump
		"TFFFT", // mark label
		"FFFLT", // push 1
		"TFFLT", // mark label
		"FTT",   // discard
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileCall(t *testing.T) {
	input := "a(1, 2, 3);"
	instructions := compile(input, t)
	expects := []string{
		"FFFLT",                                  // push 1
		"FFFLFT",                                 // push 2
		"FFFLLT",                                 // push 3
		"TFLLFLLLFFLFFFFFFLLFFFFLFLFFLFFLFLLFFT", // call sub
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileIndex(t *testing.T) {
	input := "[1][0];"
	instructions := compile(input, t)
	expects := []string{
		// allocate 4
		"FFFLFFFFFFFFFFFFFFFFFT", // push last heap allocate address
		"LLL",                    // retrieve
		"FTF",                    // dup
		"FFFLFFT",                // push 4 ( length * 2 + 2 )
		"LFFF",                   // add
		"FFFLFFFFFFFFFFFFFFFFFT", // push last heap allocate address
		"FTL",                    // swap
		"LLF",                    // store

		// setup array
		"FTF",    // dup
		"FFFLT",  // push 1
		"LLF",    // store
		"FTF",    // dup
		"FFFLT",  // push 1
		"LFFF",   // add
		"FFFLFT", // push 2
		"LLF",    // store

		// assign array[0]
		"FTF",    // dup
		"FFFLFT", // push 2
		"LFFF",   // add
		"FFFLT",  // push 1
		"LLF",    // store

		// fetch array[0]
		"FFFFT",  // push 0
		"FFFLFT", // push 2
		"LFFF",   // add
		"LFFF",   // add
		"LLL",    // retrieve
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
		"FFFLT",  // push 1
		"FFFLFT", // push 2
		"FTL",    // swap
		"LFFL",   // sub
		"TLLFT",  // jump label when negative
		"FFFFT",  // push 0
		"TFTLT",  // jump label
		"TFFFT",  // mark label
		"FFFLT",  // push 1
		"TFFLT",  // mark label
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileComparisonWithEqual(t *testing.T) {
	input := "1 >= 2;"
	instructions := compile(input, t)
	expects := []string{
		"FFFLT",  // push lhs
		"FFFLFT", // push rhs
		"FTL",    // swap
		"LFFL",   // sub
		"FTF",    // dup
		"TLFFT",  // jump label when zero
		"TLLLT",  // jump label when negative
		"FFFFT",  // push 0
		"TFTLFT", // jump label to end
		"TFFFT",  // mark label zero
		"FTT",    // discard
		"TFFLT",  // mark label negative
		"FFFLT",  // push 1
		"TFFLFT", // mark label end
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileEquality(t *testing.T) {
	input := "1 != 2;"
	instructions := compile(input, t)
	expects := []string{
		"FFFLT",  // push lhs
		"FFFLFT", // push rhs
		"LFFL",   // sub
		"TLFFT",  // jump label when zero
		"FFFLT",  // push 1
		"TFTLT",  // jump label to end
		"TFFFT",  // mark label when zero
		"FFFFT",  // push 0
		"TFFLT",  // mark label end
		"FTT",    // discard
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileAnd(t *testing.T) {
	input := "true && false;"
	instructions := compile(input, t)
	expects := []string{
		"FFFLT", // push lhs
		"TLFFT", // jump label when zero
		"FFFFT", // push rhs
		"TLFFT", // jump label when zero
		"FFFLT", // push 1
		"TFTLT", // jump label to end
		"TFFFT", // mark label when zero
		"FFFFT", // push 0
		"TFFLT", // mark label end
		"FTT",   // discard
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileOr(t *testing.T) {
	input := "true || false;"
	instructions := compile(input, t)
	expects := []string{
		"FFFLT",  // push lhs
		"TLFFT",  // jump label when zero
		"FFFLT",  // push 1
		"TFTLFT", // jump label to end

		"TFFFT",  // mark label when zero
		"FFFFT",  // push rhs
		"TLFLT",  // jump label when zero
		"FFFLT",  // push 1
		"TFTLFT", // jump label to end

		"TFFLT",  // mark label when zero
		"FFFFT",  // push 0
		"TFFLFT", // mark label end
		"FTT",    // discard
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileGlobalVariableAssign(t *testing.T) {
	input := "var a = 1; a = 2;"
	instructions := compile(input, t)
	expects := []string{
		"FFFLFLLLFFLFFFFFFLLFFFFLFLFFLFFLFLLFFT", // push "a" address
		"FFFLT",                                  // push 1
		"LLF",                                    // store

		"FFFLFLLLFFLFFFFFFLLFFFFLFLFFLFFLFLLFFT", // push "a" address
		"FTF",                                    // dup
		"FFFLFT",                                 // push 2
		"LLF",                                    // store

		"LLL", // retrieve
		"FTT", // discard
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileAssignLocalVariable(t *testing.T) {
	input := "{ var a = 1; a = 2; }"
	instructions := compile(input, t)
	expects := []string{

		"FFFLFFFFFFFFFFFFFFFFFFFFFFFFFLFFFFFFFFT", // push local variable addr (scope 1, index 0)
		"FFFLT", // push 1
		"LLF",   // store

		"FFFLFFFFFFFFFFFFFFFFFFFFFFFFFLFFFFFFFFT", // push local variable addr (scope 1, index 0)
		"FTF",    // dup
		"FFFLFT", // push 2
		"LLF",    // store

		"LLL", // retrieve
		"FTT", // discard
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileAssignIndex(t *testing.T) {
	input := "[1][0] = 2;"
	instructions := compile(input, t)
	expects := []string{
		// array literal
		// allocate 4
		"FFFLFFFFFFFFFFFFFFFFFT", // push last heap allocate address
		"LLL",                    // retrieve
		"FTF",                    // dup
		"FFFLFFT",                // push 4 ( length * 2 + 2 )
		"LFFF",                   // add
		"FFFLFFFFFFFFFFFFFFFFFT", // push last heap allocate address
		"FTL",                    // swap
		"LLF",                    // store

		// setup array
		"FTF",    // dup
		"FFFLT",  // push 1
		"LLF",    // store
		"FTF",    // dup
		"FFFLT",  // push 1
		"LFFF",   // add
		"FFFLFT", // push 2
		"LLF",    // store
		"FTF",    // dup
		"FFFLFT", // push 2
		"LFFF",   // add
		"FFFLT",  // push 1
		"LLF",    // store

		// fetch array[0] address
		"FFFFT",  // push 0
		"FFFLFT", // push 2
		"LFFF",   // add
		"LFFF",   // add

		// assign to array[0]
		"FTF",    // dup
		"FFFLFT", // push 2
		"LLF",    // store

		"LLL", // retrieve
		"FTT", // discard

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
		"FFFLT",  // condition
		"TLFFT",  // jump label when zero
		"FFFLT",  // then statement
		"FTT",    // then statement
		"TFTLT",  // jump label to end
		"TFFFT",  // mark label zero
		"FFFLFT", // else statement
		"FTT",    // else statement
		"TFFLT",  //mark label end
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileWhile(t *testing.T) {
	input := "while (true) true;"
	instructions := compile(input, t)
	expects := []string{
		"TFFFT", // mark label loop
		"FFFLT", // condition
		"TLFLT", // jump label when zero
		"FFFLT", // body statement
		"FTT",   // body statement
		"TFTFT", // jump label to loop
		"TFFLT", // mark label zero
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileGlobalVariable(t *testing.T) {
	input := "var a = 1; a;"
	instructions := compile(input, t)
	expects := []string{
		"FFFLFLLLFFLFFFFFFLLFFFFLFLFFLFFLFLLFFT",
		"FFFLT",
		"LLF",
		"FFFLFLLLFFLFFFFFFLLFFFFLFLFFLFFLFLLFFT",
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
		"TFFLFLLLFFFLLFFFFLLFFFFLFFLLLLFFLLFFLT",
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

func TestCompileLocalVariable(t *testing.T) {
	input := "{ var a = 1; { var b = 2; var c = 3;  a; b; }}"
	instructions := compile(input, t)
	expects := []string{
		"FFFLFFFFFFFFFFFFFFFFFFFFFFFFFLFFFFFFFFT", // push local variable addr (scope 1, index 0)
		"FFFLT", // push 1
		"LLF",   // store
		"FFFLFFFFFFFFFFFFFFFFFFFFFFFFLFFFFFFFFFT", // push local variable addr (scope 2, index 0)
		"FFFLFT", // push 2
		"LLF",    // store
		"FFFLFFFFFFFFFFFFFFFFFFFFFFFFLFFFFFFFFLT", // push local variable addr (scope 2, index 1)
		"FFFLLT", // push 3
		"LLF",    // store
		"FFFLFFFFFFFFFFFFFFFFFFFFFFFFFLFFFFFFFFT", // push local variable addr (scope 1, index 0)
		"LLL", // retrieve
		"FTT", // discard
		"FFFLFFFFFFFFFFFFFFFFFFFFFFFFLFFFFFFFFFT", // push local variable addr (scope 2, index 0)
		"LLL", // retrieve
		"FTT", // discard
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileFunction(t *testing.T) {
	input := "func a() { 1;} a();"
	instructions := compile(input, t)
	expects := []string{
		"TFLLFLLLFFLFFFFFFLLFFFFLFLFFLFFLFLLFFT",
		"FTT",
		"TTT",
		"TFFLFLLLFFLFFFFFFLLFFFFLFLFFLFFLFLLFFT",
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
		"TFLLFLLLFFLFFFFFFLLFFFFLFLFFLFFLFLLFFT",
		"FTT",
		"TTT",
		"TFFLFLLLFFLFFFFFFLLFFFFLFLFFLFFLFLLFFT",
		"FFFLT",
		"TLT",
		"FFFFT",
		"TLT",
	}

	assertInstructions(instructions, expects, t)
}

func TestCompileEmptyReturn(t *testing.T) {
	input := "func a() { return; } a();"
	lexer := lexer.New("script", input)
	parser := parser.New(lexer)
	exprs := parser.ParseProgram()
	compiler := New(exprs)

	instructions := compiler.Compile()
	expects := []string{
		"TFLLFLLLFFLFFFFFFLLFFFFLFLFFLFFLFLLFFT",
		"FTT",
		"TTT",
		"TFFLFLLLFFLFFFFFFLLFFFFLFLFFLFFLFLLFFT",
		"FFFFT",
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
		"TFFFT",  // mark label loop
		"FFFLT",  // condition
		"TLFLLT", // jump label when zero

		// inner while
		"TFFLT",  // mark label loop
		"FFFLT",  // condition
		"TLFLFT", // jump label when zero
		"TFTLFT", // break, jump label to end
		"TFTLT",  // jump label to loop
		"TFFLFT", // mark label zero

		// outer while
		"TFTLLT", // break, jump label to end
		"TFTFT",  // jump label to loop
		"TFFLLT", // mark label zero
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
	// init vm heap instructions
	if len(actuals) < 3 {
		t.Fatal("Initialize vm heap instructions do not exists")
	}
	initExpects := []string{
		"FFFLFFFFFFFFFFFFFFFFFT",                  // push last heap allocate address
		"FFFLLFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFT", // push init heap address
		"LLF", // store
	}
	for i, expect := range initExpects {
		if actuals[i] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, actuals[i])
		}
	}

	for i, expect := range expects {
		if len(actuals) <= i+3 {
			t.Fatalf("tests[%d] - expected instruction does not exists. expected=%q", i, expect)
		}
		if actuals[i+3] != expect {
			t.Fatalf("tests[%d] - instruction wrong. expected=%q, got=%q", i, expect, actuals[i+3])
		}
	}
}
