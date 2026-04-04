package main

import (
	"fmt"
)

// Symbol types
const (
	SYM_VAR = iota + 1
	SYM_FUN
	SYM_CLASS
	SYM_PARAM
	SYM_LOCAL
)

// AST node types
const (
	AST_NULL = iota
	AST_FUN_DECL
	AST_FUN_CALL
	AST_VAR_DECL
	AST_ASSIGN
	AST_BINARY_OP
	AST_UNARY_OP
	AST_IF
	AST_WHILE
	AST_FOR
	AST_DO
	AST_RETURN
	AST_BREAK
	AST_CONTINUE
	AST_BLOCK
	AST_IDENT
	AST_LITERAL
	AST_STRING
	AST_MEMBER_ACCESS
	AST_ARRAY_ACCESS
	AST_CLASS_DECL
	AST_EXTERN_DECL
	AST_SWITCH
	AST_CASE
	AST_GOTO
	AST_LABEL
)

// Type codes
const (
	RT_VOID = iota
	RT_U8
	RT_U16
	RT_U32
	RT_U64
	RT_I8
	RT_I16
	RT_I32
	RT_I64
	RT_F64
	RT_PTR
	RT_CLASS
	RT_FUN
	RT_BOOL
)

// ASTNode represents a node in the abstract syntax tree
type ASTNode struct {
	Type      int
	DataType  int
	I64Val    int64
	F64Val    float64
	StrVal    string
	Ident     string
	Child     *ASTNode
	Sibling   *ASTNode
	Next      *ASTNode
	Sym       *Symbol
	Fun       *Function
	ClassDef  *Class
	Body      *ASTNode  // For FOR loop body
	Line      int
}

// Symbol represents a symbol table entry
type Symbol struct {
	Name     string
	Type     int // SYM_VAR, SYM_FUN, SYM_CLASS, SYM_PARAM, SYM_LOCAL
	DataType int
	Offset   int
	Size     int
	Flags    int
	UseCnt   int
	Next     *Symbol
	Members  *Symbol
}

// Class represents a class definition
type Class struct {
	Name      string
	Size      int
	MemberCnt int
	Members   *Symbol
	BaseClass *Class
	Next      *Class
}

// Function represents a function definition
type Function struct {
	Name       string
	ReturnType int
	ParamCnt   int
	LocalCnt   int
	Size       int
	Flags      int
	ExeAddr    int
	Params     *Symbol
	Locals     *Symbol
	Next       *Function
	CodeStart  int
	CodeEnd    int
}

// Parser parses tokens into an AST
type Parser struct {
	lexer       *Lexer
	currentToken Token
	astRoot     *ASTNode
	symTable    *Symbol
	globalSyms  *Symbol
	localSyms   *Symbol
	classes     *Class
	functions   *Function
	currentFun  *Function
	errorCnt    int
	warningCnt  int
}

// NewParser creates a new parser
func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		lexer: lexer,
	}
}

// Error reports an error
func (p *Parser) Error(line int, format string, args ...interface{}) {
	p.errorCnt++
	fmt.Printf("Error at line %d: ", line)
	fmt.Printf(format+"\n", args...)
}

// Warning reports a warning
func (p *Parser) Warning(line int, format string, args ...interface{}) {
	p.warningCnt++
	fmt.Printf("Warning at line %d: ", line)
	fmt.Printf(format+"\n", args...)
}

// Parse starts parsing and returns the AST
func (p *Parser) Parse() (*ASTNode, error) {
	p.lexer.NextToken()
	p.astRoot = p.parseTranslationUnit()
	return p.astRoot, nil
}

