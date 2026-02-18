package parser

import (
	"experimental/ast"
	"experimental/lexer"
	"fmt"
	"strings"
)

type parser struct {
	tokens []lexer.Token
	pos    int
}

func createParser(tokens []lexer.Token) *parser {
	createTokenLookups()

	return &parser{
		tokens: tokens,
	}
}

type BTParser struct {
	tokens      []lexer.Token
	pos         int
	farthest    int
	farthestErr error
}

func NewBTParser(tokens []lexer.Token) *BTParser {
	return &BTParser{
		tokens: tokens,
	}
}

func (p *BTParser) ParseProgram() (*ast.Program, error) {
	prog, ok := p.program()
	if ok && p.at(lexer.EOF) {
		return prog, nil
	}

	if p.farthestErr != nil {
		return nil, p.farthestErr
	}

	return nil, fmt.Errorf("parse error at token %d: %v", p.pos, p.cur())
}

func (p *BTParser) cur() lexer.Token {
	return p.tokens[p.pos]
}

func (p *BTParser) at(t lexer.Kind) bool {
	return p.cur().Kind == t
}

func (p *BTParser) accept(k lexer.Kind) bool {
	if p.at(k) {
		p.pos++
		return true
	}

	return false
}

func (p *BTParser) setError(msg string) {
	if p.pos >= p.farthest {
		p.farthest = p.pos
		p.farthestErr = fmt.Errorf("parse error at token %d (%s): %s",
			p.pos, lexer.TokenKindString(p.cur().Kind), msg)
	}
}

func (p *BTParser) skipNewLines() {
	for p.at(lexer.NEW_LINE) {
		p.pos++
	}
}

func (p *BTParser) program() (*ast.Program, bool) {
	start := p.pos
	prog := &ast.Program{}
	p.skipNewLines()

	for {
		if !p.at(lexer.FUN) {
			break
		}

		fn, ok := p.funcDecl()
		if !ok {
			p.pos = start
			return nil, false
		}

		prog.Func = append(prog.Func, fn)
		p.skipNewLines()
	}

	return prog, true
}

func (p *BTParser) funcDecl() (*ast.FunctionDeclStmt, bool) {
	start := p.pos
	if !p.accept(lexer.FUN) {
		p.pos = start
		return nil, false
	}

	if !p.at(lexer.IDENTIFIER) {
		p.setError("expected function name")
		p.pos = start
		return nil, false
	}

	name := p.cur().Literal
	p.pos++

	if !p.accept(lexer.OPEN_PARENTHESIS) {
		p.setError("expected '(' after function identifier")
		p.pos = start
		return nil, false
	}

	params, ok := p.paramList()
	if !ok {
		p.pos = start
		return nil, false
	}

	if !p.accept(lexer.CLOSE_PARENTHESIS) {
		p.setError("expected ')'")
		p.pos = start
		return nil, false
	}

	body, ok := p.block()
	if !ok {
		p.pos = start
		return nil, false
	}

	fn := &ast.FunctionDeclStmt{
		Name:       name,
		Parameters: params,
		Body:       body,
		Visibility: "public",
	}

	return fn, true
}

func (p *BTParser) paramList() ([]ast.Parameter, bool) {
	start := p.pos

	if !p.at(lexer.IDENTIFIER) {
		return []ast.Parameter{}, true
	}

	var params []ast.Parameter
	name := p.cur().Literal
	p.pos++
	params = append(params, ast.Parameter{Name: name})

	for p.accept(lexer.COMMA) {
		if !p.at(lexer.IDENTIFIER) {
			p.setError("expected paramter name after ','")
			p.pos = start
			return nil, false
		}

		name = p.cur().Literal
		p.pos++
		params = append(params, ast.Parameter{Name: name})
	}

	return params, true
}

func (p *BTParser) printLnStmt() (ast.Stmt, bool) {
	start := p.pos

	if !(p.accept(lexer.PRINT) || (p.at(lexer.IDENTIFIER) && p.cur().Literal == "println")) {
		p.pos = start
		return nil, false
	}

	if p.cur().Kind == lexer.IDENTIFIER && p.cur().Literal == "println" {
		p.pos++
	}

	if !p.accept(lexer.OPEN_PARENTHESIS) || !(p.at(lexer.STRING) || p.at(lexer.IDENTIFIER)) {
		p.setError("expected string literal or identifier after print(")
		p.pos = start
		return nil, false
	}

	raw := p.cur().Literal
	isString := p.at(lexer.STRING)
	p.pos++

	if !p.accept(lexer.CLOSE_PARENTHESIS) {
		p.setError("expected ')'")
		p.pos = start
		return nil, false
	}

	// optional separators
	if p.accept(lexer.SEMI_COLON) || p.accept(lexer.NEW_LINE) || p.at(lexer.EOF) {
		// ok
	}

	val := raw
	if isString {
		val = strings.Trim(raw, `"`)
	}
	stmt := ast.PrintLnStmt{Value: val}
	return stmt, true
}

func (p *BTParser) stmt() (ast.Stmt, bool) {
	return p.printLnStmt()
}

func (p *BTParser) block() ([]ast.Stmt, bool) {
	start := p.pos

	if !p.accept(lexer.OPEN_CURLY) {
		p.pos = start
		return nil, false
	}

	var stmts []ast.Stmt
	p.skipNewLines()
	for !p.at(lexer.CLOSE_CURLY) && !p.at(lexer.EOF) {
		st, ok := p.stmt()
		if !ok {
			p.pos = start
			return nil, false
		}
		stmts = append(stmts, st)
		p.skipNewLines()
	}

	if !p.accept(lexer.CLOSE_CURLY) {
		p.setError("expected '}'")
		p.pos = start
		return nil, false
	}

	return stmts, true
}
