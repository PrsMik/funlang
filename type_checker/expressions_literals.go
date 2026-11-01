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
func (chk *TypeChecker) checkArrayType(expr ast.ExpressionNode) types.Type {
	return &types.TypeType{Underlying: chk.resolveType(expr)}
}

func (chk *TypeChecker) checkHashMapType(expr ast.ExpressionNode) types.Type {
	return &types.TypeType{Underlying: chk.resolveType(expr)}
}

func (chk *TypeChecker) checkFuncType(expr ast.ExpressionNode) types.Type {
	return &types.TypeType{Underlying: chk.resolveType(expr)}
}

func (chk *TypeChecker) checkArrayLiteral(expr ast.ExpressionNode) types.Type {
	arrType := &types.ArrayType{}
	arrLit := expr.(*ast.ArrayLiteral)

	if len(arrLit.Elements) == 0 {
		return arrType
	}

	oldType := chk.curExpectedType
	defer func() { chk.curExpectedType = oldType }()

	var expectedElemType types.Type
	// ожидается тип массива
	if expType, ok := chk.curExpectedType.(*types.ArrayType); ok {
		chk.curExpectedType = expType.ElementsType
		if _, isIllegal := expType.ElementsType.(*types.IllegalType); !isIllegal {
			expectedElemType = expType.ElementsType
		}
	}

	// фактические типы элементов
	elementTypes := make([]types.Type, len(arrLit.Elements))
	for i, element := range arrLit.Elements {
		elementTypes[i] = chk.checkExpression(element)
	}

	if expectedElemType != nil {
		// тип элементов - ожидаемый тип
		for i, elementType := range elementTypes {
			if _, isIllegal := elementType.(*types.IllegalType); isIllegal {
				// ошибка уже в элементе
				return &types.IllegalType{}
			}

			if !types.IsAssignable(expectedElemType, elementType) {
				chk.typeError(fmt.Sprintf("array literal has element of type %s, but expected %s",
					elementType.Signature(), expectedElemType.Signature()), arrLit.Elements[i])
				return &types.IllegalType{}
			}
		}

		arrType.ElementsType = expectedElemType
	} else {
		// тип элементов - фактический тип
		for _, curParamType := range elementTypes[1:] {
			if !types.Equals(elementTypes[0], curParamType) {
				chk.typeError(fmt.Sprintf("array literal has elements of different types %s & %s",
					elementTypes[0].Signature(), curParamType.Signature()), expr)
				return &types.IllegalType{}
			}
		}

		arrType.ElementsType = elementTypes[0]
	}

	if _, ok := arrLit.Elements[0].(*ast.BadExpression); len(arrLit.Elements) == 1 && ok {
		arrLit.Elements = arrLit.Elements[1:]
	}

	return arrType
}

func (chk *TypeChecker) checkHashMapLiteral(expr ast.ExpressionNode) types.Type {
	hashMapType := &types.HashMapType{}
	hashMapLiteral := expr.(*ast.HashMapLiteral)

	if len(hashMapLiteral.Pairs) == 0 {
		return hashMapType
	}

	oldExpectedType := chk.curExpectedType
	defer func() { chk.curExpectedType = oldExpectedType }()

	var keyExpectedType types.Type = nil
	var elemExpectedType types.Type = nil

	if tp, ok := chk.curExpectedType.(*types.HashMapType); ok {
		keyExpectedType = tp.KeyType
		elemExpectedType = tp.ElementType
	}

	useExpected := keyExpectedType != nil && elemExpectedType != nil
	if _, isIllegal := keyExpectedType.(*types.IllegalType); isIllegal {
		useExpected = false
	}
	if _, isIllegal := elemExpectedType.(*types.IllegalType); isIllegal {
		useExpected = false
	}

	var firstHashMapKeyType types.Type = nil
	var firstHashMapElementType types.Type = nil

	for key, elem := range hashMapLiteral.Pairs {
		if _, isBad := key.(*ast.BadExpression); isBad {
			continue
		}

		chk.curExpectedType = keyExpectedType
		curHashMapKeyType := chk.checkExpression(key)

		chk.curExpectedType = elemExpectedType
		curHashMapElementType := chk.checkExpression(elem)

		if useExpected {
			if !types.IsAssignable(keyExpectedType, curHashMapKeyType) {
				chk.typeError(fmt.Sprintf("map literal has key of type %s, but expected %s",
					curHashMapKeyType.Signature(), keyExpectedType.Signature()), key)
				return &types.IllegalType{}
			}

			if !types.IsAssignable(elemExpectedType, curHashMapElementType) {
				chk.typeError(fmt.Sprintf("map literal has element of type %s, but expected %s",
					curHashMapElementType.Signature(), elemExpectedType.Signature()), elem)
				return &types.IllegalType{}
			}

			continue
		}

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

	if useExpected {
		hashMapType.KeyType = keyExpectedType
		hashMapType.ElementType = elemExpectedType
	} else {
		hashMapType.KeyType = firstHashMapKeyType
		hashMapType.ElementType = firstHashMapElementType
	}

	return hashMapType
}