// parseTranslationUnit parses the top-level declarations
func (p *Parser) parseTranslationUnit() *ASTNode {
	node := &ASTNode{Type: AST_BLOCK}
	var lastDecl *ASTNode

	for p.lexer.token.Type != TK_EOF {
		// Handle preprocessor directives
		if p.lexer.token.Type == KW_INCLUDE {
			p.lexer.NextToken()
			if p.lexer.token.Type == TK_STR {
				p.lexer.NextToken()
			}
			continue
		}

		if p.lexer.token.Type == KW_DEFINE ||
			p.lexer.token.Type == KW_IFDEF ||
			p.lexer.token.Type == KW_IFNDEF {
			// Skip until endif
			depth := 1
			for depth > 0 && p.lexer.token.Type != TK_EOF {
				if p.lexer.token.Type == KW_IFDEF || p.lexer.token.Type == KW_IFNDEF {
					depth++
				} else if p.lexer.token.Type == KW_ENDIF {
					depth--
				}
				p.lexer.NextToken()
			}
			continue
		}

		var decl *ASTNode

		if p.lexer.token.Type == KW_CLASS {
			decl = p.parseClassDeclaration()
		} else if p.lexer.token.Type == KW_EXTERN {
			decl = p.parseExternDeclaration()
		} else {
			decl = p.parseFunctionDeclaration()
		}

		if decl != nil {
			if lastDecl != nil {
				lastDecl.Next = decl
			} else {
				node.Child = decl
			}
			lastDecl = decl
		}
	}

	return node
}

// parseFunctionDeclaration parses a function declaration
func (p *Parser) parseFunctionDeclaration() *ASTNode {
	node := &ASTNode{Type: AST_FUN_DECL, Line: p.lexer.token.Line}

	// Parse return type
	returnType := p.parseType()

	// Parse function name
	if p.lexer.token.Type != TK_IDENT {
		p.Error(p.lexer.token.Line, "Expected function name")
		return node
	}
	name := p.lexer.token.Ident
	node.Ident = name
	p.lexer.NextToken()

	// Create function
	f := &Function{
		Name:       name,
		ReturnType: returnType,
	}
	p.addFunction(f)
	node.Fun = f

	// Parse parameters
	if err := p.lexer.Match(TK_LPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	var lastParam *Symbol
	if p.lexer.token.Type != TK_RPAREN {
		for {
			paramType := p.parseType()

			if p.lexer.token.Type != TK_IDENT {
				break
			}

			param := &Symbol{
				Name:     p.lexer.token.Ident,
				Type:     4, // SYM_PARAM
				DataType: paramType,
				Offset:   16 + f.ParamCnt*8,
				Size:     typeSize(paramType),
			}

			if lastParam != nil {
				lastParam.Next = param
			} else {
				f.Params = param
			}
			lastParam = param
			f.ParamCnt++

			p.lexer.NextToken()

			if p.lexer.token.Type != TK_COMMA {
				break
			}
			p.lexer.NextToken()
		}
	}

	if err := p.lexer.Match(TK_RPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	// Check for extern (no body)
	if p.lexer.token.Type == TK_SEMICOLON {
		p.lexer.NextToken()
		f.Flags |= 0x01 // Extern flag
		return node
	}

	// Parse function body
	p.currentFun = f
	p.localSyms = nil

	// Add parameters to local symbols
	param := f.Params
	for param != nil {
		p.addSymbol(param.Name, 5, param.DataType, param.Offset, param.Size) // SYM_LOCAL
		param = param.Next
	}

	node.Child = p.parseBlock()
	p.currentFun = nil

	return node
}

// parseType parses a type name
func (p *Parser) parseType() int {
	tokenType := p.lexer.token.Type
	p.lexer.NextToken()

	switch tokenType {
	case KW_VOID:
		return RT_VOID
	case KW_U8:
		return RT_U8
	case KW_U16:
		return RT_U16
	case KW_U32:
		return RT_U32
	case KW_U64:
		return RT_U64
	case KW_I8:
		return RT_I8
	case KW_I16:
		return RT_I16
	case KW_I32:
		return RT_I32
	case KW_I64:
		return RT_I64
	case KW_F64:
		return RT_F64
	case KW_BOOL:
		return RT_BOOL
	case KW_U0:
		return RT_VOID
	default:
		return RT_I64
	}
}

// parseClassDeclaration parses a class declaration
func (p *Parser) parseClassDeclaration() *ASTNode {
	node := &ASTNode{Type: AST_CLASS_DECL, Line: p.lexer.token.Line}

	p.lexer.NextToken() // Skip 'class'

	if p.lexer.token.Type != TK_IDENT {
		p.Error(p.lexer.token.Line, "Expected class name")
		return node
	}

	name := p.lexer.token.Ident
	node.Ident = name
	p.lexer.NextToken()

	// Create class
	c := &Class{Name: name}
	p.addClass(c)
	node.ClassDef = c

	if err := p.lexer.Match(TK_LBRACE); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	offset := 0
	for p.lexer.token.Type != TK_RBRACE && p.lexer.token.Type != TK_EOF {
		memberType := p.parseType()
		memberSize := typeSize(memberType)

		for {
			if p.lexer.token.Type != TK_IDENT {
				break
			}

			memberName := p.lexer.token.Ident
			p.lexer.NextToken()

			// Array size
			arraySize := 1
			if p.lexer.token.Type == TK_LBRACKET {
				p.lexer.NextToken()
				if p.lexer.token.Type == TK_I64 {
					arraySize = int(p.lexer.token.I64Val)
				}
				p.lexer.NextToken()
				if err := p.lexer.Match(TK_RBRACKET); err != nil {
					p.Error(p.lexer.token.Line, "%v", err)
				}
			}

			// Add member
			member := &Symbol{
				Name:     memberName,
				Type:     1, // SYM_VAR
				DataType: memberType,
				Offset:   offset,
				Size:     memberSize * arraySize,
			}
			member.Next = c.Members
			c.Members = member

			offset += member.Size
			c.MemberCnt++

			if p.lexer.token.Type != TK_COMMA {
				break
			}
			p.lexer.NextToken()
		}

		if err := p.lexer.Match(TK_SEMICOLON); err != nil {
			p.Error(p.lexer.token.Line, "%v", err)
		}
	}

	if err := p.lexer.Match(TK_RBRACE); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
	}
	if err := p.lexer.Match(TK_SEMICOLON); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
	}

	c.Size = offset
	return node
}

