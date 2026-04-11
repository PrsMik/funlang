package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/token"
	"funlang/types"
)

func (chk *TypeChecker) resolveType(inType ast.TypeNode) types.Type {
	chk.TypeNodes[inType] = true
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
