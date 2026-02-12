package ast

type BlockStmt struct {
	Body []Stmt
}

func (b BlockStmt) stmt() {}

type MethodStmt struct {
	Value Expr // or string if you want to simplify
}

func (m MethodStmt) stmt() {}

type EndpointStmt struct {
	Value Expr // or string if you want to simplify
}

func (m EndpointStmt) stmt() {}

// headers section
type HeadersStmt struct {
	Properties map[string]string
}

func (m HeadersStmt) stmt() {}

// body section
type BodyStmt struct {
	Properties map[string]string
}

func (m BodyStmt) stmt() {}

type ExpressionStmt struct {
	Expression Expr
}

func (n ExpressionStmt) stmt() {}

type VarDeclStmt struct {
	VariableName  string
	IsConstant    bool
	AssignedValue Expr
}

func (n VarDeclStmt) stmt() {}

type Parameter struct {
	Name string
}

type FunctionDeclStmt struct {
	Name       string
	Parameters []Parameter
	Body       []Stmt
}

func (n FunctionDeclStmt) stmt() {}
