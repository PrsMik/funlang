package evaluator

import (
	"funlang/ast"
	"funlang/object"
	"funlang/types"
)

func evalIfExpression(node *ast.IfExpression, env *object.Environment, info *types.Info) object.Object {
	condition := Eval(node.Condition, env, info)
	if isError(condition) {
		return condition
	}

	var res object.Object
	if condition == TRUE {
		res = Eval(node.Consequence, env, info)
	} else {
		res = Eval(node.Alternative, env, info)
	}
	return unwrapReturnValue(res)
}

func evalSwitchExpression(node *ast.SwitchExpression, env *object.Environment, info *types.Info) object.Object {
	condition := Eval(node.Value, env, info)
	if isError(condition) {
		return condition
	}

	var res object.Object
	for _, caseBlock := range node.Cases {
		pattern := Eval(caseBlock.Pattern, env, info)
		if object.Equals(pattern, condition) {
			res = Eval(caseBlock.Body, env, info)
			return unwrapReturnValue(res)
		}
	}

	res = Eval(node.Default, env, info)

	return unwrapReturnValue(res)
}
