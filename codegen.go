package main

import (
	"encoding/binary"
)

// x64 Opcodes
const (
	OC_PUSH_RBP       = 0x55
	OC_MOV_RBP_RSP    = 0x4889E5
	OC_LEAVE          = 0xC9
	OC_RET            = 0xC3
	OC_PUSH_RAX       = 0x50
	OC_POP_RAX        = 0x58
	OC_PUSH_RCX       = 0x51
	OC_POP_RCX        = 0x59
	OC_PUSH_RDX       = 0x52
	OC_POP_RDX        = 0x5A
	OC_PUSH_RSI       = 0x56
	OC_POP_RSI        = 0x5E
	OC_PUSH_RDI       = 0x57
	OC_POP_RDI        = 0x5F
	OC_PUSH_R8        = 0x4150
	OC_PUSH_R9        = 0x4151
	OC_MOV_RAX_IMM    = 0x48B8
	OC_MOV_RCX_IMM    = 0x48B9
	OC_MOV_RDX_IMM    = 0x48BA
	OC_MOV_RSI_IMM    = 0x48BE
	OC_MOV_RDI_IMM    = 0x48BF
	OC_MOV_RAX_RCX    = 0x4889C8
	OC_MOV_RCX_RAX    = 0x4889C1
	OC_ADD_RAX_RCX    = 0x4801C8
	OC_SUB_RAX_RCX    = 0x4829C8
	OC_IMUL_RAX_RCX   = 0x480FAFC8
	OC_XOR_RAX_RCX    = 0x4831C8
	OC_AND_RAX_RCX    = 0x4821C8
	OC_OR_RAX_RCX     = 0x4809C8
	OC_CMP_RAX_RCX    = 0x4839C8
	OC_TEST_RAX_RCX   = 0x4885C0
	OC_NEG_RAX        = 0x48F7D8
	OC_NOT_RAX        = 0x48F7D0
	OC_INC_RAX        = 0x48FFC0
	OC_DEC_RAX        = 0x48FFC8
	OC_NOP            = 0x90
	OC_INT3           = 0xCC
	OC_CALL_REL32     = 0xE8
	OC_JMP_REL32      = 0xE9
	OC_JZ_REL32       = 0x0F84
	OC_JNZ_REL32      = 0x0F85
	OC_JL_REL32       = 0x0F8C
	OC_JLE_REL32      = 0x0F8E
	OC_JG_REL32       = 0x0F8F
	OC_JGE_REL32      = 0x0F8D
	OC_JZ_REL8        = 0x74
	OC_JNZ_REL8       = 0x75
	OC_SUB_RSP_IMM32  = 0x4881EC
	OC_ADD_RSP_IMM32  = 0x4881C4
)

// CodeGenerator generates x64 machine code from AST
type CodeGenerator struct {
	code         []byte
	codeOffset   int
	stringPool   []byte
	stringOffset int
	functions    map[string]int
	labels       map[string]int
	patches      []Patch
	labelCount   int
	breakStack   []int
	continueStack []int
}

// Patch represents a location to patch with an address
type Patch struct {
	Offset    int
	Label     string
	Type      int // 1 = relative32
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		code:          make([]byte, 0x100000), // 1MB initial capacity
		functions:     make(map[string]int),
		labels:        make(map[string]int),
		breakStack:    make([]int, 0),
		continueStack: make([]int, 0),
	}
}

// EmitByte emits a single byte
func (cg *CodeGenerator) EmitByte(b byte) {
	if cg.codeOffset >= len(cg.code) {
		newCode := make([]byte, len(cg.code)*2)
		copy(newCode, cg.code)
		cg.code = newCode
	}
	cg.code[cg.codeOffset] = b
	cg.codeOffset++
}

// EmitBytes emits multiple bytes
func (cg *CodeGenerator) EmitBytes(bytes []byte) {
	for _, b := range bytes {
		cg.EmitByte(b)
	}
}

// EmitU16 emits a 16-bit value
func (cg *CodeGenerator) EmitU16(v uint16) {
	cg.EmitByte(byte(v))
	cg.EmitByte(byte(v >> 8))
}

