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
		chk.typeError(fmt.Sprintf("type %s has no field named *%s*", leftType.Signature(), fieldName), expr)
		return &types.IllegalType{}
	}

	return fieldType
}

func (chk *TypeChecker) getFieldType(objectType types.Type, fieldName string) (types.Type, bool) {
	switch rawType := objectType.(type) {
	case *types.InterfaceType:
		fieldType, ok := rawType.Fields[fieldName]
		if fieldType == nil {
			return &types.IllegalType{}, ok
		}
		return fieldType, ok

	case *types.ImplType:
		fieldType, ok := rawType.Fields[fieldName]
		if fieldType == nil {
			return &types.IllegalType{}, ok
		}
		return fieldType, ok

	case *types.IntersectionType:
		if fieldType, ok := chk.getFieldType(rawType.Left, fieldName); ok {
			return fieldType, true
		}
		if fieldType, ok := chk.getFieldType(rawType.Right, fieldName); ok {
			return fieldType, true
		}
		return &types.IllegalType{}, false
	}

	return &types.IllegalType{}, false
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

func (chk *TypeChecker) checkImplInstantiationExpression(expr ast.ExpressionNode) types.Type {
	instantationExpr := expr.(*ast.ImplInstantiationExpression)

	leftType := chk.checkExpression(instantationExpr.Left)
	if _, isIllegal := leftType.(*types.IllegalType); isIllegal {
		return leftType
	}

	var impl *types.ImplType

	if typeType, isType := leftType.(*types.TypeType); !isType {
		chk.typeError("cannot use instantation .{} on a no-type expression, it can only be used on impl types", expr)
		return &types.IllegalType{}
	} else {
		under, isImpl := typeType.Underlying.(*types.ImplType)
		if !isImpl {
			chk.typeError("cannot use instantation .{} on a no-impl types, it can only be used on impl type", expr)
			return &types.IllegalType{}
		}
		impl = under
	}

	fields := instantationExpr.Fields.(*ast.HashMapLiteral)

	chk.checkFieldsInstantation(impl, fields, instantationExpr)

	return impl
}

func (chk *TypeChecker) checkFieldsInstantation(implType *types.ImplType, fields *ast.HashMapLiteral, node ast.Node) {
	for key, elem := range fields.Pairs {
		instElementType := chk.checkExpression(elem)

		keyIndent, ok := key.(*ast.Identifier)
		if !ok {
			chk.typeError(fmt.Sprintf("%s is not an identifier", keyIndent.String()), node)
			return
		}

		fieldType, ok := chk.getFieldType(implType, keyIndent.Value)
		if !ok {
			chk.typeError(fmt.Sprintf("type %s has no field named *%s*", implType.Signature(), keyIndent.Value), node)
		}

		if !types.IsAssignable(fieldType, instElementType) {
			chk.typeError(fmt.Sprintf(`type mismatch in instatation expression of %s between field named *%s* with type %s, and value %s with type %s`,
				implType.Signature(), keyIndent.Value, fieldType.Signature(), elem.String(), instElementType.Signature()), node)
		}
	}
}
