package evaluator

import (
	"funlang/ast"
	"funlang/object"
)

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
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
	case "-", "+", "*", "/":
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Integer).Value

		switch operator {
		case "+":
			return &object.Integer{Value: leftVal + rightVal}
		case "-":
			return &object.Integer{Value: leftVal - rightVal}
		case "*":
			return &object.Integer{Value: leftVal * rightVal}
		case "/":
			return &object.Integer{Value: leftVal / rightVal}
		default:
			return NULL
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
			return NULL
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
				return NULL
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
				return NULL
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
			return NULL
		}
	}
	return NULL
}

func evalIfExpression(node *ast.IfExpression) object.Object {
	condition := Eval(node.Condition)
	if condition == TRUE {
		return Eval(node.Consequence)
	} else {
		return Eval(node.Alternative)
	}
}
