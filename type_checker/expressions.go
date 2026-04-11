package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/token"
	"funlang/types"
	"reflect"
)

func (chk *TypeChecker) checkExpression(expr ast.ExpressionNode) types.Type {
	exprType := reflect.TypeOf(expr)
	if chk.curExpectedType != nil {
		chk.ExpectedTypes[expr] = chk.curExpectedType
	}
	if checkFun, ok := chk.expressionCheckFns[exprType]; ok {
		tp := checkFun(expr)
		chk.TypesInfo[expr] = tp
		return tp
	}
	chk.typeError("unknown expression type", expr)
	return &types.IllegalType{}
}

func (chk *TypeChecker) checkVirtualNode(expr ast.ExpressionNode) types.Type {
	if chk.curExpectedType != nil {
		chk.ExpectedTypes[expr] = chk.curExpectedType
		return chk.curExpectedType
	}
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

	oldType := chk.curExpectedType

	expType, ok := chk.curExpectedType.(*types.ArrayType)
	if ok {
		chk.curExpectedType = expType.ElementsType
	}

	firstArrLitType := chk.checkExpression(arrLit.Elements[0])

	for _, param := range arrLit.Elements[1:] {
		curParamType := chk.checkExpression(param)
		if !types.Equals(firstArrLitType, curParamType) {
			// if tp, ok := arrLit.Elements[index].(*ast.UnparsedNode); ok {

			// }
			chk.typeError(fmt.Sprintf("array literal has elements of different types %s & %s",
				firstArrLitType.Signature(), curParamType.Signature()), expr)
			return &types.IllegalType{}
		}
	}

	chk.curExpectedType = oldType

	if _, ok := arrLit.Elements[0].(*ast.VirtualNode); len(arrLit.Elements) == 1 && ok {
		arrLit.Elements = arrLit.Elements[1:]
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

	oldExpectedType := chk.curExpectedType
	var keyExpectedType types.Type = nil
	var elemExpectedType types.Type = nil

	if tp, ok := chk.curExpectedType.(*types.HashMapType); ok {
		keyExpectedType = tp.KeyType
		elemExpectedType = tp.ElementType
	}

	var firstHashMapKeyType types.Type = nil
	var firstHashMapElementType types.Type = nil

	for key, elem := range hashMapLiteral.Pairs {

		chk.curExpectedType = keyExpectedType
		curHashMapKeyType := chk.checkExpression(key)

		chk.curExpectedType = elemExpectedType
		curHashMapElementType := chk.checkExpression(elem)

		if firstHashMapKeyType == nil && firstHashMapElementType == nil {
			firstHashMapKeyType = curHashMapKeyType
			firstHashMapElementType = curHashMapElementType
		} else {
			if !types.Equals(firstHashMapKeyType, curHashMapKeyType) {
				chk.typeError(fmt.Sprintf("map literal has keys of different types %s & %s",
					firstHashMapKeyType.Signature(), curHashMapKeyType.Signature()), expr)
				return &types.IllegalType{}
			}

			if !types.Equals(firstHashMapElementType, curHashMapElementType) {
				chk.typeError(fmt.Sprintf("map literal has elements of different types %s & %s",
					firstHashMapElementType.Signature(), curHashMapElementType.Signature()), expr)
				return &types.IllegalType{}
			}
		}
	}

	chk.curExpectedType = oldExpectedType
	hashMapType.KeyType = firstHashMapKeyType
	hashMapType.ElementType = firstHashMapElementType

	return hashMapType
}

func (chk *TypeChecker) checkIndexExpression(expr ast.ExpressionNode) types.Type {
	leftExpr := chk.checkExpression(expr.(*ast.IndexExpression).Left)

	switch leftType := leftExpr.(type) {
	case *types.ArrayType:
		oldType := chk.curExpectedType
		chk.curExpectedType = &types.IntType{}

		res := chk.checkArrayIndexExpression(expr)

		chk.curExpectedType = oldType
		return res
	case *types.HashMapType:
		oldType := chk.curExpectedType
		chk.curExpectedType = leftType.KeyType

		// fmt.Fprintf(os.Stderr, "Checking map with expected %T\n", chk.curExpectedType)

		res := chk.checkHashMapIndexExpression(expr)

		chk.curExpectedType = oldType
		return res
	default:
		chk.typeError("index expression has wrong left operand", expr)
		return &types.IllegalType{}
	}
}

func (chk *TypeChecker) checkArrayIndexExpression(expr ast.ExpressionNode) types.Type {
	arrType := chk.checkExpression(expr.(*ast.IndexExpression).Left).(*types.ArrayType)

	indexType := chk.checkExpression(expr.(*ast.IndexExpression).Index)

	if !types.Equals(indexType, &types.IntType{}) {
		chk.typeError("array index expression has non-integer index", expr)
		return &types.IllegalType{}
	}

	return arrType.ElementsType
}

func (chk *TypeChecker) checkHashMapIndexExpression(expr ast.ExpressionNode) types.Type {
	hashMapType := chk.checkExpression(expr.(*ast.IndexExpression).Left).(*types.HashMapType)

	indexType := chk.checkExpression(expr.(*ast.IndexExpression).Index)

	if hashMapType.KeyType == nil {
		chk.typeError("index operator usage for map with unkwown key type", expr)
		return &types.IllegalType{}
	}

	if !types.Equals(indexType, hashMapType.KeyType) {
		chk.typeError(fmt.Sprintf("type mismatch between %s in index & %s in keys for index operator in hash map",
			indexType.Signature(), hashMapType.KeyType.Signature()), expr)
		return &types.IllegalType{}
	}

	return hashMapType.ElementType
}

func (chk *TypeChecker) checkIdentifier(expr ast.ExpressionNode) types.Type {
	symbolInfo, ok := chk.env.Get(expr.(*ast.Identifier).Value)
	if !ok {
		chk.typeError(fmt.Sprintf("unknown identifier: %s", expr.(*ast.Identifier).Value), expr)
		return &types.IllegalType{}
	}

	if symbolInfo.DeclNode != nil {
		chk.Definitions[expr] = symbolInfo.DeclNode
	}

	return symbolInfo.SymbolType
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
	}

	chk.typeError(fmt.Sprintf("type mismatch between: %s & %s; for operator %s",
		leftType.Signature(), rightType.Signature(), op), expr)
	return &types.IllegalType{}
}

