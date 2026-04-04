package main

import (
	"fmt"
	"unicode"
)

// Token types
const (
	TK_EOF = iota
	TK_IDENT
	TK_I64
	TK_F64
	TK_STR
	TK_CHAR_CONST

	// Operators
	TK_PLUS
	TK_MINUS
	TK_STAR
	TK_SLASH
	TK_PERCENT
	TK_EQUAL
	TK_EQUAL2
	TK_NOT_EQUAL
	TK_LESS
	TK_GREATER
	TK_LESS_EQU
	TK_GREATER_EQU
	TK_AND
	TK_OR
	TK_XOR
	TK_NOT
	TK_AND_AND
	TK_OR_OR
	TK_SEMICOLON
	TK_COMMA
	TK_LPAREN
	TK_RPAREN
	TK_LBRACE
	TK_RBRACE
	TK_LBRACKET
	TK_RBRACKET
	TK_DOT
	TK_ARROW
	TK_PLUS_PLUS
	TK_MINUS_MINUS
	TK_PLUS_EQU
	TK_MINUS_EQU
	TK_STAR_EQU
	TK_SLASH_EQU
	TK_QUESTION
	TK_COLON
	TK_TILDE
	TK_SHIFT_LEFT
	TK_SHIFT_RIGHT

	// Keywords
	KW_IF
	KW_ELSE
	KW_FOR
	KW_WHILE
	KW_DO
	KW_RETURN
	KW_BREAK
	KW_CONTINUE
	KW_CLASS
	KW_UNION
	KW_ENUM
	KW_EXTERN
	KW_PUBLIC
	KW_STATIC
	KW_VOID
	KW_U8
	KW_U16
	KW_U32
	KW_U64
	KW_I8
	KW_I16
	KW_I32
	KW_I64
	KW_F64
	KW_BOOL
	KW_U0
	KW_ASM
	KW_SIZEOF
	KW_INCLUDE
	KW_DEFINE
	KW_IFDEF
	KW_IFNDEF
	KW_ENDIF
	KW_SWITCH
	KW_CASE
	KW_DEFAULT
	KW_GOTO
	KW_TRY
	KW_CATCH
)

// Token represents a lexical token
type Token struct {
	Type    int
	Line    int
	Col     int
	I64Val  int64
	F64Val  float64
	StrVal  string
	Ident   string
}

// Lexer performs lexical analysis
type Lexer struct {
	src      string
	pos      int
	line     int
	col      int
	token    Token
	prevToken Token
	currentToken Token // Alias for token
}

// NewLexer creates a new lexer
func NewLexer(src string) *Lexer {
	return &Lexer{
		src:  src,
		line: 1,
		col:  0,
	}
}

// skipWhitespace skips whitespace and comments
func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.src) {
		ch := l.src[l.pos]
		if ch == ' ' || ch == '\t' {
			l.pos++
			l.col++
		} else if ch == '\n' {
			l.pos++
			l.line++
			l.col = 0
		} else if ch == '/' && l.pos+1 < len(l.src) && l.src[l.pos+1] == '/' {
			// Single line comment
			for l.pos < len(l.src) && l.src[l.pos] != '\n' {
				l.pos++
			}
		} else if ch == '/' && l.pos+1 < len(l.src) && l.src[l.pos+1] == '*' {
			// Multi-line comment
			l.pos += 2
			for l.pos+1 < len(l.src) && !(l.src[l.pos] == '*' && l.src[l.pos+1] == '/') {
				if l.src[l.pos] == '\n' {
					l.line++
					l.col = 0
				}
				l.pos++
			}
			if l.pos+1 < len(l.src) {
				l.pos += 2
			}
		} else {
			break
		}
	}
}

// isAlphaNumeric checks if a character is alphanumeric or underscore
func isAlphaNumeric(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_' || unicode.IsDigit(rune(ch))
}

