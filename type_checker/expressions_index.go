package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/types"
)

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
