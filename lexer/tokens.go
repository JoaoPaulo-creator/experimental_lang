package lexer

import "fmt"

type Kind int

const (
	EOF Kind = iota
	ILLEGAL
	GIVEN
	LET
	SET
	IDENTIFIER
	OPEN_BRACKET
	CLOSE_BRACKET

	INT
	STRING
	NUMBER

	ASSIGNMENT
	EQ
	COMMA
	SEMI_COLON
)

type Token struct {
	Kind    Kind
	Literal string
}

var keywords = map[string]Kind{
	"given": GIVEN,
	"set":   SET,
}

func TokenKindString(kind Kind) string {
	switch kind {
	case EOF:
		return "eof"
	case LET:
		return "let"
	case GIVEN:
		return "given"
	case SET:
		return "set"
	case IDENTIFIER:
		return "identifier"
	case STRING:
		return "string"
	case NUMBER:
		return "number"
	case ASSIGNMENT:
		return "assignment"
	case SEMI_COLON:
		return "semi_colon"
	case OPEN_BRACKET:
		return "open_bracket"
	case CLOSE_BRACKET:
		return "close_bracket"
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
