package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/types"
	"reflect"
)

func (chk *TypeChecker) checkExpression(expr ast.ExpressionNode) types.Type {
	exprTp := reflect.TypeOf(expr)
	if fn, ok := chk.expressionCheckFns[exprTp]; ok {
		return fn(expr)
	}
	chk.errors = append(chk.errors, fmt.Errorf("unknown expression type"))
	return &types.IllegalType{}
}

func (chk *TypeChecker) checkIntegerLiteral(expr ast.ExpressionNode) types.Type {
	return &types.IntType{}
}

func (chk *TypeChecker) checkBooleanLiteral(expr ast.ExpressionNode) types.Type {
	return &types.BoolType{}
}

func (chk *TypeChecker) checkStringLiteral(expr ast.ExpressionNode) types.Type {
	return &types.StringType{}
}

func (chk *TypeChecker) checkArrayLiteral(expr ast.ExpressionNode) types.Type {
	arrType := &types.ArrayType{}
	arrLit := expr.(*ast.ArrayLiteral)

	if len(arrLit.Elements) == 0 {
		return arrType
	}

	firstArrLitType := chk.checkExpression(arrLit.Elements[0])

	for _, param := range arrLit.Elements[1:] {
		curParamType := chk.checkExpression(param)
		if !types.Equals(firstArrLitType, curParamType) {
			chk.errors = append(chk.errors,
				fmt.Errorf("array literal has elements of different types %s & %s",
					firstArrLitType.Signature(), curParamType.Signature()))
			return &types.IllegalType{}
		}
	}

	arrType.ElementsType = firstArrLitType

	return arrType
}

func (chk *TypeChecker) checkHashMapLiteral(expr ast.ExpressionNode) types.Type {
	hashMapType := &types.HashMapType{}
	hashMapLiteral := expr.(*ast.HashMapLiteral)

	if len(hashMapLiteral.Pairs) == 0 {
		return hashMapType
	}

	var firstHashMapKeyType types.Type = nil
	var firstHashMapElementType types.Type = nil

	for key, elem := range hashMapLiteral.Pairs {
		curHashMapKeyType := chk.checkExpression(key)
		curHashMapElementType := chk.checkExpression(elem)

		if firstHashMapKeyType == nil && firstHashMapElementType == nil {
			firstHashMapKeyType = curHashMapKeyType
			firstHashMapElementType = curHashMapElementType
		} else {
			if !types.Equals(firstHashMapKeyType, curHashMapKeyType) {
				chk.errors = append(chk.errors,
					fmt.Errorf("map literal has keys of different types %s & %s",
						firstHashMapKeyType.Signature(), curHashMapKeyType.Signature()))
				return &types.IllegalType{}
			}

			if !types.Equals(firstHashMapElementType, curHashMapElementType) {
				chk.errors = append(chk.errors,
					fmt.Errorf("map literal has elements of different types %s & %s",
						firstHashMapElementType.Signature(), curHashMapElementType.Signature()))
				return &types.IllegalType{}
			}
		}
	}

	hashMapType.KeyType = firstHashMapKeyType
	hashMapType.ElementType = firstHashMapElementType

	return hashMapType
}

func (chk *TypeChecker) checkIndexExpression(expr ast.ExpressionNode) types.Type {
	indexType := chk.checkExpression(expr.(*ast.IndexExpression).Index)
	if !types.Equals(indexType, &types.IntType{}) {
		chk.errors = append(chk.errors, fmt.Errorf("index expression has non-integer index"))
		return &types.IllegalType{}
	}

	arrType, ok := chk.checkExpression(expr.(*ast.IndexExpression).Left).(*types.ArrayType)
	if !ok {
		chk.errors = append(chk.errors, fmt.Errorf("index expression has non-array left operand"))
		return &types.IllegalType{}
	}

	return arrType.ElementsType
}

func (chk *TypeChecker) checkIdentifier(expr ast.ExpressionNode) types.Type {
	identType, ok := chk.env.Get(expr.(*ast.Identifier).Value)
	if !ok {
		chk.errors = append(chk.errors, fmt.Errorf("unknown identifier: %s", expr.(*ast.Identifier).Value))
		return &types.IllegalType{}
	}
	return identType
}

func (chk *TypeChecker) checkPrefixExpression(expr ast.ExpressionNode) types.Type {
	op := expr.(*ast.PrefixExpression).Operator
	rightType := chk.checkExpression(expr.(*ast.PrefixExpression).Right)
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
	chk.errors = append(chk.errors, fmt.Errorf("unknown operator: %s for type %s", op, rightType.Signature()))
	return &types.IllegalType{}
}

func (chk *TypeChecker) checkInfixExpression(expr ast.ExpressionNode) types.Type {
	leftType := chk.checkExpression(expr.(*ast.InfixExpression).Left)
	rightType := chk.checkExpression(expr.(*ast.InfixExpression).Right)
	op := expr.(*ast.InfixExpression).Operator

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
	}

	chk.errors = append(chk.errors, fmt.Errorf("type mismatch between: %s & %s; for operator %s",
		leftType.Signature(), rightType.Signature(), op))
	return &types.IllegalType{}
}

