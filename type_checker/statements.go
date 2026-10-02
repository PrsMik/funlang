package type_checker

import (
	"fmt"
	"funlang/ast"
	"funlang/types"
)

func (chk *TypeChecker) checkStatement(stmt ast.StatementNode) types.Type {
	switch curStmt := stmt.(type) {
	case *ast.LetStatement:
		return chk.checkLetStatement(curStmt)
	case *ast.ReturnStatement:
		return chk.checkReturnStatement(curStmt)
	case *ast.BlockStatement:
		return chk.checkBlockStatement(curStmt)
	}
	chk.typeError("unknown statement type", stmt)
	return &types.IllegalType{}
}

func (chk *TypeChecker) checkLetStatement(stmt *ast.LetStatement) types.Type {
	if stmt == nil {
		return &types.IllegalType{}
	}

	if stmt.Type == nil {
		chk.typeError("missing type in let statement", stmt)
		return &types.IllegalType{}
	}

	expectedType := chk.resolveType(stmt.Type)

	chk.curExpectedType = expectedType

	chk.recordType(stmt.Name, expectedType)
	chk.recordExpectedType(stmt.Value, expectedType)

	_, isDeclaringType := expectedType.(*types.TypeType)

	// для самоссылающихся типов
	if isDeclaringType {
		var forwardDecl types.Type

		switch val := stmt.Value.(type) {
		case *ast.InterfaceLiteral:
			forwardDecl = &types.InterfaceType{
				Name:   stmt.Name.Value,
				Fields: make(map[string]types.Type)}
		case *ast.ImplLiteral:
			var targetType types.Type
			if val.TargetType != nil {
				targetType = chk.resolveType(val.TargetType)
			}

			forwardDecl = &types.ImplType{
				Name:                stmt.Name.Value,
				ImplementedType:     targetType,
				Fields:              make(map[string]types.Type),
				FieldsDefaultValues: make(map[string]ast.ExpressionNode)}
		}

		if forwardDecl != nil {
			chk.env.Set(stmt.Name.Value, &types.TypeType{Underlying: forwardDecl}, stmt.Name)
		}
	} else if _, isFunc := expectedType.(*types.FuncType); isFunc {
		chk.env.Set(stmt.Name.Value, expectedType, stmt.Name)
	}

	// установка имен для параметров функции
	funcLit, isFuncLit := stmt.Value.(*ast.FunctionLiteral)
	if isFuncLit {
		if expFnType, ok := expectedType.(*types.FuncType); ok {
			for i, param := range funcLit.Parameters {
				if i < len(expFnType.Params) {
					expFnType.Params[i].Name = param.Value
				}
			}
		}
		// chk.env.Set(stmt.Name.Value, expectedType, stmt.Name)
	}

	actualType := chk.checkExpression(stmt.Value)

	if cachedVal, ok := chk.ComptimeValues[stmt.Value]; ok {
		chk.comptimeEnv.Set(stmt.Name.Value, cachedVal)
	}

	if !types.IsAssignable(expectedType, actualType) {
		if expectedType != nil && actualType != nil {
			// if len(chk.errors) == 0 {
			chk.typeError(fmt.Sprintf("in let expected variable of type %s, got %s", expectedType.Signature(), actualType.Signature()), stmt)
			// }
		}
	}

	if isDeclaringType {
		if actTypeType, ok := actualType.(*types.TypeType); ok {
			actualUnderlying := actTypeType.Underlying

			if symbol, isDeclared := chk.env.Get(stmt.Name.Value); isDeclared {
				if forwardTypeType, ok := symbol.SymbolType.(*types.TypeType); ok && forwardTypeType.Underlying != nil {
					switch fwd := forwardTypeType.Underlying.(type) {
					case *types.InterfaceType:
						if act, ok := actualUnderlying.(*types.InterfaceType); ok {
							fwd.Fields = act.Fields
							act.Name = stmt.Name.Value
							actualType = symbol.SymbolType
						}
					case *types.ImplType:
						if act, ok := actualUnderlying.(*types.ImplType); ok {
							fwd.ImplementedType = act.ImplementedType
							fwd.Fields = act.Fields
							fwd.FieldsDefaultValues = act.FieldsDefaultValues
							act.Name = stmt.Name.Value
							actualType = symbol.SymbolType
						}
					}
				}
			}

			// если не объявлен предварительно
			if interfaceType, ok := actualUnderlying.(*types.InterfaceType); ok && interfaceType.Name == "" {
				interfaceType.Name = stmt.Name.Value
			}

			if implType, ok := actualUnderlying.(*types.ImplType); ok && implType.Name == "" {
				implType.Name = stmt.Name.Value
			}
		}
		chk.env.Set(stmt.Name.Value, actualType, stmt.Name)
		chk.recordType(stmt.Name, actualType)
	} else if isFuncLit {
		if funcExpType, isDeclaringFunc := expectedType.(*types.FuncType); isDeclaringFunc {
			chk.env.Set(stmt.Name.Value, funcExpType, stmt.Name)
			chk.recordType(stmt.Name, funcExpType)
		} else {
			chk.env.Set(stmt.Name.Value, actualType, stmt.Name)
			chk.recordType(stmt.Name, actualType)
		}
	} else {
		chk.env.Set(stmt.Name.Value, expectedType, stmt.Name)
	}

	chk.curExpectedType = nil

	return actualType
}

func (chk *TypeChecker) checkReturnStatement(stmt *ast.ReturnStatement) types.Type {
	if stmt == nil {
		return &types.IllegalType{}
	}

	if stmt.Value == nil {
		chk.typeError("missing val in return statement", stmt)
		return &types.IllegalType{}
	}

	if chk.curExpectedType != nil {
		chk.recordExpectedType(stmt.Value, chk.curExpectedType)
	}

	returnType := chk.checkExpression(stmt.Value)
	return returnType
}

func (chk *TypeChecker) checkBlockStatement(stmt *ast.BlockStatement) types.Type {
	if stmt == nil {
		chk.typeError("missing block statement", stmt)
		return &types.IllegalType{}
	}

	if stmt.Statements == nil {
		chk.typeError("missing body in block statement", stmt)
		return &types.IllegalType{}
	}

	chk.env = types.NewEnclosedTypeEviroment(chk.env)
	chk.recordScope(stmt, chk.env)

	var returnType types.Type = &types.IllegalType{}

	for _, stmt := range stmt.Statements {
		switch stmt.(type) {
		case *ast.ReturnStatement:
			returnType = chk.checkStatement(stmt)
		default:
			chk.checkStatement(stmt)
		}
	}

	chk.env = chk.env.Outer

	return returnType
}
