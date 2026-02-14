package parser

import (
	"experimental/ast"
	"experimental/lexer"
)

type bindingPower int

const (
	default_bp bindingPower = iota
	function_association
	assignment
	primary
	aditive
	call
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
	// led(lexer.ASSIGNMENT, assignment, parseAssignmentExpr)
	// led(lexer.PLUS_EQUALS, assignment, parseAssignmentExpr)
	// led(lexer.MINUS_EQUALS, assignment, parseAssignmentExpr)

	// Logical
	// led(lexer.AND, logical, parseBinaryExpr)
	// led(lexer.OR, logical, parseBinaryExpr)
	// led(lexer.DOT_DOT, logical, parseBinaryExpr)

	// Relational
	// led(lexer.LESS, relational, parseBinaryExpr)
	// led(lexer.LESS_EQUALS, relational, parseBinaryExpr)
	// led(lexer.GREATER, relational, parseBinaryExpr)
	// led(lexer.GREATER_EQUALS, relational, parseBinaryExpr)
	// led(lexer.EQUALS, relational, parseBinaryExpr)
	// led(lexer.NOT_EQUALS, relational, parseBinaryExpr)

	// Additive & Multiplicative
	// led(lexer.PLUS, additive, parseBinaryExpr)
	// led(lexer.DASH, additive, parseBinaryExpr)
	// led(lexer.SLASH, multiplicative, parseBinaryExpr)
	// led(lexer.STAR, multiplicative, parseBinaryExpr)
	// led(lexer.PERCENT, multiplicative, parseBinaryExpr)

	// Literals & Symbols
	nud(lexer.NUMBER, parsePrimaryExpr)
	nud(lexer.STRING, parsePrimaryExpr)
	nud(lexer.IDENTIFIER, parsePrimaryExpr)

	//Unary/Prefix
	// nud(lexer.TYPEOF, parsePrimaryExpr)
	// nud(lexer.DASH, parsePrefixExpr)
	// nud(lexer.NOT, parsePrimaryExpr)
	nud(lexer.OPEN_BRACKET, parsePrimaryExpr)

	// Member / Computed // Call
	// led(lexer.DOT, member, parseMemberExpr)
	led(lexer.OPEN_BRACKET, member, parseMemberExpr)
	// led(lexer.OPEN_PAREN, member, parseMemberExpr)

	// Grouping Expr
	// nud(lexer.OPEN_PAREN, parseGroupingExpr)
	// nud(lexer.FN, parseGroupingExpr)
	// nud(lexer.NEW, func(p *parser) ast.Expr {
	// 	p.advance()
	// 	classInstantiation := parseExpr(p, default_bp)

	// 	return ast.NewExpr{
	// 		Instantiation: ast.ExpectExpr[ast.CallExpr](classInstantiation),
	// 	}
	// })

	// nud(lexer.TYPEOF, parsePrefixExpr)
	// nud(lexer.DASH, parsePrefixExpr)
	// nud(lexer.NOT, parsePrefixExpr)
	// nud(lexer.OPEN_BRACKET, parseArrayInstantiationExpr)

	// Call/Member/Arrays expressions
	// led(lexer.DOT, member, parseMemberExpr)
	led(lexer.OPEN_BRACKET, member, parseMemberExpr)
	led(lexer.OPEN_PARENTHESIS, call, parseCallExpr)

	// Grouping Expr
	nud(lexer.OPEN_PARENTHESIS, parseGroupingExpr)
	nud(lexer.FUN, parseFunExpr)
	// nud(lexer.NEW, func(p *parser) ast.Expr {
	// 	p.advance()
	// 	classInstation := parseExpr(p, default_bp)
	// 	return ast.NewExpr{
	// 		Instantiation: ast.ExpectExpr[ast.CallExpr](classInstation),
	// 	}
	// })

	// Statements
	stmt(lexer.OPEN_CURLY, parseBlockStmtV2)
	stmt(lexer.LET, parseVarDeclStmt)
	// stmt(lexer.CONST, parseVarDeclStmt)
	stmt(lexer.FUN, parseFunctionDeclStmt)
	// stmt(lexer.IF, parseIfStmt)
	// stmt(lexer.IMPORT, parseImportStmt)
	// stmt(lexer.FOREACH, parseForEarchStmt)
	// stmt(lexer.CLASS, parseClassDeclStmt)
	// stmt(lexer.STRUCT, parseStructDeclStmt)
}
