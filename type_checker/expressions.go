package type_checker

import (
	"fmt"
	"funlang/ast"
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
		if (types.Equals(leftType, &types.BoolType{}) && types.Equals(rightType, &types.BoolType{})) ||
			(types.Equals(leftType, &types.IntType{}) && types.Equals(rightType, &types.IntType{})) {
			return &types.BoolType{}
		}
	case ">", "<", ">=", "<=":
		if types.Equals(leftType, &types.IntType{}) && types.Equals(rightType, &types.IntType{}) {
			return &types.BoolType{}
		}
	case "|":
		if types.Equals(types.TrueTypeType, leftType) && types.Equals(types.TrueTypeType, rightType) {
			return &types.TypeType{
				Underlying: &types.UnionType{
					Left:  leftType.(*types.TypeType).Underlying,
					Right: rightType.(*types.TypeType).Underlying,
				},
			}
		}
	case "&":
		if types.Equals(types.TrueTypeType, leftType) && types.Equals(types.TrueTypeType, rightType) {
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

func (chk *TypeChecker) checkCallExpression(expr ast.ExpressionNode) types.Type {
	callExpr := expr.(*ast.CallExpression)
	rawCallType := chk.checkExpression(callExpr.Function)

	if _, err := rawCallType.(*types.IllegalType); err {
		return &types.IllegalType{}
	}

	switch callFuncType := rawCallType.(type) {
	case *types.FuncType:
		if len(callExpr.Arguments) != len(callFuncType.Params) {
			if len(callFuncType.Params) != 0 {
				chk.typeError(fmt.Sprintf(`wrong number of arguments for "%s"`, callExpr.Function.String()), expr)
				return &types.IllegalType{}
			} else {
				if _, ok := callExpr.Arguments[0].(*ast.BadExpression); ok {
					callExpr.Arguments = callExpr.Arguments[1:]
				}
			}
		}

		for i, arg := range callExpr.Arguments {
			argType := chk.checkExpression(arg)

			if !types.Equals(argType, callFuncType.Params[i].Type) {
				chk.typeError(fmt.Sprintf("wrong type for argument %d: %s in func call \" %s \"; expected %s",
					i+1, argType.Signature(), callExpr.Function.String(), callFuncType.Params[i].Type.Signature()), expr)
				return &types.IllegalType{}
			}
		}

		return callFuncType.ReturnType
	case *types.BuiltinFunc:
		argTypes := make([]types.Type, len(callExpr.Arguments))
		for i, argExpr := range callExpr.Arguments {
			argTypes[i] = chk.checkExpression(argExpr)
		}

		returnType, err := callFuncType.CheckFunc(argTypes)

		if err != nil {
			// chk.errors = append(chk.errors, err)
			chk.typeError(err.Error(), expr)
			return &types.IllegalType{}
		}

		return returnType
	default:
		chk.typeError(fmt.Sprintf(`identifier "%s" is no a funciton`, callExpr.Function.String()), expr)
		return &types.IllegalType{}
	}

}
