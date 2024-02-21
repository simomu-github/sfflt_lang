package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/simomu-github/sfflt_lang/compiler"
	"github.com/simomu-github/sfflt_lang/formatter"
	"github.com/simomu-github/sfflt_lang/lexer"
	"github.com/simomu-github/sfflt_lang/parser"
)

var (
	versionOpt = flag.Bool("v", false, "display version information")
	formatOpt  = flag.String("format", "64", "output code format. [oneline, pretty, (number of column)]")
)

const version = "v0.0.1"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: sfflt (option) [FILE]\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if *versionOpt {
		fmt.Printf("sfflt version %s\n", version)
		os.Exit(0)
	}

	if len(flag.Args()) == 1 {
		os.Exit(Compile(flag.Args()[0]))
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

	resolver := compiler.NewResolver(path, statements)
	resolver.Resolve()
	if resolver.HadErrors() {
		for _, err := range resolver.Errors {
			fmt.Fprintf(os.Stderr, err)
		}
		return 1
	}

	compiler := compiler.New(statements)
	instructions := compiler.Compile()
	outputFilename := getFilenameWithoutExt(path) + ".fflt"
	var output string
	switch *formatOpt {
	case "oneline":
		output = formatter.FormatOneLine(instructions)
	case "pretty":
		output = formatter.FormatRaw(instructions)
	default:
		column, err := strconv.Atoi(*formatOpt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid format option. [oneline, pretty, (number of column)]\n")
			return 1
		}
		if column <= 0 {
			fmt.Fprintf(os.Stderr, "Invalid format option. [oneline, pretty, (number of column)]\n")
			return 1
		}
		output = formatter.FormatSquere(instructions, column)
	}
	os.WriteFile(outputFilename, []byte(output), 0644)

	return 0
}

func getFilenameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
