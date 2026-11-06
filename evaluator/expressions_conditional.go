package evaluator

import (
	"funlang/ast"
	"funlang/object"
	"funlang/types"
)

func evalIfExpression(node *ast.IfExpression, env *object.Environment, info *types.Info) object.Object {
	return forceTailCall(evalIfExpressionTail(node, env, info), info)
}

func evalIfExpressionTail(node *ast.IfExpression, env *object.Environment, info *types.Info) object.Object {
	condition := Eval(node.Condition, env, info)
	if isError(condition) {
		return condition
	}

	var res object.Object
	if condition == TRUE {
		res = Eval(node.Consequence, object.NewEnclosedEnvironment(env), info)
	} else if node.Alternative != nil {
		res = Eval(node.Alternative, object.NewEnclosedEnvironment(env), info)
	}
	return unwrapReturnValue(res)
}

func evalSwitchExpression(node *ast.SwitchExpression, env *object.Environment, info *types.Info) object.Object {
	return forceTailCall(evalSwitchExpressionTail(node, env, info), info)
}

func evalSwitchExpressionTail(node *ast.SwitchExpression, env *object.Environment, info *types.Info) object.Object {
	condition := Eval(node.Value, env, info)
	if isError(condition) {
		return condition
	}

	for _, caseBlock := range node.Cases {
		pattern := Eval(caseBlock.Pattern, env, info)
		if isError(pattern) {
			return pattern
		}

		if object.Equals(pattern, condition) {
			res := Eval(caseBlock.Body, object.NewEnclosedEnvironment(env), info)
			return unwrapReturnValue(res)
		}
	}

	res := Eval(node.Default, object.NewEnclosedEnvironment(env), info)
	return unwrapReturnValue(res)
}

func forceTailCall(obj object.Object, info *types.Info) object.Object {
	if tc, ok := obj.(*object.TailCall); ok {
		return applyFunction(tc.Function, tc.Arguments, info)
	}
	return obj
}
