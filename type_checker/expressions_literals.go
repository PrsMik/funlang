package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/types"
)

func (chk *TypeChecker) checkIntegerLiteral(expr ast.ExpressionNode) types.Type {
	return &types.IntType{}
}

func (chk *TypeChecker) checkBooleanLiteral(expr ast.ExpressionNode) types.Type {
	return &types.BoolType{}
}

func (chk *TypeChecker) checkStringLiteral(expr ast.ExpressionNode) types.Type {
	return &types.StringType{}
}

func (chk *TypeChecker) checkIdentifier(expr ast.ExpressionNode) types.Type {
	symbolInfo, ok := chk.env.Get(expr.(*ast.Identifier).Value)
	if !ok {
		chk.typeError(fmt.Sprintf("unknown identifier: %s", expr.(*ast.Identifier).Value), expr)
		return &types.IllegalType{}
	}

	if symbolInfo.DeclNode != nil {
		chk.recordDefinition(expr, symbolInfo.DeclNode)
	}

	return symbolInfo.SymbolType
}

func (chk *TypeChecker) checkTypeIdentifier(expr ast.ExpressionNode) types.Type {
	symbolInfo, ok := chk.env.Get(expr.String())
	if !ok {
		chk.typeError(fmt.Sprintf("unknown identifier: %s", expr.String()), expr)
		return &types.IllegalType{}
	}

	if symbolInfo.DeclNode != nil {
		chk.recordDefinition(expr, symbolInfo.DeclNode)
	}

	return symbolInfo.SymbolType
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

	if _, ok := arrLit.Elements[0].(*ast.BadExpression); len(arrLit.Elements) == 1 && ok {
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
	chk.recordScope(expr, chk.env)

	// спуск типов из ожидаемого типа функции на переменные в литерале
	for i, param := range funLiteral.Parameters {
		parameter := types.FuncParam{Name: param.Value, Type: expectedFuncType.Params[i].Type}
		resFuncType.Params = append(resFuncType.Params, parameter)

		chk.recordType(param, expectedFuncType.Params[i].Type)
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
