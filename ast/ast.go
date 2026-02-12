package ast

import (
	"fmt"
	"reflect"
)

type Stmt interface {
	stmt()
}

type Expr interface {
	expr()
}

func ExpectStmt[T Stmt](stmt Stmt) T {
	return expectType[T](stmt)
}

func ExpectExpr[T Expr](expr Expr) T {
	return expectType[T](expr)
}

func expectType[T any](r any) T {
	expectedType := reflect.TypeOf((*T)(nil)).Elem()
	receivedType := reflect.TypeOf(r)

	if expectedType == receivedType {
		return r.(T)
	}

	panic(fmt.Sprintf("Expected %T but received %T instead.", expectedType, receivedType))
}