// isDigit checks if a character is a digit
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// isHexDigit checks if a character is a hex digit
func isHexDigit(ch byte) bool {
	return isDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

// NextToken fetches the next token
func (l *Lexer) NextToken() Token {
	l.prevToken = l.token
	l.skipWhitespace()

	l.token.Line = l.line
	l.token.Col = l.col

	if l.pos >= len(l.src) {
		l.token.Type = TK_EOF
		l.currentToken = l.token
		return l.token
	}

	ch := l.src[l.pos]

	// Identifiers and keywords
	if unicode.IsLetter(rune(ch)) || ch == '_' {
		start := l.pos
		for l.pos < len(l.src) && isAlphaNumeric(l.src[l.pos]) {
			l.pos++
		}
		ident := l.src[start:l.pos]
		l.token.Ident = ident
		l.token.Type = TK_IDENT
		l.col += l.pos - start

		// Check for keywords
		keywords := map[string]int{
			"if":       KW_IF,
			"else":     KW_ELSE,
			"for":      KW_FOR,
			"while":    KW_WHILE,
			"do":       KW_DO,
			"return":   KW_RETURN,
			"break":    KW_BREAK,
			"continue": KW_CONTINUE,
			"class":    KW_CLASS,
			"union":    KW_UNION,
			"enum":     KW_ENUM,
			"extern":   KW_EXTERN,
			"public":   KW_PUBLIC,
			"static":   KW_STATIC,
			"void":     KW_VOID,
			"U8":       KW_U8,
			"U16":      KW_U16,
			"U32":      KW_U32,
			"U64":      KW_U64,
			"I8":       KW_I8,
			"I16":      KW_I16,
			"I32":      KW_I32,
			"I64":      KW_I64,
			"F64":      KW_F64,
			"Bool":     KW_BOOL,
			"U0":       KW_U0,
			"asm":      KW_ASM,
			"sizeof":   KW_SIZEOF,
			"include":  KW_INCLUDE,
			"define":   KW_DEFINE,
			"ifdef":    KW_IFDEF,
			"ifndef":   KW_IFNDEF,
			"endif":    KW_ENDIF,
			"switch":   KW_SWITCH,
			"case":     KW_CASE,
			"default":  KW_DEFAULT,
			"goto":     KW_GOTO,
			"try":      KW_TRY,
			"catch":    KW_CATCH,
		}

		if kw, ok := keywords[ident]; ok {
			l.token.Type = kw
		}

		l.currentToken = l.token
		return l.token
	}

	// Numbers
	if isDigit(ch) {
		l.token.I64Val = 0
		l.token.Type = TK_I64
		start := l.pos

		// Hex
		if ch == '0' && l.pos+1 < len(l.src) && (l.src[l.pos+1] == 'x' || l.src[l.pos+1] == 'X') {
			l.pos += 2
			for l.pos < len(l.src) && isHexDigit(l.src[l.pos]) {
				digit := l.src[l.pos]
				if digit >= '0' && digit <= '9' {
					l.token.I64Val = l.token.I64Val*16 + int64(digit-'0')
				} else {
					digit = toUpper(digit)
					l.token.I64Val = l.token.I64Val*16 + int64(digit-'A'+10)
				}
				l.pos++
			}
		} else if ch == '0' && l.pos+1 < len(l.src) && l.src[l.pos+1] == 'b' {
			// Binary
			l.pos += 2
			for l.pos < len(l.src) && (l.src[l.pos] == '0' || l.src[l.pos] == '1') {
				l.token.I64Val = l.token.I64Val*2 + int64(l.src[l.pos]-'0')
				l.pos++
			}
		} else {
			// Decimal
			for l.pos < len(l.src) && isDigit(l.src[l.pos]) {
				l.token.I64Val = l.token.I64Val*10 + int64(l.src[l.pos]-'0')
				l.pos++
			}

			// Float
			if l.pos < len(l.src) && l.src[l.pos] == '.' {
				l.token.F64Val = float64(l.token.I64Val)
				l.token.Type = TK_F64
				l.pos++
				decimalPlaces := 0
				for l.pos < len(l.src) && isDigit(l.src[l.pos]) {
					l.token.F64Val = l.token.F64Val*10 + float64(l.src[l.pos]-'0')
					decimalPlaces++
					l.pos++
				}
				for decimalPlaces > 0 {
					l.token.F64Val /= 10
					decimalPlaces--
				}
			}
		}

		l.col += l.pos - start
		return l.token
	}

	// String
	if ch == '"' {
		l.pos++
		start := l.pos
		var str []byte
		for l.pos < len(l.src) && l.src[l.pos] != '"' {
			if l.src[l.pos] == '\\' && l.pos+1 < len(l.src) {
				l.pos++
				switch l.src[l.pos] {
				case 'n':
					str = append(str, '\n')
				case 'r':
					str = append(str, '\r')
				case 't':
					str = append(str, '\t')
				case '\\':
					str = append(str, '\\')
				case '"':
					str = append(str, '"')
				case '0':
					str = append(str, 0)
				default:
					str = append(str, l.src[l.pos])
				}
			} else {
				str = append(str, l.src[l.pos])
			}
			l.pos++
		}
		if l.pos < len(l.src) {
			l.pos++ // Skip closing quote
		}
		l.token.Type = TK_STR
		l.token.StrVal = string(str)
		l.col += l.pos - start + 2
		return l.token
	}

	// Char constant
	if ch == '\'' {
		l.pos++
		l.token.I64Val = 0
		if l.pos < len(l.src) {
			if l.src[l.pos] == '\\' && l.pos+1 < len(l.src) {
				l.pos++
				switch l.src[l.pos] {
				case 'n':
					l.token.I64Val = '\n'
				case 'r':
					l.token.I64Val = '\r'
				case 't':
					l.token.I64Val = '\t'
				case '\\':
					l.token.I64Val = '\\'
				case '\'':
					l.token.I64Val = '\''
				case '0':
					l.token.I64Val = 0
				default:
					l.token.I64Val = int64(l.src[l.pos])
				}
				l.pos++
			} else {
				l.token.I64Val = int64(l.src[l.pos])
				l.pos++
			}
		}
		if l.pos < len(l.src) && l.src[l.pos] == '\'' {
			l.pos++
		}
		l.token.Type = TK_CHAR_CONST
		l.col += 2
		return l.token
	}

	// Operators
	start := l.pos
	switch ch {
	case '+':
		if l.pos+1 < len(l.src) && l.src[l.pos+1] == '+' {
			l.token.Type = TK_PLUS_PLUS
			l.pos += 2
		} else if l.pos+1 < len(l.src) && l.src[l.pos+1] == '=' {
			l.token.Type = TK_PLUS_EQU
			l.pos += 2
		} else {
			l.token.Type = TK_PLUS
			l.pos++
		}
	case '-':
		if l.pos+1 < len(l.src) && l.src[l.pos+1] == '-' {
			l.token.Type = TK_MINUS_MINUS
			l.pos += 2
		} else if l.pos+1 < len(l.src) && l.src[l.pos+1] == '=' {
			l.token.Type = TK_MINUS_EQU
			l.pos += 2
		} else if l.pos+1 < len(l.src) && l.src[l.pos+1] == '>' {
			l.token.Type = TK_ARROW
			l.pos += 2
		} else {
			l.token.Type = TK_MINUS
			l.pos++
		}
	case '*':
		if l.pos+1 < len(l.src) && l.src[l.pos+1] == '=' {
			l.token.Type = TK_STAR_EQU
			l.pos += 2
		} else {
			l.token.Type = TK_STAR
			l.pos++
		}
	case '/':
		if l.pos+1 < len(l.src) && l.src[l.pos+1] == '=' {
			l.token.Type = TK_SLASH_EQU
			l.pos += 2
		} else {
			l.token.Type = TK_SLASH
			l.pos++
		}
	case '%':
		l.token.Type = TK_PERCENT
		l.pos++
	case '=':
		if l.pos+1 < len(l.src) && l.src[l.pos+1] == '=' {
			l.token.Type = TK_EQUAL2
			l.pos += 2
		} else {
			l.token.Type = TK_EQUAL
			l.pos++
		}
	case '!':
		if l.pos+1 < len(l.src) && l.src[l.pos+1] == '=' {
			l.token.Type = TK_NOT_EQUAL
			l.pos += 2
		} else {
			l.token.Type = TK_NOT
			l.pos++
		}
	case '<':
		if l.pos+1 < len(l.src) && l.src[l.pos+1] == '=' {
			l.token.Type = TK_LESS_EQU
			l.pos += 2
		} else if l.pos+1 < len(l.src) && l.src[l.pos+1] == '<' {
			l.token.Type = TK_SHIFT_LEFT
			l.pos += 2
		} else {
			l.token.Type = TK_LESS
			l.pos++
		}
	case '>':
		if l.pos+1 < len(l.src) && l.src[l.pos+1] == '=' {
			l.token.Type = TK_GREATER_EQU
			l.pos += 2
		} else if l.pos+1 < len(l.src) && l.src[l.pos+1] == '>' {
			l.token.Type = TK_SHIFT_RIGHT
			l.pos += 2
		} else {
			l.token.Type = TK_GREATER
			l.pos++
		}
	case '&':
		if l.pos+1 < len(l.src) && l.src[l.pos+1] == '&' {
			l.token.Type = TK_AND_AND
			l.pos += 2
		} else {
			l.token.Type = TK_AND
			l.pos++
		}
	case '|':
		if l.pos+1 < len(l.src) && l.src[l.pos+1] == '|' {
			l.token.Type = TK_OR_OR
			l.pos += 2
		} else {
			l.token.Type = TK_OR
			l.pos++
		}
	case '^':
		l.token.Type = TK_XOR
		l.pos++
	case '~':
		l.token.Type = TK_TILDE
		l.pos++
	case ';':
		l.token.Type = TK_SEMICOLON
		l.pos++
	case ',':
		l.token.Type = TK_COMMA
		l.pos++
	case '(':
		l.token.Type = TK_LPAREN
		l.pos++
	case ')':
		l.token.Type = TK_RPAREN
		l.pos++
	case '{':
		l.token.Type = TK_LBRACE
		l.pos++
	case '}':
		l.token.Type = TK_RBRACE
		l.pos++
	case '[':
		l.token.Type = TK_LBRACKET
		l.pos++
	case ']':
		l.token.Type = TK_RBRACKET
		l.pos++
	case '.':
		l.token.Type = TK_DOT
		l.pos++
	case '?':
		l.token.Type = TK_QUESTION
		l.pos++
	case ':':
		l.token.Type = TK_COLON
		l.pos++
	default:
		fmt.Printf("Unknown character '%c' at line %d, col %d\n", ch, l.line, l.col)
		l.pos++
		l.token.Type = TK_EOF
	}

	l.col += l.pos - start
	l.currentToken = l.token
	return l.token
}

// toUpper converts a byte to uppercase
func toUpper(ch byte) byte {
	if ch >= 'a' && ch <= 'z' {
		return ch - 'a' + 'A'
	}
	return ch
}

// Match checks if current token matches expected type and advances
func (l *Lexer) Match(expected int) error {
	if l.token.Type != expected {
		return fmt.Errorf("expected token type %d, got %d at line %d", expected, l.token.Type, l.line)
	}
	l.NextToken()
	return nil
}
