package evaluator

import (
	"funlang/ast"
	"funlang/object"
	"funlang/type_checker"
	"funlang/types"
)

func evalExpressions(exps []ast.ExpressionNode, env *object.Environment,
	chk *type_checker.TypeChecker) []object.Object {

	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env, chk)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_MAP_OBJ:
		return evalHashIndexExpression(left, index)
	}
	return newError("runtime error. index operator not supported: %s", left.Inspect())
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := len(arrayObject.Elements) - 1
	if idx < 0 || idx > max {
		// return NULL
		return newError("runtime error: array index %d is out of bounds (max is %d)",
			idx, len(arrayObject.Elements)-1)
	}
	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.HashMap)

	key := index.(object.Hashable)

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
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

func evalInfixExpression(operator string, left, right object.Object) (res object.Object) {
	defer func() {
		if r := recover(); r != nil {
			res = newError("%s", r)
		}
	}()

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
		// TODO: поменять на ADD_INT и CONCAT_STR, добавить поле с инструкцией в node AST
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
	case "==":
		if _, isComparable := left.(object.Comparable); left.Type() == right.Type() && isComparable {
			return nativeBoolToBooleanObject(object.Equals(left, right))
		}
	case "!=":
		if _, isComparable := left.(object.Comparable); left.Type() == right.Type() && isComparable {
			return nativeBoolToBooleanObject(!object.Equals(left, right))
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
	case "&", "|":
		leftVal := left.(*object.TypeObject).TypeValue
		rightVal := right.(*object.TypeObject).TypeValue
		switch operator {
		case "&":
			return &object.TypeObject{TypeValue: &types.IntersectionType{Left: leftVal, Right: rightVal}}
		case "|":
			return &object.TypeObject{TypeValue: &types.UnionType{Left: leftVal, Right: rightVal}}
		default:
			return newError("runtime error. operator %s for %s", operator, right.Inspect())
		}
	}
	return newError("runtime error. operator %s for %s", operator, right.Inspect())
}
