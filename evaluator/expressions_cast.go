package evaluator

import (
	"funlang/ast"
	"funlang/object"
	"funlang/types"
)

func evalTypeCastExpression(node *ast.SafeTypeCastExpression, env *object.Environment, info *types.Info) object.Object {
	left := Eval(node.Left, env, info)
	if isError(left) {
		return left
	}

	targetVal := Eval(node.TargetType, env, info)
	if isError(targetVal) {
		return targetVal
	}

	tObj, ok := targetVal.(*object.TypeObject)
	if !ok {
		return newError("runtime error: cast target is not a type")
	}

	if types.IsAssignable(tObj.TypeValue, object.GetTypeFromObject(left)) {
		return left
	}

	return NIL
}

func evalUnsafeTypeCastExpression(node *ast.UnsafeTypeCastExpression, env *object.Environment, info *types.Info) object.Object {
	targetVal := Eval(node.TargetType, env, info)
	if isError(targetVal) {
		return targetVal
	}

	tObj, ok := targetVal.(*object.TypeObject)
	if !ok {
		return newError("runtime error: cast target is not a type")
	}

	val := Eval(node.Value, env, info)
	if isError(val) {
		return val
	}

	actualType := object.GetTypeFromObject(val)
	if !types.IsAssignable(tObj.TypeValue, actualType) {
		return newError("runtime error: cannot unsafe cast %s to %s",
			actualType.Signature(), tObj.TypeValue.Signature())
	}

	return val
}
