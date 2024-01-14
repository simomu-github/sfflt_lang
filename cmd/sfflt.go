package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/simomu-github/sfflt_lang/compiler"
	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/parser"
)

var (
	versionOpt = flag.Bool("v", false, "display version information")
)

const version = "v0.0.1"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: sfflt [FILE]\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if *versionOpt {
		fmt.Printf("sfflt version %s\n", version)
		os.Exit(0)
	}

	if len(os.Args) >= 2 {
		os.Exit(Compile(os.Args[1]))
	} else {
		flag.Usage()
		os.Exit(1)
	}
}

func Compile(path string) int {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 1
	}

	lexer := lexer.New(path, string(bytes))
	parser := parser.New(lexer)
	statements := parser.ParseProgram()
	if parser.HadErrors() {
		for _, err := range parser.Errors {
			fmt.Fprintf(os.Stderr, err)
		}
		return 1
	}

	compiler := compiler.New(statements)
	instructions := compiler.Compile()
	outputFilename := getFilenameWithoutExt(path) + ".fflt"
	os.WriteFile(outputFilename, []byte(strings.Join(instructions, "\n")), 0644)

	return 0
}

func getFilenameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
