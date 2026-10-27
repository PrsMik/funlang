package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/evaluator"
	"funlang/object"
	"funlang/types"
	"reflect"
)

func (chk *TypeChecker) checkExpression(expr ast.ExpressionNode) types.Type {
	exprType := reflect.TypeOf(expr)
	if chk.curExpectedType != nil {
		chk.recordExpectedType(expr, chk.curExpectedType)
	}
	if checkFun, ok := chk.expressionCheckFns[exprType]; ok {
		tp := checkFun(expr)
		chk.recordType(expr, tp)
		return tp
	}
	chk.typeError(fmt.Sprintf("unknown expression type %T", expr), expr)
	return &types.IllegalType{}
}

func (chk *TypeChecker) checkBadExpression(expr ast.ExpressionNode) types.Type {
	if chk.curExpectedType != nil {
		chk.recordExpectedType(expr, chk.curExpectedType)
		return chk.curExpectedType
	}
	return &types.IllegalType{}
}

func (chk *TypeChecker) checkPrefixExpression(expr ast.ExpressionNode) types.Type {
	op := expr.(*ast.PrefixExpression).Operator

	oldType := chk.curExpectedType

	switch op {
	case "-":
		chk.curExpectedType = &types.IntType{}
	case "!":
		chk.curExpectedType = &types.BoolType{}
	}

	rightType := chk.checkExpression(expr.(*ast.PrefixExpression).Right)
	chk.curExpectedType = oldType

	switch op {
	case "-":
		if types.Equals(rightType, &types.IntType{}) {
			return &types.IntType{}
		}
	case "!":
		if types.Equals(rightType, &types.BoolType{}) {
			return &types.BoolType{}
		}
	}
	chk.typeError(fmt.Sprintf("unknown operator: %s for type %s", op, rightType.Signature()), expr)
	return &types.IllegalType{}
}

func (chk *TypeChecker) checkInfixExpression(expr ast.ExpressionNode) types.Type {
	leftType := chk.checkExpression(expr.(*ast.InfixExpression).Left)
	op := expr.(*ast.InfixExpression).Operator

	oldType := chk.curExpectedType

	switch op {
	case "-", "*", "/", ">", "<", ">=", "<=":
		chk.curExpectedType = &types.IntType{}
	case "&&", "||":
		chk.curExpectedType = &types.BoolType{}
	case "+":
		if types.Equals(leftType, &types.StringType{}) {
			chk.curExpectedType = &types.StringType{}
		} else {
			chk.curExpectedType = &types.IntType{}
		}
	case "==", "!=":
		chk.curExpectedType = leftType
	case "|", "&":
		chk.curExpectedType = &types.TypeType{}
	}

	rightType := chk.checkExpression(expr.(*ast.InfixExpression).Right)
	chk.curExpectedType = oldType

	switch op {
	case "-", "+", "*", "/":
		if types.Equals(leftType, &types.IntType{}) && types.Equals(rightType, &types.IntType{}) {
			return &types.IntType{}
		} else if op == "+" &&
			types.Equals(leftType, &types.StringType{}) &&
			types.Equals(rightType, &types.StringType{}) {
			return &types.StringType{}
		}
	case "&&", "||":
		if types.Equals(leftType, &types.BoolType{}) && types.Equals(rightType, &types.BoolType{}) {
			return &types.BoolType{}
		}
	case "==", "!=":
		if _, isComparable := leftType.(types.ComparableType); types.Equals(leftType, rightType) && isComparable {
			return &types.BoolType{}
		}
	case ">", "<", ">=", "<=":
		if types.Equals(leftType, &types.IntType{}) && types.Equals(rightType, &types.IntType{}) {
			return &types.BoolType{}
		}
	case "|":
		if types.Equals(types.MetaType, leftType) && types.Equals(types.MetaType, rightType) {
			return &types.TypeType{
				Underlying: &types.UnionType{
					Left:  leftType.(*types.TypeType).Underlying,
					Right: rightType.(*types.TypeType).Underlying,
				},
			}
		}
	case "&":
		if types.Equals(types.MetaType, leftType) && types.Equals(types.MetaType, rightType) {
			return &types.TypeType{
				Underlying: &types.IntersectionType{
					Left:  leftType.(*types.TypeType).Underlying,
					Right: rightType.(*types.TypeType).Underlying,
				},
			}
		}
	}

	chk.typeError(fmt.Sprintf("type mismatch between: %s & %s; for operator %s",
		leftType.Signature(), rightType.Signature(), op), expr)
	return &types.IllegalType{}
}

func (chk *TypeChecker) checkComptimeExpression(expr ast.ExpressionNode) types.Type {
	comptimeExpr := expr.(*ast.ComptimeExpression)

	valType := chk.checkExpression(comptimeExpr.Value)
	if _, isIllegal := valType.(*types.IllegalType); isIllegal {
		return valType
	}

	val := evaluator.Eval(comptimeExpr.Value, chk.comptimeEnv, chk.Info)
	if val == nil || val.Type() == object.ERROR_OBJ {
		errMsg := "unknown comptime error"
		if errObj, ok := val.(*object.Error); ok {
			errMsg = errObj.Message
		}
		chk.typeError(fmt.Sprintf("comptime evaluation error: %s", errMsg), expr)
		return &types.IllegalType{}
	}

	chk.ComptimeValues[expr] = val

	if typeObj, ok := val.(*object.TypeObject); ok {
		return &types.TypeType{Underlying: typeObj.TypeValue}
	}

	evaluatedType := object.GetTypeFromObject(val)
	if _, isIllegal := evaluatedType.(*types.IllegalType); !isIllegal {
		chk.recordType(expr, evaluatedType)
		return evaluatedType
	}

	return valType
}
