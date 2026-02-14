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

func parseBlockStmtV2(p *parser) ast.Stmt {
	p.expect(lexer.OPEN_CURLY)
	body := []ast.Stmt{}
	for p.hasTokens() {
		p.skipSeparators()
		if p.currentTokenKind() == lexer.CLOSE_CURLY {
			break
		}

		body = append(body, parseStmt(p))
	}

	p.expect(lexer.CLOSE_CURLY)
	return ast.BlockStmt{
		Body: body,
	}
}

func parsePublicStmt(p *parser) ast.Stmt {
	v := p.expect(lexer.PUBLIC).Literal

	return ast.PublicStmt{
		Visibility: v,
	}
}

func parseFunctionDeclStmt(p *parser) ast.Stmt {
	p.advance()
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
	p.expect(lexer.OPEN_PARENTHESIS)
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PARENTHESIS {
		paramName := p.expect(lexer.IDENTIFIER).Literal
		p.expect(lexer.COLON)
		fnParams = append(fnParams, ast.Parameter{
			Name: paramName,
		})

		if !p.currentToken().IsOneOfMany(lexer.CLOSE_PARENTHESIS, lexer.EOF) {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_PARENTHESIS)
	fnBody := ast.ExpectStmt[ast.BlockStmt](parseBlockStmtV2(p)).Body

	return fnParams, fnBody
}
