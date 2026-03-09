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
	chk.errors = append(chk.errors, fmt.Errorf("unknown statement type"))
	return &types.IllegalType{}
}

func (chk *TypeChecker) checkLetStatement(stmt *ast.LetStatement) types.Type {
	expectedType := chk.resolveType(stmt.Type)

	chk.curExpectedType = expectedType

	chk.env.Set(stmt.Name.Value, expectedType)

	actualType := chk.checkExpression(stmt.Value)

	if !types.Equals(expectedType, actualType) {
		if expectedType != nil && actualType != nil {
			if len(chk.errors) == 0 {
				chk.errors = append(chk.errors, fmt.Errorf("expected type %s, got %s",
					expectedType.Signature(), actualType.Signature()))
			}
		}
	}

	return actualType
}

func (chk *TypeChecker) checkReturnStatement(stmt *ast.ReturnStatement) types.Type {
	returnType := chk.checkExpression(stmt.Value)
	return returnType
}

func (chk *TypeChecker) checkBlockStatement(stmt *ast.BlockStatement) types.Type {
	chk.env = types.NewEnclosedTypeEviroment(chk.env)

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
