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

type PEGParser struct {
	tokens      []lexer.Token
	pos         int
	memo        map[memoKey]memoVal
	farthest    int
	farthestErr error
}

type memoKey struct {
	rule string
	pos  int
}

type memoVal struct {
	ok  bool
	pos int
	val any
}

func NewPEGParser(tokens []lexer.Token) *PEGParser {
	return &PEGParser{
		tokens: tokens,
		memo:   map[memoKey]memoVal{},
	}
}

func (p *PEGParser) ParseProgram() (*ast.Program, error) {
	prog, ok := p.program()
	if ok && p.at(lexer.EOF) {
		return prog.(*ast.Program), nil
	}

	if p.farthestErr != nil {
		return nil, p.farthestErr
	}

	return nil, fmt.Errorf("parse error at token %d: %v", p.pos, p.cur())
}

func (p *PEGParser) cur() lexer.Token {
	return p.tokens[p.pos]
}

func (p *PEGParser) at(t lexer.Kind) bool {
	return p.cur().Kind == t
}

func (p *PEGParser) accept(k lexer.Kind) bool {
	if p.at(k) {
		p.pos++
		return true
	}

	return false
}

func (p *PEGParser) setError(msg string) {
	if p.pos >= p.farthest {
		p.farthest = p.pos
		p.farthestErr = fmt.Errorf("parse error at token %d (%s): %s",
			p.pos, lexer.TokenKindString(p.cur().Kind), msg)
	}
}

func (p *PEGParser) skipNewLines() {
	for p.at(lexer.NEW_LINE) {
		p.pos++
	}
}

func (p *PEGParser) memoized(rule string, pos int) (memoVal, bool) {
	v, ok := p.memo[memoKey{rule, pos}]
	if ok {
		p.pos = v.pos
	}

	return v, ok
}

func (p *PEGParser) memoize(rule string, start int, v memoVal) {
	p.memo[memoKey{rule, start}] = v
}

func (p *PEGParser) program() (any, bool) {
	start := p.pos
	if v, ok := p.memoized("Program", start); ok {
		if v.ok {
			return v.val.(*ast.Program), true
		}
		return nil, false
	}

	prog := &ast.Program{}
	p.skipNewLines()

	for p.at(lexer.FUN) {
		if !p.at(lexer.FUN) {
			break
		}

		fnAny, ok := p.funcDecl()
		if !ok {
			p.pos = start
			p.memoize("Program", start, memoVal{ok: false, pos: start})
			return nil, false
		}

		prog.Func = append(prog.Func, fnAny.(*ast.FunctionDeclStmt))
		p.skipNewLines()
	}

	p.memoize("Program", start, memoVal{ok: true, pos: p.pos, val: prog})
	return prog, true
}

func (p *PEGParser) funcDecl() (any, bool) {
	start := p.pos
	if v, ok := p.memoized("FunDecl", start); ok {
		if v.ok {
			return v.val.(*ast.FunctionDeclStmt), true
		}
		return nil, false
	}

	if !p.accept(lexer.FUN) {
		p.memoize("FunDecl", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	if !p.accept(lexer.IDENTIFIER) {
		p.setError("expected function name")
		p.pos = start
		p.memoize("FunDecl", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	name := p.cur().Literal
	p.pos++
	if !p.accept(lexer.OPEN_PARENTHESIS) && p.accept(lexer.CLOSE_PARENTHESIS) {
		p.pos = start
		p.memoize("FunDecl", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	bodyAny, ok := p.block()
	if !ok {
		p.pos = start
		p.memoize("FunDecl", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	fn := &ast.FunctionDeclStmt{
		Name:       name,
		Parameters: []ast.Parameter{},
		Body:       bodyAny.([]ast.Stmt),
	}

	p.memoize("FunDecl", start, memoVal{ok: true, pos: p.pos, val: fn})
	return fn, true
}

func (p *PEGParser) printLnStmt() (any, bool) {
	start := p.pos
	if v, ok := p.memoized("PrintStmt", start); ok {
		return v.val, v.ok
	}

	if !(p.accept(lexer.PRINT) || (p.at(lexer.IDENTIFIER) && p.cur().Literal == "println")) {
		p.memoize("PrintStmt", start, memoVal{ok: false, pos: start})
		return nil, false
	}
	if p.cur().Kind == lexer.IDENTIFIER && p.cur().Literal == "println" {
		p.pos++
	}

	if !p.accept(lexer.OPEN_PARENTHESIS) || !p.at(lexer.STRING) {
		p.setError("expected string literal after print(")
		p.pos = start
		p.memoize("PrintStmt", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	raw := p.cur().Literal
	p.pos++

	if !p.accept(lexer.CLOSE_PARENTHESIS) {
		p.setError("expected ')'")
		p.pos = start
		p.memoize("PrintStmt", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	// optional separators
	if p.accept(lexer.SEMI_COLON) || p.accept(lexer.NEW_LINE) || p.at(lexer.EOF) {
		// ok
	}

	val := strings.Trim(raw, `"`)
	stmt := ast.PrintLnStmt{Value: val}
	p.memoize("PrintStmt", start, memoVal{ok: true, pos: p.pos, val: stmt})
	return stmt, true
}

func (p *PEGParser) stmt() (any, bool) {
	start := p.pos
	if v, ok := p.memoized("Stmt", start); ok {
		return v.val, v.ok
	}

	// only PrintStmt for now
	st, ok := p.printLnStmt()
	if !ok {
		p.memoize("Stmt", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	p.memoize("Stmt", start, memoVal{ok: true, pos: p.pos, val: st})
	return st, true
}

func (p *PEGParser) block() (any, bool) {
	start := p.pos
	if v, ok := p.memoized("Block", start); ok {
		return v.val, v.ok
	}

	if !p.accept(lexer.OPEN_CURLY) {
		p.memoize("Block", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	var stmts []ast.Stmt
	p.skipNewLines()
	for !p.at(lexer.CLOSE_CURLY) && !p.at(lexer.EOF) {
		st, ok := p.stmt()
		if !ok {
			p.pos = start
			p.memoize("Block", start, memoVal{ok: false, pos: start})
			return nil, false
		}
		stmts = append(stmts, st.(ast.Stmt))
		p.skipNewLines()
	}

	if !p.accept(lexer.CLOSE_CURLY) {
		p.setError("expected '}'")
		p.pos = start
		p.memoize("Block", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	p.memoize("Block", start, memoVal{ok: true, pos: p.pos, val: stmts})
	return stmts, true
}
