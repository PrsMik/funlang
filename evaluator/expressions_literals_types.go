package evaluator

import (
	"funlang/ast"
	"funlang/object"
	"funlang/types"
)

func evalInterfaceLiteral(node *ast.InterfaceLiteral, env *object.Environment, info *types.Info) object.Object {
	ifaceType, ok := info.TypesInfo[node].(*types.TypeType)
	// TODO: перепроверить, что здесь всё безопасно для рантайма 1
	if !ok {
		return &object.Error{Message: "runtime error. interface is not a type. in theory this should never happen :/"}
	}

	var dynamicInterfaceType types.Type = ifaceType
	if statIfaceType, ok := ifaceType.Underlying.(*types.InterfaceType); ok {
		dynIfaceType := &types.InterfaceType{
			Name:   statIfaceType.Name,
			Decl:   statIfaceType.Decl,
			Fields: make(map[string]types.Type),
		}

		if node.Body != nil {
			for _, statement := range node.Body.Statements {
				letStmt, isLet := statement.(*ast.LetStatement)
				if !isLet {
					return newError("runtime error. interface block can only contain 'let' statements")
				}

				fieldVal := Eval(letStmt.Value, env, info)
				if isError(fieldVal) {
					return fieldVal
				}

				tObj, isType := fieldVal.(*object.TypeObject)
				if !isType {
					return newError("runtime error. field '%s' of interface is not a type", letStmt.Name.Value)
				}

				dynIfaceType.Fields[letStmt.Name.Value] = tObj.TypeValue
			}
		}

		dynamicInterfaceType = dynIfaceType
	}

	return &object.TypeObject{TypeValue: dynamicInterfaceType, Env: env}
}

func evalImplLiteral(node *ast.ImplLiteral, env *object.Environment, info *types.Info) object.Object {
	implType, ok := info.TypesInfo[node].(*types.TypeType)
	// TODO: перепроверить, что здесь всё безопасно для рантайма 2
	if !ok {
		return &object.Error{Message: "runtime error. impl is not a type. in theory this should never happen :/"}
	}
	var dynamicImplType types.Type = implType
	if statImplType, ok := implType.Underlying.(*types.ImplType); ok {
		dynImplType := &types.ImplType{
			Name:                statImplType.Name,
			Decl:                statImplType.Decl,
			ImplementedType:     statImplType.ImplementedType,
			Env:                 env,
			Fields:              make(map[string]types.Type),
			FieldsDefaultValues: make(map[string]ast.ExpressionNode),
		}

		if node.TargetType != nil {
			rtVal := Eval(node.TargetType, env, info)
			if isError(rtVal) {
				return rtVal
			}

			tObj, isType := rtVal.(*object.TypeObject)
			if !isType {
				return newError("runtime error. implemented type of impl is not a type")
			}

			dynImplType.ImplementedType = tObj.TypeValue
		}

		if node.Body != nil {
			for _, stmt := range node.Body.Statements {
				letStmt, isLet := stmt.(*ast.LetStatement)
				if !isLet {
					return newError("runtime error. impl block can only contain 'let' statements")
				}

				fieldType := Eval(letStmt.Type, env, info)
				if isError(fieldType) {
					return fieldType
				}

				tObj, isType := fieldType.(*object.TypeObject)
				if !isType {
					return newError("runtime error. type of field '%s' of impl is not a type", letStmt.Name.Value)
				}

				dynImplType.Fields[letStmt.Name.Value] = tObj.TypeValue
				dynImplType.FieldsOrder = append(dynImplType.FieldsOrder, letStmt.Name.Value)
				dynImplType.FieldsDefaultValues[letStmt.Name.Value] = letStmt.Value
			}
		}

		dynamicImplType = dynImplType
	}
	return &object.TypeObject{TypeValue: dynamicImplType, Env: env}
}
