package evaluator

import (
	"funlang/ast"
	"funlang/object"
)

func evalExpressions(exps []ast.ExpressionNode, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}
	return val
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("runtime error. unknown operator: %s for %s", operator, right.Inspect())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return NULL
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch operator {
	case "-", "*", "/":
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Integer).Value

		switch operator {
		case "-":
			return &object.Integer{Value: leftVal - rightVal}
		case "*":
			return &object.Integer{Value: leftVal * rightVal}
		case "/":
			return &object.Integer{Value: leftVal / rightVal}
		default:
			return newError("runtime error. operator %s for %s", operator, right.Inspect())
		}
	case "+":
		// TODO: поменять на ADD_INT и CONCAT_STR, добавить поле в node AST
		switch leftVal := left.(type) {
		case *object.Integer:
			return &object.Integer{Value: leftVal.Value + right.(*object.Integer).Value}
		case *object.String:
			return &object.String{Value: leftVal.Value + right.(*object.String).Value}
		}
	case ">", "<", ">=", "<=":
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Integer).Value
		switch operator {
		case ">":
			return nativeBoolToBooleanObject(leftVal > rightVal)
		case "<":
			return nativeBoolToBooleanObject(leftVal < rightVal)
		case ">=":
			return nativeBoolToBooleanObject(leftVal >= rightVal)
		case "<=":
			return nativeBoolToBooleanObject(leftVal <= rightVal)
		default:
			return newError("runtime error. operator %s for %s", operator, right.Inspect())
		}
	case "==", "!=":
		if left.Type() == object.BOOLEAN_OBJ {
			leftVal := left.(*object.Boolean).Value
			rightVal := right.(*object.Boolean).Value
			switch operator {
			case "==":
				return nativeBoolToBooleanObject(leftVal == rightVal)
			case "!=":
				return nativeBoolToBooleanObject(leftVal != rightVal)
			default:
				return newError("runtime error. operator %s for %s", operator, right.Inspect())
			}
		} else if left.Type() == object.INTEGER_OBJ {
			leftVal := left.(*object.Integer).Value
			rightVal := right.(*object.Integer).Value
			switch operator {
			case "==":
				return nativeBoolToBooleanObject(leftVal == rightVal)
			case "!=":
				return nativeBoolToBooleanObject(leftVal != rightVal)
			default:
				return newError("runtime error. operator %s for %s", operator, right.Inspect())
			}
		}
	case "&&", "||":
		leftVal := left.(*object.Boolean).Value
		rightVal := right.(*object.Boolean).Value
		switch operator {
		case "&&":
			return nativeBoolToBooleanObject(leftVal && rightVal)
		case "||":
			return nativeBoolToBooleanObject(leftVal || rightVal)
		default:
			return newError("runtime error. operator %s for %s", operator, right.Inspect())
		}
	}
	return newError("runtime error. operator %s for %s", operator, right.Inspect())
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if condition == TRUE {
		return Eval(node.Consequence, env)
	} else {
		return Eval(node.Alternative, env)
	}
}