// EmitU32 emits a 32-bit value
func (cg *CodeGenerator) EmitU32(v uint32) {
	cg.EmitByte(byte(v))
	cg.EmitByte(byte(v >> 8))
	cg.EmitByte(byte(v >> 16))
	cg.EmitByte(byte(v >> 24))
}

// EmitU64 emits a 64-bit value
func (cg *CodeGenerator) EmitU64(v uint64) {
	for i := 0; i < 8; i++ {
		cg.EmitByte(byte(v >> (i * 8)))
	}
}

// EmitPushRAX emits PUSH RAX
func (cg *CodeGenerator) EmitPushRAX() {
	cg.EmitByte(OC_PUSH_RAX)
}

// EmitPopRAX emits POP RAX
func (cg *CodeGenerator) EmitPopRAX() {
	cg.EmitByte(OC_POP_RAX)
}

// EmitMovRaxImm emits MOV RAX, imm64
func (cg *CodeGenerator) EmitMovRaxImm(imm int64) {
	cg.EmitByte(0x48)
	cg.EmitByte(0xB8)
	cg.EmitU64(uint64(imm))
}

// EmitCall emits CALL rel32
func (cg *CodeGenerator) EmitCall(addr int) {
	rel := int32(addr - cg.codeOffset - 5)
	cg.EmitByte(OC_CALL_REL32)
	cg.EmitU32(uint32(rel))
}

// EmitJmp emits JMP rel32
func (cg *CodeGenerator) EmitJmp(addr int) {
	rel := int32(addr - cg.codeOffset - 5)
	cg.EmitByte(OC_JMP_REL32)
	cg.EmitU32(uint32(rel))
}

// EmitJz emits JZ rel32
func (cg *CodeGenerator) EmitJz(addr int) {
	cg.EmitByte(0x0F)
	cg.EmitByte(0x84)
	rel := int32(addr - cg.codeOffset - 4)
	cg.EmitU32(uint32(rel))
}

// EmitJnz emits JNZ rel32
func (cg *CodeGenerator) EmitJnz(addr int) {
	cg.EmitByte(0x0F)
	cg.EmitByte(0x85)
	rel := int32(addr - cg.codeOffset - 4)
	cg.EmitU32(uint32(rel))
}

// EmitRet emits RET
func (cg *CodeGenerator) EmitRet() {
	cg.EmitByte(OC_RET)
}

// EmitLeave emits LEAVE
func (cg *CodeGenerator) EmitLeave() {
	cg.EmitByte(OC_LEAVE)
}

// EmitNop emits NOP
func (cg *CodeGenerator) EmitNop() {
	cg.EmitByte(OC_NOP)
}

// CurrentOffset returns the current code offset
func (cg *CodeGenerator) CurrentOffset() int {
	return cg.codeOffset
}

// AddString adds a string to the string pool
func (cg *CodeGenerator) AddString(s string) int {
	offset := cg.stringOffset
	for _, c := range []byte(s) {
		if cg.stringOffset >= len(cg.stringPool) {
			cg.stringPool = append(cg.stringPool, 0)
		}
		cg.stringPool[cg.stringOffset] = c
		cg.stringOffset++
	}
	// Null terminator
	if cg.stringOffset >= len(cg.stringPool) {
		cg.stringPool = append(cg.stringPool, 0)
	}
	cg.stringPool[cg.stringOffset] = 0
	cg.stringOffset++
	return offset
}

// Generate generates code from AST
func (cg *CodeGenerator) Generate(ast *ASTNode) error {
	if ast == nil || ast.Child == nil {
		return nil
	}

	// Generate code for each function
	node := ast.Child
	for node != nil {
		if node.Type == AST_FUN_DECL {
			if node.Fun != nil && (node.Fun.Flags&0x01) == 0 { // Not extern
				cg.GenerateFunction(node)
			}
		}
		node = node.Next
	}

	// Patch all labels
	cg.ApplyPatches()

	return nil
}

