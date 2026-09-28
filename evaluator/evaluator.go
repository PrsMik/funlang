package evaluator

import (
	"fmt"
	"funlang/ast"
	"funlang/object"
	"funlang/types"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment, info *types.Info) object.Object {
	switch node := node.(type) {
	// инструкции
	case *ast.Program:
		return evalProgram(node, env, info)
	case *ast.LetStatement:
		val := Eval(node.Value, env, info)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val
	case *ast.ReturnStatement:
		if callNode, ok := node.Value.(*ast.CallExpression); ok {
			function := Eval(callNode.Function, env, info)
			if isError(function) {
				return function
			}

			args := evalExpressions(callNode.Arguments, env, info)
			if len(args) == 1 && isError(args[0]) {
				return args[0]
			}

			return &object.ReturnValue{Value: &object.TailCall{Function: function, Arguments: args}}
		}

		val := Eval(node.Value, env, info)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	case *ast.BlockStatement:
		return evalBlockStatement(node, env, info)

	// ВЫРАЖЕНИЯ
	// литералы
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env, info)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{
			Elements: elements,
			TypeObj:  &object.TypeObject{TypeValue: info.TypesInfo[node]}}

	case *ast.HashMapLiteral:
		return evalHashMapLiteral(node, env, info)
	// лиетралы типов
	case *ast.InterfaceLiteral:
		return evalInterfaceLiteral(node, env, info)
	case *ast.ImplLiteral:
		return evalImplLiteral(node, env, info)
	// выражения
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.TypeType:
		return evalTypeTypeIdentifier(node, env)
	case *ast.SimpleType:
		return evalTypeIdentifier(node, env)
	case *ast.IndexExpression:
		left := Eval(node.Left, env, info)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env, info)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.FunctionLiteral:
		resolvedType := info.TypesInfo[node]

		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
			TypeObj:    &object.TypeObject{TypeValue: resolvedType},
		}

	case *ast.CallExpression:
		function := Eval(node.Function, env, info)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env, info)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args, info)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env, info)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env, info)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env, info)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env, info)
	case *ast.SwitchExpression:
		return evalSwitchExpression(node, env, info)
	case *ast.MemberAccessExpression:
		return evalMemberAccessExpression(node, env, info)
	case *ast.TypeAccessExpression:
		return evalTypeAccessExpression(node, env, info)
	case *ast.ImplInstantiationExpression:
		return evalImplInstantiationExpression(node, env, info)
	}
	return nil
}

func evalProgram(program *ast.Program, env *object.Environment, info *types.Info) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env, info)
		switch result := result.(type) {
		case *object.ReturnValue:
			if tc, ok := result.Value.(*object.TailCall); ok {
				return applyFunction(tc.Function, tc.Arguments, info)
			}
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment,
	info *types.Info) object.Object {

	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env, info)
		if result != nil {
			resType := result.Type()
			if resType == object.RETURN_VALUE_OBJ || resType == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
