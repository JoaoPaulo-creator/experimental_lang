package parser

import (
	"experimental/ast"
	"experimental/lexer"
)

func parseStmt(p *parser) ast.Stmt {
	stmtFn, exists := stmtLu[p.currentTokenKind()]
	if exists {
		return stmtFn(p)
	}

	expression := parseExpr(p, default_bp)
	return ast.ExpressionStmt{
		Expression: expression,
	}
}

func parseVarDeclStmt(p *parser) ast.Stmt {
	var assinedValue ast.Expr

	isConstant := p.advance().Kind == lexer.GIVEN
	varName := p.expectError(lexer.IDENTIFIER, "Inside variable declaration expected to find variable name").Literal

	// if p.currentTokenKind() == lexer.COLON {
	// 	p.advance() // eat the colon
	// }

	if p.currentTokenKind() != lexer.SEMI_COLON {
		p.expect(lexer.ASSIGNMENT)
		assinedValue = parseExpr(p, assignment)
	}

	p.expect(lexer.SEMI_COLON)

	// if isConstant && assinedValue == nil {
	// 	panic("Cannot define constant without providing a value!")
	// }

	return ast.VarDeclStmt{
		IsConstant:    isConstant,
		VariableName:  varName,
		AssignedValue: assinedValue,
	}
}