// GenerateFunction generates code for a function
func (cg *CodeGenerator) GenerateFunction(node *ASTNode) {
	f := node.Fun
	f.CodeStart = cg.codeOffset
	cg.functions[f.Name] = cg.codeOffset

	// Function prologue
	cg.EmitByte(OC_PUSH_RBP)           // PUSH RBP
	cg.EmitByte(0x48)                  // REX.W
	cg.EmitByte(0x89)                  // MOV
	cg.EmitByte(0xE5)                  // RBP, RSP

	// Allocate stack space for locals
	stackSize := 32 // Minimum stack frame
	if f.LocalCnt > 0 {
		stackSize += f.LocalCnt * 8
	}
	f.Size = stackSize

	if stackSize > 0 {
		cg.EmitByte(0x48)
		cg.EmitByte(0x81)
		cg.EmitByte(0xEC)
		cg.EmitU32(uint32(stackSize)) // SUB RSP, stackSize
	}

	// Generate function body
	cg.GenerateStatement(node.Child)

	// Function epilogue (if not already returned)
	cg.EmitLeave()
	cg.EmitRet()

	f.CodeEnd = cg.codeOffset
}

// GenerateStatement generates code for a statement
func (cg *CodeGenerator) GenerateStatement(node *ASTNode) {
	if node == nil {
		return
	}

	switch node.Type {
	case AST_BLOCK:
		stmt := node.Child
		for stmt != nil {
			cg.GenerateStatement(stmt)
			stmt = stmt.Next
		}

	case AST_IF:
		cg.GenerateExpression(node.Child)
		jzOffset := cg.codeOffset
		cg.EmitJz(0) // Placeholder

		cg.GenerateStatement(node.Sibling)

		if node.Next != nil {
			// Has else
			jmpOffset := cg.codeOffset
			cg.EmitJmp(0) // Placeholder

			// Patch JZ
			endOffset := cg.codeOffset
			cg.PatchRelative32(jzOffset, endOffset)

			cg.GenerateStatement(node.Next)

			// Patch JMP
			endOffset = cg.codeOffset
			cg.PatchRelative32(jmpOffset, endOffset)
		} else {
			// No else
			endOffset := cg.codeOffset
			cg.PatchRelative32(jzOffset, endOffset)
		}

	case AST_WHILE:
		condOffset := cg.codeOffset
		breakOffset := 0
		cg.breakStack = append(cg.breakStack, breakOffset)
		cg.continueStack = append(cg.continueStack, condOffset)

		cg.GenerateExpression(node.Child)
		jzOffset := cg.codeOffset
		cg.EmitJz(0) // Placeholder

		cg.GenerateStatement(node.Sibling)

		// JMP to condition
		cg.EmitJmp(condOffset)

		// Patch JZ
		endOffset := cg.codeOffset
		cg.PatchRelative32(jzOffset, endOffset)

		// Pop stacks
		cg.breakStack = cg.breakStack[:len(cg.breakStack)-1]
		cg.continueStack = cg.continueStack[:len(cg.continueStack)-1]

	case AST_FOR:
		// Init
		if node.Child != nil {
			cg.GenerateExpression(node.Child)
		}

		condOffset := cg.codeOffset
		breakOffset := 0
		cg.breakStack = append(cg.breakStack, breakOffset)
		cg.continueStack = append(cg.continueStack, 0) // Will be set after body

		// Condition
		var jzOffset int
		if node.Sibling != nil {
			cg.GenerateExpression(node.Sibling)
			jzOffset = cg.codeOffset
			cg.EmitJz(0) // Placeholder
		}

		cg.GenerateStatement(node.Body) // Body

		// Increment
		incrOffset := cg.codeOffset
		cg.continueStack[len(cg.continueStack)-1] = incrOffset
		if node.Next != nil {
			cg.GenerateExpression(node.Next)
		}

		// JMP to condition
		cg.EmitJmp(condOffset)

		// Patch JZ
		if node.Sibling != nil {
			endOffset := cg.codeOffset
			cg.PatchRelative32(jzOffset, endOffset)
		}

		// Pop stacks
		cg.breakStack = cg.breakStack[:len(cg.breakStack)-1]
		cg.continueStack = cg.continueStack[:len(cg.continueStack)-1]

	case AST_RETURN:
		if node.Child != nil {
			cg.GenerateExpression(node.Child)
		}
		cg.EmitLeave()
		cg.EmitRet()

	case AST_BREAK:
		// Store break location for patching
		breakOffset := cg.codeOffset
		cg.EmitJmp(0) // Placeholder
		// Add to break stack
		if len(cg.breakStack) > 0 {
			cg.breakStack[len(cg.breakStack)-1] = breakOffset
		}

	case AST_CONTINUE:
		// JMP to continue label
		if len(cg.continueStack) > 0 {
			continueOffset := cg.continueStack[len(cg.continueStack)-1]
			cg.EmitJmp(continueOffset)
		}

	case AST_BINARY_OP:
		cg.GenerateExpression(node)

	case AST_ASSIGN:
		cg.GenerateExpression(node)
	}
}

