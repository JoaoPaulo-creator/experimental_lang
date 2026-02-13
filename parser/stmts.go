package parser

import (
	"experimental/ast"
	"experimental/lexer"
)

func parseStmt(p *parser) ast.Stmt {
	p.skipSeparators()
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

	isConstant := p.advance().Kind == lexer.LET
	varName := p.expectError(lexer.IDENTIFIER, "Inside variable declaration expected to find variable name").Literal

	if p.currentTokenKind() != lexer.SEMI_COLON {
		p.expect(lexer.FUNCTION_ASSOCIATION)
		assinedValue = parseExpr(p, function_association)
	}

	return ast.VarDeclStmt{
		IsConstant:    isConstant,
		VariableName:  varName,
		AssignedValue: assinedValue,
	}
}

func parseBlockStmt(p *parser) ast.Stmt {
	body := []ast.Stmt{}
	p.skipSeparators()
	for p.hasTokens() && p.currentTokenKind() != lexer.NEW_LINE {
		body = append(body, parseStmt(p))
		p.skipSeparators()
	}
	if p.hasTokens() && p.currentTokenKind() == lexer.NEW_LINE {
		p.advance()
	}

	return ast.BlockStmt{
		Body: body,
	}

}

func parseFunctionDeclStmt(p *parser) ast.Stmt {
	p.advance()
	fnName := p.expect(lexer.IDENTIFIER).Literal
	fnParameters, fnBody := parseFnParamsAndBody(p)

	return ast.FunctionDeclStmt{
		Name:       fnName,
		Parameters: fnParameters,
		Body:       fnBody,
	}
}

func parsePipeStmt(p *parser) ast.Stmt {
	p.advance()
	condition := parseExpr(p, assignment)
	p.expect(lexer.THEN)
	consequent := parseBlockStmt(p)

	return ast.PipeStmt{
		Condition:  condition,
		Consequent: consequent,
	}

}

func parseFnParamsAndBody(p *parser) ([]ast.Parameter, []ast.Stmt) {
	fnParams := make([]ast.Parameter, 0)
	for p.hasTokens() && p.currentTokenKind() != lexer.FUNCTION_ASSOCIATION {
		paramName := p.expect(lexer.IDENTIFIER).Literal
		fnParams = append(fnParams, ast.Parameter{
			Name: paramName,
		})
	}

	p.expect(lexer.FUNCTION_ASSOCIATION)
	fnBody := ast.ExpectStmt[ast.BlockStmt](parseBlockStmt(p)).Body

	return fnParams, fnBody
}