// parseExternDeclaration parses an extern declaration
func (p *Parser) parseExternDeclaration() *ASTNode {
	node := &ASTNode{Type: AST_EXTERN_DECL, Line: p.lexer.token.Line}

	p.lexer.NextToken() // Skip 'extern'

	returnType := p.parseType()

	if p.lexer.token.Type != TK_IDENT {
		p.Error(p.lexer.token.Line, "Expected function name")
		return node
	}

	node.Ident = p.lexer.token.Ident
	node.DataType = returnType
	p.lexer.NextToken()

	// Skip parameters
	if err := p.lexer.Match(TK_LPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}
	for p.lexer.token.Type != TK_RPAREN && p.lexer.token.Type != TK_EOF {
		p.lexer.NextToken()
	}
	if err := p.lexer.Match(TK_RPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
	}
	if err := p.lexer.Match(TK_SEMICOLON); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
	}

	return node
}

// parseBlock parses a block of statements
func (p *Parser) parseBlock() *ASTNode {
	node := &ASTNode{Type: AST_BLOCK, Line: p.lexer.token.Line}

	if err := p.lexer.Match(TK_LBRACE); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	var lastStmt *ASTNode
	for p.lexer.token.Type != TK_RBRACE && p.lexer.token.Type != TK_EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			if lastStmt != nil {
				lastStmt.Next = stmt
			} else {
				node.Child = stmt
			}
			lastStmt = stmt
		}
	}

	if err := p.lexer.Match(TK_RBRACE); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
	}

	return node
}

