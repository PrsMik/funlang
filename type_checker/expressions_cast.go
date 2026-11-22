package type_checker

import (
	"funlang/ast"
	"funlang/types"
)

func (chk *TypeChecker) checkSafeTypeCastExpression(expr ast.ExpressionNode) types.Type {
	cast := expr.(*ast.SafeTypeCastExpression)
	leftType := chk.checkExpression(cast.Left)

	if _, isIllegal := leftType.(*types.IllegalType); isIllegal {
		return leftType
	}

	targetType := chk.resolveType(cast.TargetType)
	if _, isIllegal := targetType.(*types.IllegalType); isIllegal {
		return targetType
	}

	return &types.UnionType{Left: targetType, Right: &types.NilType{}}
}

func (chk *TypeChecker) checkUnsafeTypeCastExpression(expr ast.ExpressionNode) types.Type {
	cast := expr.(*ast.UnsafeTypeCastExpression)
	targetType := chk.resolveType(cast.TargetType)
	if _, isIllegal := targetType.(*types.IllegalType); isIllegal {
		return targetType
	}

	_ = chk.checkExpression(cast.Value)

	return targetType
}
