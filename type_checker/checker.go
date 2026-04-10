package type_checker

import (
	"funlang/ast"
	"funlang/types"
	"reflect"
)

type expressionCheckFn func(ast.ExpressionNode) types.Type

func (chk *TypeChecker) registerExpressionCheckFn(exprType reflect.Type, fn expressionCheckFn) {
	chk.expressionCheckFns[exprType] = fn
}

type TypeError struct {
	Msg  string
	Node ast.Node
}

type TypeChecker struct {
	env                *types.TypeEviroment
	expressionCheckFns map[reflect.Type]expressionCheckFn
	curExpectedType    types.Type
	errors             []TypeError

	// для сервера
	TypesInfo     map[ast.Node]types.Type
	Definitions   map[ast.Node]ast.Node // лучше хранить по самому идентификатору, но так проще перебирать в сервере
	Scopes        map[ast.Node]*types.TypeEviroment
	ExpectedTypes map[ast.Node]types.Type
}

func New(curEnv *types.TypeEviroment) *TypeChecker {
	chk := &TypeChecker{env: curEnv,
		TypesInfo:     make(map[ast.Node]types.Type),
		Definitions:   make(map[ast.Node]ast.Node),
		Scopes:        make(map[ast.Node]*types.TypeEviroment),
		ExpectedTypes: make(map[ast.Node]types.Type)}

	chk.expressionCheckFns = make(map[reflect.Type]expressionCheckFn)

	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.VirtualNode](), chk.checkUnparsedNode)

	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.IntegerLiteral](), chk.checkIntegerLiteral)
	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.BooleanLiteral](), chk.checkBooleanLiteral)
	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.StringLiteral](), chk.checkStringLiteral)
	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.ArrayLiteral](), chk.checkArrayLiteral)
	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.HashMapLiteral](), chk.checkHashMapLiteral)

	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.Identifier](), chk.checkIdentifier)

	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.PrefixExpression](), chk.checkPrefixExpression)
	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.InfixExpression](), chk.checkInfixExpression)

	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.IfExpression](), chk.checkIfExpression)

	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.FunctionLiteral](), chk.checkFunctionLiteral)
	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.CallExpression](), chk.checkCallExpression)
	chk.registerExpressionCheckFn(reflect.TypeFor[*ast.IndexExpression](), chk.checkIndexExpression)

	return chk
}

func (chk *TypeChecker) CheckProgram(prog *ast.Program) {
	for _, stmt := range prog.Statements {
		chk.checkStatement(stmt)
	}
}

func (chk *TypeChecker) Errors() []TypeError {
	return chk.errors
}

func (chk *TypeChecker) GetEnv() types.TypeEviroment {
	return *chk.env
}
