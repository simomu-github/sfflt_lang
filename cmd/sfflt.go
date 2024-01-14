package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/simomu-github/sfflt_lang/compiler"
	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/parser"
)

func main() {
	if len(os.Args) >= 2 {
		Compile(os.Args[1])
	} else {
		fmt.Fprintf(os.Stderr, "Usage: sfflt [script]")
		os.Exit(64)
	}
}

func Compile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s can not read\n", path)
		return err
	}

	lexer := lexer.New(path, string(bytes))
	parser := parser.New(lexer)
	statements := parser.ParseProgram()
	if parser.HadErrors() {
		for _, err := range parser.Errors {
			fmt.Fprintf(os.Stderr, err)
		}
		return nil
	}

	compiler := compiler.New(statements)
	instructions := compiler.Compile()
	outputFilename := getFilenameWithoutExt(path) + ".fflt"
	os.WriteFile(outputFilename, []byte(strings.Join(instructions, "\n")), 0644)

	return nil
}

func getFilenameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
