package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/types"
)

func (chk *TypeChecker) checkInterfaceLiteral(expr ast.ExpressionNode) types.Type {
	interfaceLiteral := expr.(*ast.InterfaceLiteral)

	interfaceType := &types.InterfaceType{
		Name:   "",
		Fields: make(map[string]types.Type),
	}

	for _, stmt := range interfaceLiteral.Body.Statements {
		letStmt, ok := stmt.(*ast.LetStatement)
		if !ok {
			chk.typeError("interface block can only contain 'let' statements for type definitions", stmt)
			continue
		}

		fieldDeclType := chk.resolveType(letStmt.Type)
		if _, isTypeType := fieldDeclType.(*types.TypeType); !isTypeType {
			chk.typeError(fmt.Sprintf("interface fields must be defined as types, got %s", fieldDeclType.Signature()), letStmt)
			continue
		}

		if _, exists := interfaceType.Fields[letStmt.Name.Value]; exists {
			chk.typeError(fmt.Sprintf("duplicate field '%s' in interface", letStmt.Name.Value), letStmt)
			continue
		}

		fieldType := chk.resolveType(letStmt.Value)
		if _, isIllegal := fieldType.(*types.IllegalType); isIllegal {
			continue
		}

		interfaceType.Fields[letStmt.Name.Value] = fieldType

		chk.recordType(letStmt.Name, fieldType)
		chk.recordExpectedType(letStmt.Value, fieldDeclType)
	}

	return &types.TypeType{Underlying: interfaceType}
}

func (chk *TypeChecker) checkImplLiteral(expr ast.ExpressionNode) types.Type {
	implLit := expr.(*ast.ImplLiteral)

	var targetType types.Type
	if implLit.TargetType != nil {
		targetType = chk.resolveType(implLit.TargetType)
		if _, isIllegal := targetType.(*types.IllegalType); isIllegal {
			return &types.IllegalType{}
		}
	}

	implType := &types.ImplType{
		Name:            "",
		ImplementedType: targetType,
		Fields:          make(map[string]types.Type),
	}

	chk.env = types.NewEnclosedTypeEviroment(chk.env)
	chk.recordScope(implLit.Body, chk.env)

	for _, stmt := range implLit.Body.Statements {
		letStmt, ok := stmt.(*ast.LetStatement)
		if !ok {
			chk.typeError("impl block can only contain 'let' statements", stmt)
			continue
		}

		fieldType := chk.checkLetStatement(letStmt)
		if _, isIllegal := fieldType.(*types.IllegalType); isIllegal {
			continue
		}

		if _, exists := implType.Fields[letStmt.Name.Value]; exists {
			chk.typeError(fmt.Sprintf("duplicate field '%s' in impl", letStmt.Name.Value), letStmt)
			continue
		}

		implType.Fields[letStmt.Name.Value] = fieldType
	}

	chk.env = chk.env.Outer

	if targetType != nil {
		chk.verifyImplementation(targetType, implType, expr)
	}

	return &types.TypeType{Underlying: implType}
}

func (chk *TypeChecker) verifyImplementation(targetType types.Type, implType *types.ImplType, node ast.Node) {
	switch rawType := targetType.(type) {
	case *types.InterfaceType:
		for interfaceFieldName, interfaceFieldType := range rawType.Fields {
			implFieldType, exists := implType.Fields[interfaceFieldName]
			if !exists {
				chk.typeError(fmt.Sprintf("impl missing required field '%s'", interfaceFieldName), node)
				continue
			}
			if !types.IsAssignable(interfaceFieldType, implFieldType) {
				chk.typeError(fmt.Sprintf("field '%s' type mismatch: expected %s, got %s",
					interfaceFieldName, interfaceFieldType.Signature(), implFieldType.Signature()), node)
			}
		}
	case *types.IntersectionType:
		chk.verifyImplementation(rawType.Left, implType, node)
		chk.verifyImplementation(rawType.Right, implType, node)
	case *types.TypeType:
		chk.verifyImplementation(rawType.Underlying, implType, node)
	}
}