func (chk *TypeChecker) checkIfExpression(expr ast.ExpressionNode) types.Type {
	condType := chk.checkExpression(expr.(*ast.IfExpression).Condition)

	if !types.Equals(condType, &types.BoolType{}) {
		chk.errors = append(chk.errors, fmt.Errorf("wrong type %s for if condition", condType.Signature()))
		return &types.IllegalType{}
	}

	conseqType := chk.checkBlockStatement(expr.(*ast.IfExpression).Consequence)

	if expr.(*ast.IfExpression).Alternative != nil {
		alterType := chk.checkBlockStatement(expr.(*ast.IfExpression).Alternative)

		if !types.Equals(conseqType, alterType) {
			chk.errors = append(chk.errors, fmt.Errorf("type mismatch between %s & %s in if/else branches",
				conseqType.Signature(), alterType.Signature()))
			return &types.IllegalType{}
		}
	}

	return conseqType
}

func (chk *TypeChecker) checkFunctionLiteral(expr ast.ExpressionNode) types.Type {
	resFuncType := &types.FuncType{}
	funLit := expr.(*ast.FunctionLiteral)
	var expectedFuncType types.FuncType

	// hasSelfParams := false
	if len(funLit.ParamTypes) != 0 && funLit.ParamTypes[0] != nil || funLit.ReturnType != nil {

		if len(funLit.ParamTypes) != 0 && funLit.ParamTypes[0] != nil {
			for _, param := range funLit.ParamTypes {
				expectedFuncType.ParamsTypes = append(expectedFuncType.ParamsTypes, chk.resolveType(param))
			}
		} else if tp, ok := chk.curExpectedType.(*types.FuncType); ok && len(tp.ParamsTypes) != 0 {
			expectedFuncType.ParamsTypes = tp.ParamsTypes
		} else if len(funLit.Parameters) != 0 {
			chk.errors = append(chk.errors, fmt.Errorf("function literal has parameters, but no type specified for them"))
			return &types.IllegalType{}
		}

		if funLit.ReturnType != nil {
			expectedFuncType.ReturnType = chk.resolveType(funLit.ReturnType)
		} else {
			chk.errors = append(chk.errors, fmt.Errorf("function literal has typed parameters, but no return type specified"))
			return &types.IllegalType{}
		}

	} else {
		tempType, ok := chk.curExpectedType.(*types.FuncType)
		if !ok {
			return &types.IllegalType{}
		}
		expectedFuncType = *tempType
	}

	if len(expectedFuncType.ParamsTypes) != len(funLit.Parameters) {
		chk.errors = append(chk.errors, fmt.Errorf("function literal has %d parameters, but expected %d",
			len(funLit.Parameters), len(expectedFuncType.ParamsTypes)))
		return &types.IllegalType{}
	}

	chk.env = types.NewEnclosedTypeEviroment(chk.env)

	for i, param := range funLit.Parameters {
		resFuncType.ParamsTypes = append(resFuncType.ParamsTypes, expectedFuncType.ParamsTypes[i])
		chk.env.Set(param.Value, expectedFuncType.ParamsTypes[i])
	}

	oldExp := chk.curExpectedType

	chk.curExpectedType = expectedFuncType.ReturnType

	resFuncType.ReturnType = chk.checkBlockStatement(funLit.Body)

	if !types.Equals(resFuncType.ReturnType, expectedFuncType.ReturnType) {
		chk.errors = append(chk.errors, fmt.Errorf("function literal has return type %s, but expected %s",
			resFuncType.ReturnType.Signature(), expectedFuncType.ReturnType.Signature()))
		return &types.IllegalType{}
	}

	chk.curExpectedType = oldExp

	chk.env = chk.env.Outer

	return resFuncType
}

func (chk *TypeChecker) checkCallExpression(expr ast.ExpressionNode) types.Type {
	callExpr := expr.(*ast.CallExpression)
	rawCallType := chk.checkExpression(callExpr.Function)

	if _, err := rawCallType.(*types.IllegalType); err {
		return &types.IllegalType{}
	}

	switch callFuncType := rawCallType.(type) {
	case *types.FuncType:
		if len(callExpr.Arguments) != len(callFuncType.ParamsTypes) {
			chk.errors = append(chk.errors, fmt.Errorf(`wrong number of arguments for "%s"`, callExpr.Function.String()))
			return &types.IllegalType{}
		}

		for i, arg := range callExpr.Arguments {
			argType := chk.checkExpression(arg)

			if !types.Equals(argType, callFuncType.ParamsTypes[i]) {
				chk.errors = append(chk.errors, fmt.Errorf("wrong type for argument %d: %s in func call \" %s \"; expected %s",
					i+1, argType.Signature(), callExpr.Function.String(), callFuncType.ParamsTypes[i].Signature()))
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
			chk.errors = append(chk.errors, err)
			return &types.IllegalType{}
		}

		return returnType
	default:
		chk.errors = append(chk.errors, fmt.Errorf(`identifier "%s" is no a funciton`, callExpr.Function.String()))
		return &types.IllegalType{}
	}

}
