package parser

import (
	"experimental/ast"
	"experimental/lexer"
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
	tokens []lexer.Token
	pos    int
	memo   map[memoKey]memoVal
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
	return &PEGParser{tokens: tokens, memo: map[memoKey]memoVal{}}
}

func (p *PEGParser) cur() lexer.Token {
	return p.tokens[p.pos]
}

func (p *PEGParser) at(t lexer.Kind) bool {
	return p.cur().Kind == t
}

func (p *PEGParser) consume(t lexer.Kind) bool {
	if p.at(t) {
		p.pos++
		return true
	}

	return false
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

func (p *PEGParser) Program() (*ast.Program, bool) {
	start := p.pos
	if v, ok := p.memoized("Program", start); ok {
		if v.ok {
			return v.val.(*ast.Program), true
		}
		return nil, false
	}

	prog := &ast.Program{}
	for p.at(lexer.FUN) {
		fn, ok := p.FuncDecl()
		if !ok {
			p.pos = start
			p.memoize("Program", start, memoVal{ok: false, pos: start})
			return nil, false
		}

		prog.Func = append(prog.Func, fn)
	}

	if !p.at(lexer.EOF) {
		p.pos = start
		p.memoize("Program", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	p.memoize("Program", start, memoVal{ok: true, pos: p.pos, val: prog})
	return prog, true
}

func (p *PEGParser) FuncDecl() (*ast.FunctionDeclStmt, bool) {
	start := p.pos
	if v, ok := p.memoized("FunDecl", start); ok {
		if v.ok {
			return v.val.(*ast.FunctionDeclStmt), true
		}
		return nil, false
	}

	if !p.consume(lexer.FUN) {
		p.memoize("FunDecl", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	if !p.consume(lexer.IDENTIFIER) {
		p.pos = start
		p.memoize("FunDecl", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	name := p.cur().Literal
	p.pos++
	if !p.consume(lexer.OPEN_PARENTHESIS) && p.consume(lexer.CLOSE_PARENTHESIS) {
		p.pos = start
		p.memoize("FunDecl", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	stmts, ok := p.Block()
	if !ok {
		p.pos = start
		p.memoize("FunDecl", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	fn := &ast.FunctionDeclStmt{
		Name: name,
		Body: stmts,
	}

	p.memoize("FunDecl", start, memoVal{ok: true, pos: p.pos, val: fn})
	return fn, true
}

func (p *PEGParser) PrintLnStmt() (ast.Stmt, bool) {
	start := p.pos
	if v, ok := p.memoized("PrintLnStmt", start); ok {
		if v.ok {
			return v.val.(ast.Stmt), true
		}
	}

	if !p.consume(lexer.PRINT) || !p.consume(lexer.OPEN_PARENTHESIS) || !p.consume(lexer.STRING) {
		p.pos = start
		p.memoize("PrintLnStmt", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	val := p.cur().Literal
	p.pos++
	if !p.consume(lexer.CLOSE_PARENTHESIS) || !p.consume(lexer.SEMI_COLON) {
		p.pos = start
		p.memoize("PrintLnStmt", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	st := ast.PrintLnStmt{Value: val}
	p.memoize("PrintLnStmt", start, memoVal{ok: false, pos: p.pos, val: st})
	return st, true
}

func (p *PEGParser) Stmt() (ast.Stmt, bool) {
	start := p.pos
	if v, ok := p.memoized("Stmt", start); ok {
		if v.ok {
			return v.val.(ast.Stmt), true
		}
	}

	st, ok := p.PrintLnStmt()
	if !ok {
		p.memoize("Stmt", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	p.memoize("Stmt", start, memoVal{ok: true, pos: p.pos, val: st})
	return st, true
}

func (p *PEGParser) Block() ([]ast.Stmt, bool) {
	start := p.pos
	if v, ok := p.memoized("Block", start); ok {
		if v.ok {
			return v.val.([]ast.Stmt), true
		}
	}

	if !p.consume(lexer.OPEN_CURLY) {
		p.memoize("Block", start, memoVal{ok: false, pos: start})
	}

	var stmts []ast.Stmt
	for p.at(lexer.PRINT) {
		st, ok := p.Stmt()
		if !ok {
			p.pos = start
			p.memoize("Block", start, memoVal{ok: false, pos: start})
			return nil, false
		}

		stmts = append(stmts, st)
	}

	if !p.consume(lexer.CLOSE_CURLY) {
		p.pos = start
		p.memoize("Block", start, memoVal{ok: false, pos: start})
		return nil, false
	}

	p.memoize("Block", start, memoVal{ok: true, pos: p.pos, val: stmts})
	return stmts, true
}
