package lexer

import (
	"fmt"
	"unicode"
)

type lexer struct {
	source string
	pos    int
	line   int
	Tokens []Token
}

func Tokenize(source string) []Token {
	lex := &lexer{
		source: source,
		pos:    0,
		line:   1,
		Tokens: make([]Token, 0),
	}

	for !lex.atEOF() {
		lex.scanToken()
	}

	lex.Tokens = append(lex.Tokens, newToken(EOF, " "))

	return lex.Tokens
}

func (lex *lexer) scanToken() {
	ch := lex.peek()

	// newline
	if ch == '\n' {
		lex.advance()
		lex.push(newToken(NEW_LINE, "\n"))
		lex.line++
		return
	}

	// skip whitespaces (excluding newline)
	if unicode.IsSpace(rune(ch)) {
		lex.skipWhitespace()
		return
	}

	// comments: //
	if ch == '/' && lex.peekNext() == '/' {
		lex.skipComment()
		return
	}

	// string literals
	if ch == '"' {
		lex.scanString()
		return
	}

	// numbers
	if unicode.IsDigit(rune(ch)) {
		lex.scanNumber()
		return
	}

	// '_' can be wildcard OR start of identifier
	if ch == '_' {
		next := lex.peekNext()
		// If next continues an identifier, treat as identifier start (e.g. _x, __foo)
		if isIdentContinue(next) {
			lex.scanIdentifier()
		} else {
			lex.advance()
			lex.push(newToken(UNDER_SCORE, "_"))
		}
		return
	}

	// identifiers and keywords (snake_case)
	if isIdentStart(ch) {
		lex.scanIdentifier()
		return
	}

	switch ch {
	case '=':
		lex.advance()
		lex.push(newToken(FUNCTION_ASSOCIATION, "="))
		return
	case '{':
		lex.advance()
		lex.push(newToken(OPEN_CURLY, "{"))
		return
	case '}':
		lex.advance()
		lex.push(newToken(CLOSE_CURLY, "}"))
		return
	case '[':
		lex.advance()
		lex.push(newToken(OPEN_BRACKET, "["))
		return
	case ']':
		lex.advance()
		lex.push(newToken(CLOSE_BRACKET, "]"))
		return
	case '(':
		lex.advance()
		lex.push(newToken(OPEN_PARENTHESIS, "("))
		return
	case ')':
		lex.advance()
		lex.push(newToken(CLOSE_PARENTHESIS, ")"))
		return
	case ':':
		lex.advance()
		lex.push(newToken(COLON, ":"))
		return
	case ';':
		lex.advance()
		lex.push(newToken(SEMI_COLON, ";"))
		return
	case ',':
		lex.advance()
		lex.push(newToken(COMMA, ","))
		return
	case '|':
		lex.advance()
		lex.push(newToken(PIPE, "|"))
		return
	case '-':
		lex.advance()
		lex.push(newToken(MINUS, "-"))
		return
	case '+':
		lex.advance()
		lex.push(newToken(PLUS, "+"))
		return
	case '\n':
		lex.advance()
		lex.push(newToken(NEW_LINE, "\n"))
		return
	}

	panic(fmt.Sprintf("lexer error: unexpected character '%c' at position %d (line %d)", ch, lex.pos, lex.line))
}

func (lex *lexer) scanString() {
	start := lex.pos
	lex.advance() // consume opening "

	for !lex.atEOF() && lex.peek() != '"' {
		lex.advance()
	}

	if lex.atEOF() {
		panic("lexer error: unterminated string literal")
	}

	lex.advance()
	value := lex.source[start:lex.pos]
	lex.push(newToken(STRING, value))
}

func (lex *lexer) scanNumber() {
	start := lex.pos
	lex.advance()

	for !lex.atEOF() && unicode.IsDigit(rune(lex.peek())) {
		lex.advance()
	}

	if !lex.atEOF() && lex.peek() == '.' && unicode.IsDigit(rune(lex.peekNext())) {
		lex.advance() // consume '.'
		for !lex.atEOF() && unicode.IsDigit(rune(lex.peek())) {
			lex.advance()
		}
	}

	value := lex.source[start:lex.pos]
	lex.push(newToken(NUMBER, value))
}

func (lex *lexer) scanIdentifier() {
	start := lex.pos
	lex.advance()

	for !lex.atEOF() && isIdentContinue(lex.peek()) {
		lex.advance()
	}

	value := lex.source[start:lex.pos]
	if kind, found := keywords[value]; found {
		lex.push(newToken(kind, value))
	} else {
		lex.push(newToken(IDENTIFIER, value))
	}
}

func isIdentStart(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

func isIdentContinue(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_'
}

func (lex *lexer) skipWhitespace() {
	for !lex.atEOF() && unicode.IsSpace(rune(lex.peek())) && lex.peek() != '\n' {
		if lex.peek() == '\t' {
			lex.pos++
		}
		lex.advance()
	}
}

func (lex *lexer) skipComment() {
	for !lex.atEOF() && lex.peek() != '\n' {
		lex.advance()
	}

	if !lex.atEOF() {
		lex.advance()
		lex.line++
	}
}

func (lex *lexer) peek() byte {
	if lex.atEOF() {
		return 0
	}

	return lex.source[lex.pos]
}

func (lex *lexer) peekNext() byte {
	if lex.pos+1 >= len(lex.source) {
		return 0
	}

	return lex.source[lex.pos+1]
}

func (lex *lexer) peekAhead(n int) byte {
	if lex.pos+n >= len(lex.source) {
		return 0
	}

	return lex.source[lex.pos+n]
}

func (lex *lexer) advance() {
	lex.pos++
}

func (lex *lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *lexer) atEOF() bool {
	return lex.pos >= len(lex.source)
}
