package main

import (
	"fmt"
	"os"
)

// Compiler is the main compiler structure
type Compiler struct {
	source     string
	lexer      *Lexer
	parser     *Parser
	codegen    *CodeGenerator
	sourceFile string
	outputFile string
	errorCnt   int
	warningCnt int
}

// NewCompiler creates a new compiler instance
func NewCompiler() *Compiler {
	return &Compiler{}
}

// Compile compiles a HolyC source file
func Compile(sourceFile string) error {
	compiler := NewCompiler()
	compiler.sourceFile = sourceFile

	// Read source file
	source, err := os.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read source file: %v", err)
	}
	compiler.source = string(source)

	fmt.Printf("Compiling %s...\n", sourceFile)
	fmt.Printf("Source size: %d bytes\n", len(source))

	// Lexical analysis
	fmt.Println("Phase 1: Lexical Analysis")
	compiler.lexer = NewLexer(compiler.source)

	// Parsing
	fmt.Println("Phase 2: Parsing")
	compiler.parser = NewParser(compiler.lexer)
	ast, err := compiler.parser.Parse()
	if err != nil {
		return fmt.Errorf("parse error: %v", err)
	}

	if compiler.parser.errorCnt > 0 {
		return fmt.Errorf("compilation failed with %d errors", compiler.parser.errorCnt)
	}

	fmt.Printf("  Found %d functions, %d classes\n",
		countFunctions(compiler.parser.functions),
		countClasses(compiler.parser.classes))

	// Code generation
	fmt.Println("Phase 3: Code Generation")
	compiler.codegen = NewCodeGenerator()
	err = compiler.codegen.Generate(ast)
	if err != nil {
		return fmt.Errorf("code generation error: %v", err)
	}

	code := compiler.codegen.GetCode()
	fmt.Printf("  Generated %d bytes of machine code\n", len(code))

	// Write output file
	outputFile := changeExtension(sourceFile, ".bin")
	err = os.WriteFile(outputFile, code, 0755)
	if err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	fmt.Printf("Output written to: %s\n", outputFile)

	// Print summary
	fmt.Println()
	fmt.Println("=== Compilation Summary ===")
	fmt.Printf("Source file:    %s (%d bytes)\n", sourceFile, len(source))
	fmt.Printf("Output file:    %s (%d bytes)\n", outputFile, len(code))
	fmt.Printf("Functions:      %d\n", countFunctions(compiler.parser.functions))
	fmt.Printf("Classes:        %d\n", countClasses(compiler.parser.classes))
	fmt.Printf("Errors:         %d\n", compiler.parser.errorCnt)
	fmt.Printf("Warnings:       %d\n", compiler.parser.warningCnt)

	// List functions
	if compiler.parser.functions != nil {
		fmt.Println()
		fmt.Println("Functions:")
		for f := compiler.parser.functions; f != nil; f = f.Next {
			extern := ""
			if (f.Flags & 0x01) != 0 {
				extern = " [extern]"
			}
			fmt.Printf("  - %s() -> %s%s\n",
				f.Name,
				typeName(f.ReturnType),
				extern)
		}
	}

	return nil
}

// countClasses counts the number of classes
func countClasses(c *Class) int {
	count := 0
	for c != nil {
		count++
		c = c.Next
	}
	return count
}

// countFunctions counts the number of functions
func countFunctions(f *Function) int {
	count := 0
	for f != nil {
		count++
		f = f.Next
	}
	return count
}

// typeName returns the name of a type
func typeName(t int) string {
	switch t {
	case RT_VOID:
		return "void"
	case RT_U8:
		return "U8"
	case RT_U16:
		return "U16"
	case RT_U32:
		return "U32"
	case RT_U64:
		return "U64"
	case RT_I8:
		return "I8"
	case RT_I16:
		return "I16"
	case RT_I32:
		return "I32"
	case RT_I64:
		return "I64"
	case RT_F64:
		return "F64"
	case RT_BOOL:
		return "Bool"
	case RT_PTR:
		return "ptr"
	default:
		return "unknown"
	}
}

// changeExtension changes the file extension
func changeExtension(filename string, newExt string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[:i] + newExt
		}
		if filename[i] == '/' || filename[i] == '\\' {
			break
		}
	}
	return filename + newExt
}

// TokenizeOnly performs lexical analysis only (for testing)
func TokenizeOnly(sourceFile string) error {
	source, err := os.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read source file: %v", err)
	}

	lexer := NewLexer(string(source))
	tokenCount := 0

	fmt.Printf("Tokenizing %s...\n\n", sourceFile)

	for {
		token := lexer.NextToken()
		tokenCount++

		if token.Type == TK_EOF {
			break
		}

		fmt.Printf("[%04d] Line %3d Col %2d: ", tokenCount, token.Line, token.Col)

		switch token.Type {
		case TK_IDENT:
			fmt.Printf("IDENT: %s\n", token.Ident)
		case TK_I64:
			fmt.Printf("NUMBER: %d\n", token.I64Val)
		case TK_F64:
			fmt.Printf("FLOAT: %f\n", token.F64Val)
		case TK_STR:
			fmt.Printf("STRING: \"%s\"\n", token.StrVal)
		case TK_CHAR_CONST:
			fmt.Printf("CHAR: %d\n", token.I64Val)
		default:
			fmt.Printf("TOKEN: %d\n", token.Type)
		}
	}

	fmt.Printf("\nTotal tokens: %d\n", tokenCount)
	return nil
}

// ParseOnly performs parsing only (for testing)
func ParseOnly(sourceFile string) error {
	source, err := os.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read source file: %v", err)
	}

	lexer := NewLexer(string(source))
	parser := NewParser(lexer)

	fmt.Printf("Parsing %s...\n\n", sourceFile)

	ast, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("parse error: %v", err)
	}

	if parser.errorCnt > 0 {
		return fmt.Errorf("parsing failed with %d errors", parser.errorCnt)
	}

	printAST(ast, 0)

	fmt.Printf("\nParsing successful!\n")
	fmt.Printf("Functions: %d\n", countFunctions(parser.functions))
	fmt.Printf("Classes: %d\n", countClasses(parser.classes))

	return nil
}

// printAST prints the AST (for debugging)
func printAST(node *ASTNode, indent int) {
	if node == nil {
		return
	}

	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "  "
	}

	fmt.Printf("%sNode Type: %d", prefix, node.Type)
	if node.Ident != "" {
		fmt.Printf(" Name: %s", node.Ident)
	}
	if node.I64Val != 0 {
		fmt.Printf(" Value: %d", node.I64Val)
	}
	fmt.Println()

	printAST(node.Child, indent+1)
	printAST(node.Sibling, indent)
	printAST(node.Next, indent)
}
