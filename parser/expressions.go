package parser

import (
	"experimental/ast"
	"experimental/lexer"
	"fmt"
	"strconv"
)

func parseExpr(p *parser, bp bindingPower) ast.Expr {
	tokenKind := p.currentTokenKind()
	nudFn, exists := nudLu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("NUD handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
	}

	left := nudFn(p)
	for bpLu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		ledFn, exists := ledLu[tokenKind]
		if !exists {
			panic(fmt.Sprintf("LED handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = ledFn(p, left, bpLu[p.currentTokenKind()])
	}

	return left
}

func parseAssignmentExpr(p *parser, left ast.Expr, bp bindingPower) ast.Expr {
	operatorToken := p.advance()
	rhs := parseExpr(p, bp)

	return ast.AssignmentExpr{
		Operator: operatorToken,
		Value:    rhs,
		Assignee: left,
	}
}

func parsePrimaryExpr(p *parser) ast.Expr {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Literal, 64)
		return ast.NumberExpr{
			Value: number,
		}
	case lexer.STRING:
		return ast.StringExpr{
			Value: p.advance().Literal,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpr{
			Value: p.advance().Literal,
		}
	default:
		panic(fmt.Sprintf("Cannot create primary expression from %s\n", lexer.TokenKindString(p.currentTokenKind())))
	}
}

func parseNumberExpr(p *parser) ast.Expr {
	value := p.advance().Literal
	n, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(err.Error())
	}

	return ast.NumberExpr{
		Value: n,
	}
}

func parseFunExpr(p *parser) ast.Expr {
	p.expect(lexer.FUN)
	funParams, funBody := parseFnParamsAndBody(p)
	return ast.FunctionExpr{
		Parameters: funParams,
		Body:       funBody,
	}
}

func parseCallExpr(p *parser, left ast.Expr, bp bindingPower) ast.Expr {
	p.advance()
	arguments := make([]ast.Expr, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PARENTHESIS {
		arguments = append(arguments, parseExpr(p, assignment))
		if !p.currentToken().IsOneOfMany(lexer.EOF, lexer.CLOSE_PARENTHESIS) {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_PARENTHESIS)
	return ast.CallExpr{
		Method:    left,
		Arguments: arguments,
	}

}

func parseMemberExpr(p *parser, left ast.Expr, bp bindingPower) ast.Expr {
	isComputed := p.advance().Kind == lexer.OPEN_BRACKET
	if isComputed {
		rhs := parseExpr(p, bp)
		p.expect(lexer.CLOSE_BRACKET)
		return ast.ComputedExpr{
			Member:   left,
			Property: rhs,
		}
	}

	return ast.MemberExpr{
		Member:   left,
		Property: p.expect(lexer.IDENTIFIER).Literal,
	}

}

func parseGroupingExpr(p *parser) ast.Expr {
	p.advance()
	expr := parseExpr(p, default_bp)
	p.expect(lexer.CLOSE_PARENTHESIS)
	return expr
}

func parseWildcardExpr(p *parser) ast.Expr {
	p.expect(lexer.UNDER_SCORE)
	return ast.WildcardExpr{}
}

// parses the match block
func parseMatchExpr(p *parser) ast.Expr {
	p.skipSeparators()
	p.expect(lexer.MATCH)
	expr := parseExpr(p, default_bp)
	arms := make([]ast.MatchArm, 0)
	for p.hasTokens() {
		p.skipSeparators()
		if p.currentTokenKind() != lexer.PIPE {
			break
		}

		p.advance() // consumes '|'
		pattern := parseExpr(p, default_bp)
		p.expect(lexer.THEN)
		body := parseExpr(p, default_bp)
		arms = append(arms, ast.MatchArm{
			Pattern: pattern,
			Body:    body,
		})
	}

	return ast.MatchExpr{
		Scrutinee: expr,
		Arms:      arms,
	}
}
