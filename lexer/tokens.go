package lexer

import (
	"fmt"
	"slices"
)

type Kind int

const (
	EOF Kind = iota
	ILLEGAL

	// keywords
	PUBLIC
	FUN
	PRINT
	LET
	MATCH
	WITH
	THEN
	MODULE
	NIL

	// identifiers + literals
	IDENTIFIER
	STRING
	NUMBER

	// punctuation
	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_PARENTHESIS
	CLOSE_PARENTHESIS
	OPEN_CURLY
	CLOSE_CURLY
	PIPE // |
	COMMA
	SEMI_COLON
	COLON

	// operators
	FUNCTION_ASSOCIATION // =
	EQ                   // == (optional; lexer below doesn't emit yet)
	PLUS
	PLUS_PLUS
	MINUS
	MINUS_MINUS

	NEW_LINE

	// patterns
	UNDER_SCORE // wildcard _
)

type Token struct {
	Kind    Kind
	Literal string
}

var keywords = map[string]Kind{
	"public": PUBLIC,
	"fun":    FUN,
	"print":  PRINT,
	"let":    LET,
	"match":  MATCH,
	"with":   WITH,
	"then":   THEN,
	"nil":    NIL,
	"module": MODULE,
}

func TokenKindString(kind Kind) string {
	switch kind {
	case EOF:
		return "eof"
	case ILLEGAL:
		return "illegal"
	case PUBLIC:
		return "public"
	case FUN:
		return "fun"
	case PRINT:
		return "print"
	case LET:
		return "let"
	case MATCH:
		return "match"
	case WITH:
		return "with"
	case THEN:
		return "then"
	case MODULE:
		return "module"
	case NIL:
		return "nil"

	case IDENTIFIER:
		return "identifier"
	case STRING:
		return "string"
	case NUMBER:
		return "number"

	case OPEN_BRACKET:
		return "open_bracket"
	case CLOSE_BRACKET:
		return "close_bracket"
	case OPEN_PARENTHESIS:
		return "open_parenthesis"
	case CLOSE_PARENTHESIS:
		return "close_parenthesis"
	case PIPE:
		return "pipe"
	case COMMA:
		return "comma"
	case SEMI_COLON:
		return "semi_colon"
	case COLON:
		return "colon"
	case FUNCTION_ASSOCIATION:
		return "function_association"
	case NEW_LINE:
		return "new_line"
	case EQ:
		return "eq"
	case PLUS:
		return "plus"
	case PLUS_PLUS:
		return "plus_plus"
	case MINUS:
		return "minus"
	case MINUS_MINUS:
		return "minus_minus"

	case UNDER_SCORE:
		return "under_score"

	default:
		return fmt.Sprintf("unknown token: %d", kind)
	}
}

func (token Token) IsOneOfMany(expectedTokens ...Kind) bool {
	return slices.Contains(expectedTokens, token.Kind)
}

func newToken(kind Kind, value string) Token {
	return Token{
		Kind:    kind,
		Literal: value,
	}
}
