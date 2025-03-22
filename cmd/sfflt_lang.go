package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/simomu-github/sfflt_lang/ast"
	"github.com/simomu-github/sfflt_lang/compiler"
	"github.com/simomu-github/sfflt_lang/formatter"
	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/parser"
)

var (
	versionOpt = flag.Bool("v", false, "display version information")
	formatOpt  = flag.String("format", "64", "output code format. [oneline, pretty, (number of column)]")
)

const version = "v0.0.2"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: sfflt_lang (option) [FILE]\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if *versionOpt {
		fmt.Printf("sfflt_lang version %s\n", version)
		os.Exit(0)
	}

	if len(flag.Args()) == 1 {
		filepath := flag.Args()[0]
		stmts, err := Parse(filepath)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(Compile(filepath, stmts))
	} else {
		flag.Usage()
		os.Exit(1)
	}
}

func Parse(path string) ([]ast.Statement, error) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil, err
	}

	lexer := lexer.New(path, string(bytes))
	parser := parser.New(lexer)
	statements := parser.ParseProgram()
	if parser.HadErrors() {
		for _, err := range parser.Errors {
			fmt.Fprintf(os.Stderr, err)
		}
		return nil, errors.New("parse error.")
	}

	resolver := compiler.NewResolver(path, statements)
	resolver.Resolve()
	if resolver.HadErrors() {
		for _, err := range resolver.Errors {
			fmt.Fprintf(os.Stderr, err)
		}
		return nil, errors.New("resolve error.")
	}

	return statements, nil
}

func Compile(path string, statements []ast.Statement) int {
	compiler := compiler.New(statements)
	output, err := FormatInstructions(compiler.Compile())
	if err != nil {
		return 1
	}
	outputFilename := getFilenameWithoutExt(path) + ".fflt"
	os.WriteFile(outputFilename, []byte(output), 0644)

	return 0
}

func FormatInstructions(instructions []string) (string, error) {
	switch *formatOpt {
	case "oneline":
		return formatter.FormatOneLine(instructions), nil
	case "pretty":
		return formatter.FormatRaw(instructions), nil
	default:
		column, err := strconv.Atoi(*formatOpt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid format option. [oneline, pretty, (number of column)]\n")
			return "", errors.New("instruction format error.")
		}
		if column <= 0 {
			fmt.Fprintf(os.Stderr, "Invalid format option. [oneline, pretty, (number of column)]\n")
			return "", errors.New("instruction format error.")
		}
		return formatter.FormatSquere(instructions, column), nil
	}
}

func getFilenameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
