package types

import (
	"funlang/ast"
	"reflect"
)

type expressionCheckFn func(ast.ExpressionNode) Type

func (chk *TypeChecker) registerExpressionCheckFn(exprType reflect.Type, fn expressionCheckFn) {
	chk.expressionCheckFns[exprType] = fn
}

type TypeChecker struct {
	env                *TypeEviroment
	expressionCheckFns map[reflect.Type]expressionCheckFn
	curExpectedType    Type
	errors             []error
}

func New(curEnv *TypeEviroment) *TypeChecker {
	chk := &TypeChecker{env: curEnv}

	chk.expressionCheckFns = make(map[reflect.Type]expressionCheckFn)

	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.IntegerLiteral](), chk.checkIntegerLiteral)
	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.BooleanLiteral](), chk.checkBooleanLiteral)
	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.Identifier](), chk.checkIdentifier)

	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.PrefixExpression](), chk.checkPrefixExpression)
	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.InfixExpression](), chk.checkInfixExpression)

	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.IfExpression](), chk.checkIfExpression)

	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.FunctionLiteral](), chk.checkFunctionLiteral)

	return chk
}

func (chk *TypeChecker) CheckProgram(prog *ast.Program) {
	chk.env = NewEnclosedTypeEviroment(chk.env)
	for _, stmt := range prog.Statements {
		chk.checkStatement(stmt)
	}
}

func (chk *TypeChecker) Errors() []error {
	return chk.errors
}