// GenerateExpression generates code for an expression
func (cg *CodeGenerator) GenerateExpression(node *ASTNode) {
	if node == nil {
		return
	}

	switch node.Type {
	case AST_LITERAL:
		cg.EmitMovRaxImm(node.I64Val)

	case AST_IDENT:
		// Load variable value (simplified - would need proper symbol lookup)
		cg.EmitMovRaxImm(0)

	case AST_FUN_CALL:
		// Push arguments (simplified - reverse order for x64 calling convention)
		arg := node.Child
		argCount := 0
		for arg != nil {
			cg.GenerateExpression(arg)
			cg.EmitPushRAX()
			arg = arg.Next
			argCount++
		}

		// Call function
		if addr, ok := cg.functions[node.Ident]; ok {
			cg.EmitCall(addr)
		} else {
			// Forward declaration or extern
			cg.EmitCall(0) // Will be patched
			cg.patches = append(cg.patches, Patch{
				Offset: cg.codeOffset - 4,
				Label:  node.Ident,
				Type:   1,
			})
		}

		// Clean up stack
		if argCount > 0 {
			cg.EmitByte(0x48)
			cg.EmitByte(0x81)
			cg.EmitByte(0xC4)
			cg.EmitU32(uint32(argCount * 8)) // ADD RSP, argCount*8
		}

	case AST_BINARY_OP:
		cg.GenerateExpression(node.Child)
		cg.EmitPushRAX()
		cg.GenerateExpression(node.Sibling)
		cg.EmitPopRAX()

		op := int(node.I64Val)
		switch op {
		case TK_PLUS:
			cg.EmitByte(0x48)
			cg.EmitByte(0x01)
			cg.EmitByte(0xC1) // ADD RCX, RAX
			cg.EmitByte(0x48)
			cg.EmitByte(0x89)
			cg.EmitByte(0xC8) // MOV RAX, RCX

		case TK_MINUS:
			cg.EmitByte(0x48)
			cg.EmitByte(0x29)
			cg.EmitByte(0xC1) // SUB RCX, RAX
			cg.EmitByte(0x48)
			cg.EmitByte(0x89)
			cg.EmitByte(0xC8) // MOV RAX, RCX

		case TK_STAR:
			cg.EmitByte(0x48)
			cg.EmitByte(0x0F)
			cg.EmitByte(0xAF)
			cg.EmitByte(0xC1) // IMUL RCX, RAX
			cg.EmitByte(0x48)
			cg.EmitByte(0x89)
			cg.EmitByte(0xC8) // MOV RAX, RCX

		case TK_SLASH:
			// Division (simplified)
			cg.EmitByte(0x48)
			cg.EmitByte(0x99) // CQO
			cg.EmitByte(0x48)
			cg.EmitByte(0xF7)
			cg.EmitByte(0xF9) // IDIV RCX

		case TK_PERCENT:
			// Modulo (simplified)
			cg.EmitByte(0x48)
			cg.EmitByte(0x99) // CQO
			cg.EmitByte(0x48)
			cg.EmitByte(0xF7)
			cg.EmitByte(0xF9) // IDIV RCX
			cg.EmitByte(0x48)
			cg.EmitByte(0x89)
			cg.EmitByte(0xD0) // MOV RAX, RDX

		case TK_EQUAL2:
			cg.EmitByte(0x48)
			cg.EmitByte(0x39)
			cg.EmitByte(0xC8) // CMP RAX, RCX
			cg.EmitByte(0x0F)
			cg.EmitByte(0x94)
			cg.EmitByte(0xC0) // SETNE AL
			cg.EmitByte(0x48)
			cg.EmitByte(0x0F)
			cg.EmitByte(0xB6)
			cg.EmitByte(0xC0) // MOVZX RAX, AL

		case TK_NOT_EQUAL:
			cg.EmitByte(0x48)
			cg.EmitByte(0x39)
			cg.EmitByte(0xC8) // CMP RAX, RCX
			cg.EmitByte(0x0F)
			cg.EmitByte(0x95)
			cg.EmitByte(0xC0) // SETNE AL
			cg.EmitByte(0x48)
			cg.EmitByte(0x0F)
			cg.EmitByte(0xB6)
			cg.EmitByte(0xC0) // MOVZX RAX, AL

		case TK_LESS:
			cg.EmitByte(0x48)
			cg.EmitByte(0x39)
			cg.EmitByte(0xC8) // CMP RAX, RCX
			cg.EmitByte(0x0F)
			cg.EmitByte(0x9C)
			cg.EmitByte(0xC0) // SETL AL
			cg.EmitByte(0x48)
			cg.EmitByte(0x0F)
			cg.EmitByte(0xB6)
			cg.EmitByte(0xC0) // MOVZX RAX, AL

		case TK_GREATER:
			cg.EmitByte(0x48)
			cg.EmitByte(0x39)
			cg.EmitByte(0xC8) // CMP RAX, RCX
			cg.EmitByte(0x0F)
			cg.EmitByte(0x9F)
			cg.EmitByte(0xC0) // SETG AL
			cg.EmitByte(0x48)
			cg.EmitByte(0x0F)
			cg.EmitByte(0xB6)
			cg.EmitByte(0xC0) // MOVZX RAX, AL

		case TK_AND:
			cg.EmitByte(0x48)
			cg.EmitByte(0x21)
			cg.EmitByte(0xC8) // AND RAX, RCX

		case TK_OR:
			cg.EmitByte(0x48)
			cg.EmitByte(0x09)
			cg.EmitByte(0xC8) // OR RAX, RCX

		case TK_XOR:
			cg.EmitByte(0x48)
			cg.EmitByte(0x31)
			cg.EmitByte(0xC8) // XOR RAX, RCX
		}

	case AST_UNARY_OP:
		cg.GenerateExpression(node.Child)
		op := int(node.I64Val)
		if op == TK_MINUS {
			cg.EmitByte(0x48)
			cg.EmitByte(0xF7)
			cg.EmitByte(0xD8) // NEG RAX
		} else if op == TK_NOT {
			cg.EmitByte(0x48)
			cg.EmitByte(0xF7)
			cg.EmitByte(0xD0) // NOT RAX
		}
	}
}

// PatchRelative32 patches a relative32 offset
func (cg *CodeGenerator) PatchRelative32(offset int, target int) {
	rel := int32(target - offset - 4)
	binary.LittleEndian.PutUint32(cg.code[offset:], uint32(rel))
}

// ApplyPatches applies all pending patches
func (cg *CodeGenerator) ApplyPatches() {
	for _, patch := range cg.patches {
		if addr, ok := cg.functions[patch.Label]; ok {
			rel := int32(addr - patch.Offset - 4)
			binary.LittleEndian.PutUint32(cg.code[patch.Offset:], uint32(rel))
		}
	}
}

// GetCode returns the generated code
func (cg *CodeGenerator) GetCode() []byte {
	return cg.code[:cg.codeOffset]
}

// GetStringPool returns the string pool
func (cg *CodeGenerator) GetStringPool() []byte {
	return cg.stringPool[:cg.stringOffset]
}
