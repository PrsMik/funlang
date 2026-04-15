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

func (chk *TypeChecker) resolveType(inType ast.TypeNode) types.Type {
	chk.recordTypeNode(inType, true)
	switch tp := inType.(type) {
	case *ast.SimpleType:
		switch tp.Token.Type {
		case token.INT_TYPE:
			return &types.IntType{}
		case token.BOOL_TYPE:
			return &types.BoolType{}
		case token.STRING_TYPE:
			return &types.StringType{}
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
	}
	chk.typeError(fmt.Sprintf("%s is not a valid type", inType.String()), inType)
	return nil
}
