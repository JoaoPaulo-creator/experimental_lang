package parser

import (
	"experimental/ast"
	"experimental/lexer"
)

type bindingPower int

const (
	default_bp bindingPower = iota
	function_association
	primary
	unary
	member
)

type (
	stmtHandler func(p *parser) ast.Stmt
	nudHandler  func(p *parser) ast.Expr
	ledHandler  func(p *parser, left ast.Expr, bp bindingPower) ast.Expr

	stmtLookup map[lexer.Kind]stmtHandler
	nudLookup  map[lexer.Kind]nudHandler
	ledLookup  map[lexer.Kind]ledHandler
	bpLookup   map[lexer.Kind]bindingPower
)

var (
	bpLu   = bpLookup{}
	nudLu  = nudLookup{}
	ledLu  = ledLookup{}
	stmtLu = stmtLookup{}
)

func led(kind lexer.Kind, bp bindingPower, ledFn ledHandler) {
	bpLu[kind] = bp
	ledLu[kind] = ledFn
}

func nud(kind lexer.Kind, nudFn nudHandler) {
	nudLu[kind] = nudFn
}

func stmt(kind lexer.Kind, stmtFn stmtHandler) {
	bpLu[kind] = default_bp
	stmtLu[kind] = stmtFn
}

func createTokenLookups() {

	nud(lexer.STRING, parsePrimaryExpr)
	nud(lexer.IDENTIFIER, parsePrimaryExpr)

	stmt(lexer.LET, parseFunctionDeclStmt)
}
