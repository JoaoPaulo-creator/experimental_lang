package ast

type Stmt interface {
	stmt()
}

type Expr interface {
	expr()
}

// func ExpectStmt[T Stmt](stmt Stmt) T {
// 	return helpers.ExpectType[T](stmt)
// }

// func ExpectExpr[T Expr](expr Expr) T {
// 	return helpers.ExpectType[T](expr)
// }
