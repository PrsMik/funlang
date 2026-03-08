package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/token"
	"funlang/types"
)

func (chk *TypeChecker) resolveType(inType ast.TypeNode) types.Type {
	switch tp := inType.(type) {
	case *ast.SimpleType:
		switch tp.Token.Type {
		case token.INT_TYPE:
			return &types.IntType{}
		case token.BOOL_TYPE:
			return &types.BoolType{}
		}
	case *ast.FunctionType:
		prmTypes := []types.Type{}

		for _, prm := range tp.ParamsTypes {
			resolvedParam := chk.resolveType(prm)
			prmTypes = append(prmTypes, resolvedParam)
		}

		rtrnType := chk.resolveType(tp.ReturnType)

		return &types.FuncType{ParamsTypes: prmTypes, ReturnType: rtrnType}
	}
	chk.errors = append(chk.errors, fmt.Errorf("unknown type"))
	return nil
}
