package evaluator

import (
	"funlang/ast"
	"funlang/object"
	"funlang/type_checker"
)

func evalIfExpression(node *ast.IfExpression, env *object.Environment, chk *type_checker.TypeChecker) object.Object {
	condition := Eval(node.Condition, env, chk)
	var res object.Object
	if condition == TRUE {
		res = Eval(node.Consequence, env, chk)
	} else {
		res = Eval(node.Alternative, env, chk)
	}
	return unwrapReturnValue(res)
}

func evalSwitchExpression(node *ast.SwitchExpression, env *object.Environment, chk *type_checker.TypeChecker) object.Object {
	condition := Eval(node.Value, env, chk)
	var res object.Object
	for _, caseBlock := range node.Cases {
		pattern := Eval(caseBlock.Pattern, env, chk)
		if object.Equals(pattern, condition) {
			res = Eval(caseBlock.Body, env, chk)
			return unwrapReturnValue(res)
		}
	}

	res = Eval(node.Default, env, chk)

	return unwrapReturnValue(res)
}
