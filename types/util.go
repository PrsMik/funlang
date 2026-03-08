package types

import (
	"fmt"
	"funlang/ast"
	"funlang/token"
)

func (chk *TypeChecker) resolveType(inType ast.TypeNode) Type {
	switch tp := inType.(type) {
	case *ast.SimpleType:
		switch tp.Token.Type {
		case token.INT_TYPE:
			return &IntType{}
		case token.BOOL_TYPE:
			return &BoolType{}
		}
	case *ast.FunctionType:
		prmTypes := []Type{}

		for _, prm := range tp.ParamsTypes {
			resolvedParam := chk.resolveType(prm)
			prmTypes = append(prmTypes, resolvedParam)
		}

		rtrnType := chk.resolveType(tp.ReturnType)

		return &FuncType{prmTypes, rtrnType}
	}
	chk.errors = append(chk.errors, fmt.Errorf("unknown type"))
	return nil
}
