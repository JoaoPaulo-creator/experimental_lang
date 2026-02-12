package lexer

import "fmt"

type Kind int

const (
	EOF Kind = iota
	ILLEGAL

	// keywords
	LET
	MATCH
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
	PIPE // |
	COMMA
	SEMI_COLON

	// operators
	ASSIGNMENT // =
	EQ         // == (optional; lexer below doesn't emit yet)
	PLUS
	PLUS_PLUS
	MINUS
	MINUS_MINUS

	// patterns
	UNDER_SCORE // wildcard _
)

type Token struct {
	Kind    Kind
	Literal string
}

var keywords = map[string]Kind{
	"let":    LET,
	"match":  MATCH,
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

	case LET:
		return "let"
	case MATCH:
		return "match"
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

	case ASSIGNMENT:
		return "assignment"
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

func newToken(kind Kind, value string) Token {
	return Token{
		Kind:    kind,
		Literal: value,
	}
}