// parseStatement parses a single statement
func (p *Parser) parseStatement() *ASTNode {
	switch p.lexer.token.Type {
	case TK_LBRACE:
		return p.parseBlock()
	case KW_IF:
		return p.parseIfStatement()
	case KW_WHILE:
		return p.parseWhileStatement()
	case KW_FOR:
		return p.parseForStatement()
	case KW_DO:
		return p.parseDoStatement()
	case KW_RETURN:
		p.lexer.NextToken()
		node := &ASTNode{Type: AST_RETURN, Line: p.lexer.token.Line}
		if p.lexer.token.Type != TK_SEMICOLON {
			node.Child = p.parseExpression()
		}
		p.lexer.NextToken() // Skip semicolon
		return node
	case KW_BREAK:
		p.lexer.NextToken()
		p.lexer.NextToken() // Skip semicolon
		return &ASTNode{Type: AST_BREAK, Line: p.lexer.token.Line}
	case KW_CONTINUE:
		p.lexer.NextToken()
		p.lexer.NextToken() // Skip semicolon
		return &ASTNode{Type: AST_CONTINUE, Line: p.lexer.token.Line}
	case KW_SWITCH:
		return p.parseSwitchStatement()
	case KW_GOTO:
		return p.parseGotoStatement()
	case KW_VOID, KW_U0, KW_U8, KW_U16, KW_U32, KW_U64, KW_I8, KW_I16, KW_I32, KW_I64, KW_F64, KW_BOOL:
		return p.parseVarDeclaration()
	default:
		node := p.parseExpression()
		p.lexer.NextToken() // Skip semicolon
		return node
	}
}

