package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"
const programName = "holyc"

func main() {
	fmt.Println("+------------------------------------------+")
	fmt.Println("|       HolyC Compiler                     |")
	fmt.Println("|       Go Implementation v" + version + "           |")
	fmt.Println("+------------------------------------------+")
	fmt.Println()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Parse flags
	sourceFile := ""
	tokenizeOnly := false
	parseOnly := false

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-h", "--help":
			printUsage()
			os.Exit(0)
		case "-v", "--version":
			fmt.Printf("%s version %s\n", programName, version)
			os.Exit(0)
		case "-t", "--tokens":
			tokenizeOnly = true
		case "-p", "--parse":
			parseOnly = true
		default:
			if arg[0] != '-' {
				sourceFile = arg
			}
		}
	}

	if sourceFile == "" {
		fmt.Fprintln(os.Stderr, "Error: No source file specified")
		printUsage()
		os.Exit(1)
	}

	var err error
	if tokenizeOnly {
		err = TokenizeOnly(sourceFile)
	} else if parseOnly {
		err = ParseOnly(sourceFile)
	} else {
		err = Compile(sourceFile)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "\nCompilation failed: %v\n", err)
		os.Exit(1)
	}

	if !tokenizeOnly && !parseOnly {
		fmt.Println("\nCompilation successful!")
	}
}

func printUsage() {
	fmt.Printf("Usage: %s [options] <source.hc>\n\n", programName)
	fmt.Println("A HolyC compiler written in Go.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help       Show this help message")
	fmt.Println("  -v, --version    Show version information")
	fmt.Println("  -t, --tokens     Tokenize only (show tokens)")
	fmt.Println("  -p, --parse      Parse only (show AST)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Printf("  %s program.hc           Compile program.hc\n", programName)
	fmt.Printf("  %s -t program.hc        Show tokens\n", programName)
	fmt.Printf("  %s -p program.hc        Show AST\n", programName)
	fmt.Println()
	fmt.Println("Output:")
	fmt.Println("  Generates <source>.bin with x64 machine code")
}
