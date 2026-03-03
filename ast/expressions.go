package ast

import "experimental/lexer"

type AssignmentExpr struct {
	Assignee Expr
	Operator lexer.Token
	Value    Expr
}

func (n AssignmentExpr) expr() {}

type PrefixExpr struct {
	Operator  lexer.Token
	RightExpr Expr
}

func (n PrefixExpr) expr() {}

type ComputedExpr struct {
	Member   Expr
	Property Expr
}

func (n ComputedExpr) expr() {}

type NumberExpr struct {
	Value float64
}

func (n NumberExpr) expr() {}

type StringExpr struct {
	Value string
}

func (n StringExpr) expr() {}

type SymbolExpr struct {
	Value string
}

func (n SymbolExpr) expr() {}

type MemberExpr struct {
	Member   Expr
	Property string
}

func (n MemberExpr) expr() {}

type WildcardExpr struct{}

func (n WildcardExpr) expr() {}

// ── declarations ─────────────────────────────────────────────────

type Program struct {
	Decls []LetDecl
}

type LetDecl struct {
	Name   string
	Params []string
	Body   Expr
}

// ── expressions ──────────────────────────────────────────────────

type MatchExpr struct {
	Subject Expr
	Arms    []MatchArm
}

type MatchArm struct {
	Pattern Expr // StrLit | Ident | Wildcard
	Body    Expr
}

type BinOp struct {
	Op    string
	Left  Expr
	Right Expr
}

type Ident struct{ Name string }
type StrLit struct{ Value string }
type NumLit struct{ Value string }
type Wildcard struct{}

func (Ident) expr()     {}
func (StrLit) expr()    {}
func (NumLit) expr()    {}
func (Wildcard) expr()  {}
func (BinOp) expr()     {}
func (MatchExpr) expr() {}