// parseIfStatement parses an if statement
func (p *Parser) parseIfStatement() *ASTNode {
	node := &ASTNode{Type: AST_IF, Line: p.lexer.token.Line}
	p.lexer.NextToken()

	if err := p.lexer.Match(TK_LPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	node.Child = p.parseExpression()

	if err := p.lexer.Match(TK_RPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	node.Sibling = p.parseStatement()

	if p.lexer.token.Type == KW_ELSE {
		p.lexer.NextToken()
		node.Next = p.parseStatement()
	}

	return node
}

// parseWhileStatement parses a while statement
func (p *Parser) parseWhileStatement() *ASTNode {
	node := &ASTNode{Type: AST_WHILE, Line: p.lexer.token.Line}
	p.lexer.NextToken()

	if err := p.lexer.Match(TK_LPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	node.Child = p.parseExpression()

	if err := p.lexer.Match(TK_RPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	node.Sibling = p.parseStatement()

	return node
}

// parseForStatement parses a for statement
func (p *Parser) parseForStatement() *ASTNode {
	node := &ASTNode{Type: AST_FOR, Line: p.lexer.token.Line}
	p.lexer.NextToken()

	if err := p.lexer.Match(TK_LPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	// Init
	var init *ASTNode
	if p.lexer.token.Type != TK_SEMICOLON {
		init = p.parseExpression()
	}
	p.lexer.NextToken() // Skip semicolon

	// Condition
	var cond *ASTNode
	if p.lexer.token.Type != TK_SEMICOLON {
		cond = p.parseExpression()
	}
	p.lexer.NextToken() // Skip semicolon

	// Increment
	var incr *ASTNode
	if p.lexer.token.Type != TK_RPAREN {
		incr = p.parseExpression()
	}

	if err := p.lexer.Match(TK_RPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	node.Child = init
	node.Sibling = cond
	node.Next = incr
	node.Body = p.parseStatement() // Body stored in Body field

	return node
}

// parseDoStatement parses a do-while statement
func (p *Parser) parseDoStatement() *ASTNode {
	node := &ASTNode{Type: AST_DO, Line: p.lexer.token.Line}
	p.lexer.NextToken()

	node.Child = p.parseStatement()

	if err := p.lexer.Match(KW_WHILE); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}
	if err := p.lexer.Match(TK_LPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	node.Sibling = p.parseExpression()

	if err := p.lexer.Match(TK_RPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}
	if err := p.lexer.Match(TK_SEMICOLON); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	return node
}

// parseSwitchStatement parses a switch statement
func (p *Parser) parseSwitchStatement() *ASTNode {
	node := &ASTNode{Type: AST_SWITCH, Line: p.lexer.token.Line}
	p.lexer.NextToken()

	if err := p.lexer.Match(TK_LPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	node.Child = p.parseExpression()

	if err := p.lexer.Match(TK_RPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	if err := p.lexer.Match(TK_LBRACE); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
		return node
	}

	var lastCase *ASTNode
	for p.lexer.token.Type == KW_CASE || p.lexer.token.Type == KW_DEFAULT {
		caseNode := &ASTNode{Type: AST_CASE}

		if p.lexer.token.Type == KW_CASE {
			p.lexer.NextToken()
			caseNode.I64Val = p.parseExpression().I64Val
		} else {
			caseNode.I64Val = -1 // Default
			p.lexer.NextToken()
		}

		if err := p.lexer.Match(TK_COLON); err != nil {
			p.Error(p.lexer.token.Line, "%v", err)
		}

		caseNode.Child = p.parseBlock()

		if lastCase != nil {
			lastCase.Next = caseNode
		} else {
			node.Sibling = caseNode
		}
		lastCase = caseNode
	}

	if err := p.lexer.Match(TK_RBRACE); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
	}

	return node
}

// parseGotoStatement parses a goto statement
func (p *Parser) parseGotoStatement() *ASTNode {
	node := &ASTNode{Type: AST_GOTO, Line: p.lexer.token.Line}
	p.lexer.NextToken()

	if p.lexer.token.Type == TK_IDENT {
		node.Ident = p.lexer.token.Ident
		p.lexer.NextToken()
	}

	p.lexer.NextToken() // Skip semicolon
	return node
}

// parseVarDeclaration parses a variable declaration
func (p *Parser) parseVarDeclaration() *ASTNode {
	node := &ASTNode{Type: AST_VAR_DECL, Line: p.lexer.token.Line}

	dataType := p.parseType()

	if p.lexer.token.Type != TK_IDENT {
		p.Error(p.lexer.token.Line, "Expected variable name")
		return node
	}

	name := p.lexer.token.Ident
	p.lexer.NextToken()

	// Array size
	arraySize := 1
	if p.lexer.token.Type == TK_LBRACKET {
		p.lexer.NextToken()
		if p.lexer.token.Type == TK_I64 {
			arraySize = int(p.lexer.token.I64Val)
		}
		p.lexer.NextToken()
		if err := p.lexer.Match(TK_RBRACKET); err != nil {
			p.Error(p.lexer.token.Line, "%v", err)
		}
	}

	node.DataType = dataType
	node.Ident = name
	node.I64Val = int64(arraySize)

	// Initialization
	if p.lexer.token.Type == TK_EQUAL {
		p.lexer.NextToken()
		node.Child = p.parseExpression()
	}

	if p.lexer.token.Type == TK_SEMICOLON {
		p.lexer.NextToken() // Skip semicolon
	}

	// Add to symbol table
	symType := SYM_LOCAL
	if p.currentFun == nil {
		symType = SYM_VAR
	}
	node.Sym = p.addSymbol(name, symType, dataType, 0, typeSize(dataType)*arraySize)

	return node
}

// parseExpression parses an expression
func (p *Parser) parseExpression() *ASTNode {
	return p.parseBinaryOp(p.parsePrimary(), 1)
}

// parsePrimary parses primary expressions
func (p *Parser) parsePrimary() *ASTNode {
	var node *ASTNode

	switch p.lexer.token.Type {
	case TK_I64:
		node = &ASTNode{
			Type:     AST_LITERAL,
			DataType: RT_I64,
			I64Val:   p.lexer.token.I64Val,
			Line:     p.lexer.token.Line,
		}
		p.lexer.NextToken()

	case TK_F64:
		node = &ASTNode{
			Type:     AST_LITERAL,
			DataType: RT_F64,
			F64Val:   p.lexer.token.F64Val,
			Line:     p.lexer.token.Line,
		}
		p.lexer.NextToken()

	case TK_STR:
		node = &ASTNode{
			Type:     AST_STRING,
			DataType: RT_PTR,
			StrVal:   p.lexer.token.StrVal,
			Line:     p.lexer.token.Line,
		}
		p.lexer.NextToken()

	case TK_CHAR_CONST:
		node = &ASTNode{
			Type:     AST_LITERAL,
			DataType: RT_U8,
			I64Val:   p.lexer.token.I64Val,
			Line:     p.lexer.token.Line,
		}
		p.lexer.NextToken()

	case TK_IDENT:
		node = &ASTNode{
			Type:  AST_IDENT,
			Ident: p.lexer.token.Ident,
			Line:  p.lexer.token.Line,
		}
		node.Sym = p.lookupSymbol(node.Ident)
		p.lexer.NextToken()

		// Member access
		for p.lexer.token.Type == TK_DOT || p.lexer.token.Type == TK_ARROW {
			isArrow := p.lexer.token.Type == TK_ARROW
			_ = isArrow
			p.lexer.NextToken()

			member := &ASTNode{Type: AST_MEMBER_ACCESS, Child: node}

			if p.lexer.token.Type == TK_IDENT {
				member.Ident = p.lexer.token.Ident
				p.lexer.NextToken()
			}

			node = member
		}

		// Array access
		for p.lexer.token.Type == TK_LBRACKET {
			p.lexer.NextToken()
			index := p.parseExpression()
			p.lexer.NextToken() // Skip ]

			array := &ASTNode{
				Type:    AST_ARRAY_ACCESS,
				Child:   node,
				Sibling: index,
			}
			node = array
		}

	case TK_LPAREN:
		p.lexer.NextToken()
		node = p.parseExpression()
		p.lexer.NextToken() // Skip )

	case TK_MINUS, TK_NOT, TK_TILDE:
		op := p.lexer.token.Type
		p.lexer.NextToken()
		node = &ASTNode{
			Type:    AST_UNARY_OP,
			I64Val:  int64(op),
			Child:   p.parsePrimary(),
			Line:    p.lexer.token.Line,
		}

	case KW_SIZEOF:
		p.lexer.NextToken()
		p.lexer.NextToken() // Skip (
		node = &ASTNode{
			Type:     AST_LITERAL,
			DataType: RT_I64,
			I64Val:   8, // Default size
			Line:     p.lexer.token.Line,
		}
		if p.lexer.token.Type == TK_IDENT {
			// Could look up class size here
			p.lexer.NextToken()
		}
		p.lexer.NextToken() // Skip )

	default:
		p.Error(p.lexer.token.Line, "Unexpected token in expression: %d", p.lexer.token.Type)
		node = &ASTNode{Type: AST_LITERAL, I64Val: 0}
	}

	// Function call
	if node.Type == AST_IDENT && p.lexer.token.Type == TK_LPAREN {
		return p.parseFunctionCall(node.Ident)
	}

	return node
}

// parseFunctionCall parses a function call
func (p *Parser) parseFunctionCall(name string) *ASTNode {
	node := &ASTNode{
		Type:  AST_FUN_CALL,
		Ident: name,
		Line:  p.lexer.token.Line,
	}

	// Lookup function
	for f := p.functions; f != nil; f = f.Next {
		if f.Name == name {
			node.Fun = f
			break
		}
	}

	p.lexer.NextToken() // Skip (

	if p.lexer.token.Type != TK_RPAREN {
		arg := p.parseExpression()
		node.Child = arg
		lastArg := arg

		for p.lexer.token.Type == TK_COMMA {
			p.lexer.NextToken()
			arg = p.parseExpression()
			if lastArg != nil {
				lastArg.Next = arg
			}
			lastArg = arg
		}
	}

	if err := p.lexer.Match(TK_RPAREN); err != nil {
		p.Error(p.lexer.token.Line, "%v", err)
	}

	return node
}

// getPrecedence returns the precedence of an operator
func getPrecedence(tokenType int) int {
	switch tokenType {
	case TK_PLUS, TK_MINUS:
		return 10
	case TK_STAR, TK_SLASH, TK_PERCENT:
		return 12
	case TK_SHIFT_LEFT, TK_SHIFT_RIGHT:
		return 11
	case TK_LESS, TK_GREATER, TK_LESS_EQU, TK_GREATER_EQU:
		return 8
	case TK_EQUAL2, TK_NOT_EQUAL:
		return 7
	case TK_AND:
		return 6
	case TK_XOR:
		return 5
	case TK_OR:
		return 4
	case TK_AND_AND:
		return 3
	case TK_OR_OR:
		return 2
	case TK_EQUAL, TK_PLUS_EQU, TK_MINUS_EQU, TK_STAR_EQU, TK_SLASH_EQU:
		return 1
	default:
		return 0
	}
}

// parseBinaryOp parses binary operations with precedence
func (p *Parser) parseBinaryOp(lhs *ASTNode, minPrec int) *ASTNode {
	for {
		prec := getPrecedence(p.lexer.token.Type)
		if prec < minPrec {
			break
		}

		op := p.lexer.token.Type
		p.lexer.NextToken()

		rhs := p.parsePrimary()

		for {
			nextPrec := getPrecedence(p.lexer.token.Type)
			if nextPrec <= prec {
				break
			}
			rhs = p.parseBinaryOp(rhs, nextPrec)
		}

		node := &ASTNode{
			Type:    AST_BINARY_OP,
			I64Val:  int64(op),
			Child:   lhs,
			Sibling: rhs,
			Line:    lhs.Line,
		}
		lhs = node
	}

	return lhs
}

// Symbol table functions
func (p *Parser) addSymbol(name string, symType int, dataType int, offset int, size int) *Symbol {
	sym := &Symbol{
		Name:     name,
		Type:     symType,
		DataType: dataType,
		Offset:   offset,
		Size:     size,
	}

	if p.currentFun != nil {
		if symType == SYM_LOCAL {
			p.currentFun.LocalCnt++
			sym.Offset = -8 * p.currentFun.LocalCnt
		}
		sym.Next = p.localSyms
		p.localSyms = sym
	} else {
		sym.Next = p.globalSyms
		p.globalSyms = sym
	}

	return sym
}

func (p *Parser) lookupSymbol(name string) *Symbol {
	// Check locals first
	for sym := p.localSyms; sym != nil; sym = sym.Next {
		if sym.Name == name {
			return sym
		}
	}
	// Check globals
	for sym := p.globalSyms; sym != nil; sym = sym.Next {
		if sym.Name == name {
			return sym
		}
	}
	return nil
}

func (p *Parser) findSymbol(name string, table *Symbol) *Symbol {
	for sym := table; sym != nil; sym = sym.Next {
		if sym.Name == name {
			return sym
		}
	}
	return nil
}

func (p *Parser) addClass(c *Class) {
	// Check if already exists
	for existing := p.classes; existing != nil; existing = existing.Next {
		if existing.Name == c.Name {
			return
		}
	}
	c.Next = p.classes
	p.classes = c
}

func (p *Parser) addFunction(f *Function) {
	// Check if already exists
	for existing := p.functions; existing != nil; existing = existing.Next {
		if existing.Name == f.Name {
			return
		}
	}
	f.Next = p.functions
	p.functions = f
}

// typeSize returns the size of a type
func typeSize(t int) int {
	switch t {
	case RT_VOID:
		return 0
	case RT_U8, RT_I8, RT_BOOL:
		return 1
	case RT_U16, RT_I16:
		return 2
	case RT_U32, RT_I32:
		return 4
	case RT_U64, RT_I64, RT_F64, RT_PTR:
		return 8
	default:
		return 8
	}
}
