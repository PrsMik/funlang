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
	// oldEnv := env
	// env = object.NewEnclosedEnvironment(env)

	// defer func() {
	// 	env = oldEnv
	// }()

	var dynamicInterfaceType types.Type = ifaceType
	if statIfaceType, ok := ifaceType.Underlying.(*types.InterfaceType); ok {
		dynIfaceType := &types.InterfaceType{
			Name:   statIfaceType.Name,
			Fields: make(map[string]types.Type),
		}

		if node.Body != nil {
			for _, statement := range node.Body.Statements {
				letStmt, _ := statement.(*ast.LetStatement)
				evalRes := Eval(statement, env, info)
				if tObj, isType := evalRes.(*object.TypeObject); isType {
					dynIfaceType.Fields[letStmt.Name.Value] = tObj.TypeValue
				}
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
			ImplementedType:     statImplType.ImplementedType,
			Fields:              make(map[string]types.Type),
			FieldsDefaultValues: make(map[string]ast.ExpressionNode),
		}

		if node.TargetType != nil {
			rtVal := Eval(node.TargetType, env, info)
			if tObj, isType := rtVal.(*object.TypeObject); isType {
				dynImplType.ImplementedType = tObj.TypeValue
			}
		}

		if node.Body != nil {
			for _, stmt := range node.Body.Statements {
				letStmt, _ := stmt.(*ast.LetStatement)
				fieldType := Eval(letStmt.Type, env, info)
				if tObj, isType := fieldType.(*object.TypeObject); isType {
					dynImplType.Fields[letStmt.Name.Value] = tObj.TypeValue
				}
				dynImplType.FieldsDefaultValues[letStmt.Name.Value] = letStmt.Value
			}
		}

		dynamicImplType = dynImplType
	}
	return &object.TypeObject{TypeValue: dynamicImplType, Env: env}
}
