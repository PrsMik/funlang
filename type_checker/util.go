package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/token"
	"funlang/types"
)

func (chk *TypeChecker) recordType(node ast.Node, tp types.Type) {
	if chk.Info != nil && node != nil {
		chk.Info.TypesInfo[node] = tp
	}
}

func (chk *TypeChecker) recordTypeNode(node ast.Node, bl bool) {
	if chk.Info != nil && node != nil {
		chk.Info.TypeNodes[node] = bl
	}
}

func (chk *TypeChecker) recordDefinition(usage ast.Node, decl ast.Node) {
	if chk.Info != nil && usage != nil && decl != nil {
		chk.Info.Definitions[usage] = decl
	}
}

func (chk *TypeChecker) recordScope(node ast.Node, env *types.TypeEviroment) {
	if chk.Info != nil && node != nil {
		chk.Info.Scopes[node] = env
	}
}

func (chk *TypeChecker) recordExpectedType(node ast.Node, tp types.Type) {
	if chk.Info != nil && node != nil && tp != nil {
		chk.Info.ExpectedTypes[node] = tp
	}
}

func (chk *TypeChecker) resolveTypeExpression(expr ast.ExpressionNode) types.Type {
	evaluatedType := chk.checkExpression(expr)

	if typeType, ok := evaluatedType.(*types.TypeType); ok {
		return typeType.Underlying
	}

	if _, ok := evaluatedType.(*types.IllegalType); !ok {
		chk.typeError(fmt.Sprintf("expected a type, but got a value of type %s", evaluatedType.Signature()), expr)
	}

	return &types.IllegalType{}
}

func (chk *TypeChecker) resolveType(inType ast.ExpressionNode) types.Type {
	chk.recordTypeNode(inType, true)
	switch tp := inType.(type) {
	case *ast.TypeType:
		res, _ := chk.env.Get("type")
		return res.SymbolType.(*types.TypeType)
	case *ast.SimpleType:
		switch tp.Token.Type {
		case token.INT_TYPE:
			res, _ := chk.env.Get("int")
			return res.SymbolType.(*types.TypeType).Underlying
		case token.BOOL_TYPE:
			res, _ := chk.env.Get("bool")
			return res.SymbolType.(*types.TypeType).Underlying
		case token.STRING_TYPE:
			res, _ := chk.env.Get("string")
			return res.SymbolType.(*types.TypeType).Underlying
		}
	case *ast.ArrayType:
		return &types.ArrayType{ElementsType: chk.resolveType(tp.ElementsType)}
	case *ast.HashMapType:
		keyType := chk.resolveType(tp.KeyType)
		if _, ok := keyType.(types.HashableType); !ok {
			chk.typeError(fmt.Sprintf("cannot use type %s for hash map key", keyType.Signature()), tp)
			return &types.IllegalType{}
		}
		return &types.HashMapType{KeyType: keyType, ElementType: chk.resolveType(tp.ElementType)}
	case *ast.FunctionType:
		funcParams := []types.FuncParam{}

		for _, prm := range tp.ParamsTypes {
			resolvedParamType := chk.resolveType(prm)
			funcParams = append(funcParams, types.FuncParam{Name: "", Type: resolvedParamType})
		}

		rtrnType := chk.resolveType(tp.ReturnType)

		return &types.FuncType{Params: funcParams, ReturnType: rtrnType}
	default:
		return chk.resolveTypeExpression(inType)
	}
	chk.typeError(fmt.Sprintf("%s is not a valid type", inType.String()), inType)
	return nil
}