func (chk *TypeChecker) checkIfExpression(expr ast.ExpressionNode) types.Type {
	// --- РАЗМЕТКА УСЛОВИЯ ---
	ifExpr := expr.(*ast.IfExpression)

	startPos := ifExpr.Start()
	var endPos token.Position

	if ifExpr.Consequence != nil {
		endPos = ifExpr.Consequence.Start()
	} else {
		endPos = ifExpr.End()
	}

	condArea := &ast.VirtualNode{
		From: startPos,
		To:   endPos,
	}

	chk.ExpectedTypes[condArea] = &types.BoolType{}

	// --- ПРОВЕРКА ТИПА ---
	oldType := chk.curExpectedType

	chk.curExpectedType = &types.BoolType{}

	condType := chk.checkExpression(expr.(*ast.IfExpression).Condition)

	chk.curExpectedType = oldType

	if !types.Equals(condType, &types.BoolType{}) {
		chk.typeError(fmt.Sprintf("wrong type %s for if condition", condType.Signature()), expr)
		return &types.IllegalType{}
	}

	conseqType := chk.checkBlockStatement(expr.(*ast.IfExpression).Consequence)

	if expr.(*ast.IfExpression).Alternative != nil {
		alterType := chk.checkBlockStatement(expr.(*ast.IfExpression).Alternative)

		if !types.Equals(conseqType, alterType) {
			chk.typeError(fmt.Sprintf("type mismatch between %s & %s in if/else branches",
				conseqType.Signature(), alterType.Signature()), expr)
			return &types.IllegalType{}
		}
	}

	return conseqType
}

func (chk *TypeChecker) checkFunctionLiteral(expr ast.ExpressionNode) types.Type {
	resFuncType := &types.FuncType{}
	funLiteral := expr.(*ast.FunctionLiteral)

	var expectedFuncType types.FuncType

	// hasSelfParams := false
	// типы параметров описаны в самом литерале
	if len(funLiteral.ParamTypes) != 0 && funLiteral.ParamTypes[0] != nil || funLiteral.ReturnType != nil {

		if len(funLiteral.ParamTypes) != 0 && funLiteral.ParamTypes[0] != nil {
			for ind, param := range funLiteral.ParamTypes {
				paramType := chk.resolveType(param)
				funcParam := &types.FuncParam{Name: funLiteral.Parameters[ind].Value, Type: paramType}
				expectedFuncType.Params = append(expectedFuncType.Params, *funcParam)
			}
		} else if expType, ok := chk.curExpectedType.(*types.FuncType); ok && len(expType.Params) != 0 {
			// у ожидаемого типа определены типы параметров
			expectedFuncType.Params = expType.Params
		} else if len(funLiteral.Parameters) != 0 {
			chk.typeError("function literal has parameters, but no type specified for them", expr)
			return &types.IllegalType{}
		}

		if funLiteral.ReturnType != nil {
			expectedFuncType.ReturnType = chk.resolveType(funLiteral.ReturnType)
		} else {
			chk.typeError("function literal has typed parameters, but no return type specified", expr)
			return &types.IllegalType{}
		}
	} else {
		tempType, ok := chk.curExpectedType.(*types.FuncType)
		if !ok {
			return &types.IllegalType{}
		}
		expectedFuncType = *tempType
	}

	if len(expectedFuncType.Params) != len(funLiteral.Parameters) {
		chk.typeError(fmt.Sprintf("function literal has %d parameters, but expected %d",
			len(funLiteral.Parameters), len(expectedFuncType.Params)), expr)
		return &types.IllegalType{}
	}

	chk.env = types.NewEnclosedTypeEviroment(chk.env)
	chk.Scopes[expr] = chk.env

	// спуск типов из ожидаемого типа функции на переменные в литерале
	for i, param := range funLiteral.Parameters {
		parameter := types.FuncParam{Name: param.Value, Type: expectedFuncType.Params[i].Type}
		resFuncType.Params = append(resFuncType.Params, parameter)

		chk.TypesInfo[param] = expectedFuncType.Params[i].Type
		chk.env.Set(param.Value, expectedFuncType.Params[i].Type, param)
	}

	oldExp := chk.curExpectedType

	chk.curExpectedType = expectedFuncType.ReturnType

	resFuncType.ReturnType = chk.checkBlockStatement(funLiteral.Body)

	if !types.Equals(resFuncType.ReturnType, expectedFuncType.ReturnType) {
		chk.typeError(fmt.Sprintf("function literal has return type %s, but expected %s",
			resFuncType.ReturnType.Signature(), expectedFuncType.ReturnType.Signature()), expr)
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
		if len(callExpr.Arguments) != len(callFuncType.Params) {
			if len(callFuncType.Params) != 0 {
				chk.typeError(fmt.Sprintf(`wrong number of arguments for "%s"`, callExpr.Function.String()), expr)
				return &types.IllegalType{}
			} else {
				if _, ok := callExpr.Arguments[0].(*ast.VirtualNode); ok {
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
