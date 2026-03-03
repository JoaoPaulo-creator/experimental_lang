package parser

import (
	"experimental/ast"
	"experimental/lexer"
	"fmt"
	"slices"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) peek() lexer.Token {
	return p.tokens[p.pos]
}

func (p *Parser) advance() lexer.Token {
	t := p.tokens[p.pos]
	p.pos++
	return t
}

func (p *Parser) expect(tokenType lexer.Kind) (lexer.Token, error) {
	t := p.advance()
	if t.Kind != tokenType {
		return t, fmt.Errorf("expected %s, got %s (%q)", lexer.TokenKindString(tokenType), lexer.TokenKindString(t.Kind), t.Literal)
	}

	return t, nil
}

func (p *Parser) at(types ...lexer.Kind) bool {
	cur := p.peek().Kind
	return slices.Contains(types, cur)
}

func (p *Parser) skipNewLines() {
	for p.at(lexer.NEW_LINE) {
		p.advance()
	}
}

// --- binding power table
func (p *Parser) infixBP(tokenKind lexer.Kind) int {
	return 0
}

func (p *Parser) parserExpr(minBP int) (ast.Expr, error) {
	p.skipNewLines()
	left, err := p.parsePrefix()
	if err != nil {
		return nil, err
	}

	for {
		bp := p.infixBP(p.peek().Kind)
		if bp <= minBP {
			break
		}

		op := p.advance()
		right, err := p.parserExpr(bp)
		if err != nil {
			return nil, err
		}

		left = ast.BinOp{Op: op.Literal, Left: left, Right: right}
	}

	return left, nil
}

func (p *Parser) parsePrefix() (ast.Expr, error) {
	t := p.peek()

	switch t.Kind {
	case lexer.MATCH:
		return p.parseMatch()
	case lexer.IDENTIFIER:
		p.advance()
		return ast.Ident{Name: t.Literal}, nil
	case lexer.STRING:
		p.advance()
		return ast.StrLit{Value: t.Literal}, nil
	case lexer.NUMBER:
		p.advance()
		return ast.NumLit{Value: t.Literal}, nil
	case lexer.UNDER_SCORE:
		p.advance()
		return ast.WildcardExpr{}, nil
	}

	return nil, fmt.Errorf("unexpected token %s %q", lexer.TokenKindString(t.Kind), t.Literal)
}

func (p *Parser) parseMatch() (ast.Expr, error) {
	if _, err := p.expect(lexer.MATCH); err != nil {
		return nil, err
	}

	subject, err := p.parserExpr(0)
	if err != nil {
		return nil, err
	}

	var arms []ast.MatchArm
	p.skipNewLines()
	for p.at(lexer.PIPE) {
		p.advance() // consume '|'

		pattern, err := p.parsePrefix()
		if err != nil {
			return nil, err
		}

		if _, err := p.expect(lexer.THEN); err != nil {
			return nil, err
		}

		body, err := p.parserExpr(0)
		if err != nil {
			return nil, err
		}

		arms = append(arms, ast.MatchArm{Pattern: pattern, Body: body})
		p.skipNewLines()
	}

	return ast.MatchExpr{Subject: subject, Arms: arms}, nil
}

// let declaration
func (p *Parser) parseLet() (ast.LetDecl, error) {
	if _, err := p.expect(lexer.LET); err != nil {
		return ast.LetDecl{}, err
	}

	nameTok, err := p.expect(lexer.IDENTIFIER)
	if err != nil {
		return ast.LetDecl{}, err
	}

	var params []string
	for !p.at(lexer.FUNCTION_ASSOCIATION, lexer.EOF, lexer.NEW_LINE) {
		t, err := p.expect(lexer.IDENTIFIER)
		if err != nil {
			return ast.LetDecl{}, err
		}

		params = append(params, t.Literal)
	}
	p.skipNewLines()

	if _, err := p.expect(lexer.FUNCTION_ASSOCIATION); err != nil {
		return ast.LetDecl{}, err
	}

	body, err := p.parserExpr(0)
	if err != nil {
		return ast.LetDecl{}, err
	}

	return ast.LetDecl{Name: nameTok.Literal, Params: params, Body: body}, nil
}

// -- program
func (p *Parser) ParseProgram() (ast.Program, error) {
	var decls []ast.LetDecl
	p.skipNewLines()

	for !p.at(lexer.EOF) {
		switch {
		case p.at(lexer.LET):
			d, err := p.parseLet()
			if err != nil {
				return ast.Program{}, err
			}
			decls = append(decls, d)
			p.skipNewLines()
		default:
			t := p.peek()
			return ast.Program{}, fmt.Errorf("unexpected token %s (%q)", t.Kind, t.Literal)
		}
	}

	return ast.Program{Decls: decls}, nil
}
