package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/types"
)

func (chk *TypeChecker) checkMemberAccessExpression(expr ast.ExpressionNode) types.Type {
	memberExpr := expr.(*ast.MemberAccessExpression)

	leftType := chk.checkExpression(memberExpr.Left)
	if _, isIllegal := leftType.(*types.IllegalType); isIllegal {
		return &types.IllegalType{}
	}

	fieldName := memberExpr.Property.Value
	fieldType, ok := chk.getFieldType(leftType, fieldName)
	if !ok {
		chk.typeError(fmt.Sprintf("type %s has no field '%s'", leftType.Signature(), fieldName), expr)
		return &types.IllegalType{}
	}

	return fieldType
}

func (chk *TypeChecker) getFieldType(objectType types.Type, fieldName string) (types.Type, bool) {
	switch rawType := objectType.(type) {
	case *types.InterfaceType:
		fieldType, ok := rawType.Fields[fieldName]
		return fieldType, ok

	case *types.ImplType:
		fieldType, ok := rawType.Fields[fieldName]
		return fieldType, ok

	case *types.IntersectionType:
		if fieldType, ok := chk.getFieldType(rawType.Left, fieldName); ok {
			return fieldType, true
		}
		if fieldType, ok := chk.getFieldType(rawType.Right, fieldName); ok {
			return fieldType, true
		}
		return nil, false
	}

	return nil, false
}

func (chk *TypeChecker) checkTypeAccessExpression(expr ast.ExpressionNode) types.Type {
	typeAccess := expr.(*ast.TypeAccessExpression)

	leftType := chk.checkExpression(typeAccess.Left)
	if _, isIllegal := leftType.(*types.IllegalType); isIllegal {
		return leftType
	}

	if _, isType := leftType.(*types.TypeType); isType {
		chk.typeError("cannot use .(type) on a type expression, it can only be used on values", expr)
		return &types.IllegalType{}
	}

	return &types.TypeType{Underlying: leftType}
}
